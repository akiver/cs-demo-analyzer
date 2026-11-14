// Command bench compares the speed (wall time) and RAM usage (peak memory) of the analyzer between two versions.
//
// Usage:
//
//	go run ./scripts/bench [flags] [base-ref] [head-ref]
//
// base-ref is the git ref used as the baseline (default: HEAD~1).
// head-ref is the git ref compared against the baseline (default: the current working tree, uncommitted changes included).
// Flags must be placed before the refs.
//
// Both versions are built from temporary git worktrees, their runs are interleaved and each demo is analyzed once
// before measuring to warm up the file cache. Medians are reported because a single slow run (e.g. the machine doing
// something else) would skew an average. Differences of a few percent are usually noise, compare them with the min-max
// spread of each version.
//
// Peak memory is the max RSS on macOS/Linux and the peak working set on Windows, compare results from the same OS only.
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"time"
)

var defaultDemos = []string{
	"cs2/esplay_ntfNCNcmKCQc_2025_mirage.dem:esplay",
	"cs2/matchzy_bleed_vs_parivision_2024_mirage.dem:matchzy",
	"cs2/5eplay_g161-20240107155319726585099_2023_nuke.dem:5eplay",
	"csgo/ebot_immortals_vs_north_iem_katowice_2017_overpass.dem:ebot",
}

type demo struct {
	path   string // Absolute path.
	name   string
	source string
}

type version struct {
	name   string
	binary string
	output string
}

type measure struct {
	wall time.Duration
	peak uint64 // Bytes.
}

func main() {
	os.Exit(run())
}

func run() int {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: go run ./scripts/bench [flags] [base-ref (default HEAD~1)] [head-ref (default working tree)]\n\nFlags:\n")
		flag.PrintDefaults()
	}
	runs := flag.Int("runs", 5, "Number of runs per demo and version")
	positions := flag.Bool("positions", false, "Analyze with -positions")
	demosFlag := flag.String("demos", strings.Join(defaultDemos, ","), "Comma-separated list of <path relative to cs-demos>:<source>")
	flag.Parse()

	baseRef := "HEAD~1"
	headRef := ""
	if flag.NArg() > 0 {
		baseRef = flag.Arg(0)
	}
	if flag.NArg() > 1 {
		headRef = flag.Arg(1)
	}
	if flag.NArg() > 2 || *runs < 1 {
		flag.Usage()
		return 2
	}

	rootDir, err := git("", "rev-parse", "--show-toplevel")
	if err != nil {
		return fail(err)
	}

	demos, err := parseDemos(rootDir, *demosFlag)
	if err != nil {
		return fail(err)
	}

	workDir, err := os.MkdirTemp("", "csda-bench-")
	if err != nil {
		return fail(err)
	}

	var worktrees []string
	cleanup := func() {
		for _, worktree := range worktrees {
			_, _ = git(rootDir, "worktree", "remove", "--force", worktree)
		}
		_ = os.RemoveAll(workDir)
	}
	defer cleanup()

	// Remove the worktrees on Ctrl+C, the running analysis receives the signal too and stops.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	go func() {
		<-interrupt
		cleanup()
		os.Exit(130)
	}()

	build := func(name string, ref string) (*version, error) {
		srcDir := rootDir
		label := "working tree"
		if ref != "" {
			hash, err := git(rootDir, "rev-parse", "--short", ref)
			if err != nil {
				return nil, err
			}
			srcDir = filepath.Join(workDir, "src-"+name)
			if _, err := git(rootDir, "worktree", "add", "--detach", srcDir, ref); err != nil {
				return nil, err
			}
			worktrees = append(worktrees, srcDir)
			label = fmt.Sprintf("%s (%s)", ref, hash)
		}

		fmt.Printf("Building %s: %s\n", name, label)
		v := &version{
			name:   name,
			binary: filepath.Join(workDir, "csda-"+name),
			output: filepath.Join(workDir, "out-"+name),
		}
		if runtime.GOOS == "windows" {
			v.binary += ".exe"
		}

		cmd := exec.Command("go", "build", "-trimpath", "-o", v.binary, "./cmd/cli")
		cmd.Dir = srcDir
		cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("failed to build %s: %w", label, err)
		}

		return v, os.MkdirAll(v.output, 0o755)
	}

	base, err := build("base", baseRef)
	if err != nil {
		return fail(err)
	}
	head, err := build("head", headRef)
	if err != nil {
		return fail(err)
	}

	results := make(map[string]map[string][]measure)
	for _, d := range demos {
		fmt.Printf("Benchmarking %s (%d runs)\n", d.name, *runs)
		results[d.name] = make(map[string][]measure)

		// Warm up the file cache.
		if _, err := analyze(base, d, *positions); err != nil {
			return fail(err)
		}

		for range *runs {
			for _, v := range []*version{base, head} {
				m, err := analyze(v, d, *positions)
				if err != nil {
					return fail(err)
				}
				results[d.name][v.name] = append(results[d.name][v.name], m)
			}
		}
	}

	fmt.Println()
	printResults(demos, results)

	return 0
}

func parseDemos(rootDir string, value string) ([]demo, error) {
	var demos []demo
	for entry := range strings.SplitSeq(value, ",") {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		path, source, found := strings.Cut(entry, ":")
		if !found || source == "" {
			return nil, fmt.Errorf("invalid demo %q, expected <path relative to cs-demos>:<source>", entry)
		}

		absolutePath := filepath.Join(rootDir, "cs-demos", filepath.FromSlash(path))
		if _, err := os.Stat(absolutePath); err != nil {
			return nil, fmt.Errorf("demo not found: cs-demos/%s, run ./download-demos.sh first", path)
		}

		demos = append(demos, demo{
			path:   absolutePath,
			name:   strings.TrimSuffix(filepath.Base(absolutePath), ".dem"),
			source: source,
		})
	}

	if len(demos) == 0 {
		return nil, errors.New("no demos provided")
	}

	return demos, nil
}

func analyze(v *version, d demo, positions bool) (measure, error) {
	args := []string{"-demo-path", d.path, "-source", d.source, "-output", v.output, "-format", "json", "-minify"}
	if positions {
		args = append(args, "-positions")
	}

	var output bytes.Buffer
	cmd := exec.Command(v.binary, args...)
	cmd.Stdout = &output
	cmd.Stderr = &output

	m, err := runAndMeasure(cmd)
	if err != nil {
		return m, fmt.Errorf("%s failed on %s: %w\n%s", v.name, d.name, err, output.String())
	}

	return m, nil
}

func printResults(demos []demo, results map[string]map[string][]measure) {
	width := len("Demo")
	for _, d := range demos {
		width = max(width, len(d.name))
	}

	fmt.Printf("%-*s %22s %22s %8s %9s %9s %8s\n", width, "Demo", "Wall base (min-max)", "Wall head (min-max)", "Diff", "RAM base", "RAM head", "Diff")
	for _, d := range demos {
		base := summarize(results[d.name]["base"])
		head := summarize(results[d.name]["head"])
		fmt.Printf("%-*s %6.2fs (%5.2f-%5.2f) %6.2fs (%5.2f-%5.2f) %+7.1f%% %7.0fMB %7.0fMB %+7.1f%%\n",
			width, d.name,
			base.wall, base.wallMin, base.wallMax,
			head.wall, head.wallMin, head.wallMax,
			percentDiff(base.wall, head.wall),
			base.peak, head.peak,
			percentDiff(base.peak, head.peak),
		)
	}
}

type summary struct {
	wall    float64 // Median, seconds.
	wallMin float64
	wallMax float64
	peak    float64 // Median, MB.
}

func summarize(measures []measure) summary {
	walls := make([]float64, len(measures))
	peaks := make([]float64, len(measures))
	for i, m := range measures {
		walls[i] = m.wall.Seconds()
		peaks[i] = float64(m.peak) / 1024 / 1024
	}
	slices.Sort(walls)

	return summary{
		wall:    median(walls),
		wallMin: walls[0],
		wallMax: walls[len(walls)-1],
		peak:    median(peaks),
	}
}

func median(values []float64) float64 {
	values = slices.Clone(values)
	slices.Sort(values)
	middle := len(values) / 2
	if len(values)%2 == 1 {
		return values[middle]
	}

	return (values[middle-1] + values[middle]) / 2
}

func percentDiff(base float64, head float64) float64 {
	return (head - base) / base * 100
}

func git(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s: %w\n%s", strings.Join(args, " "), err, output)
	}

	return strings.TrimSpace(string(output)), nil
}

func fail(err error) int {
	fmt.Fprintln(os.Stderr, err)

	return 1
}
