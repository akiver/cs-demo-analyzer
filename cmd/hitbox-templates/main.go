// Generates the hitbox templates of the visibility check (internal/visibility/hitbox_templates.go) from hitbox capsules
// captured in game.
//
// The capsules are captured with the csdm_dump_hitboxes command of the CS Demo Manager CS2 plugin, which only exists in
// its dev builds.
// It appends the world-space hitbox capsules of the alive players at the current tick to csdm_hitboxes.jsonl, next to
// csdm.log. Run it at many ticks of a demo, for example from the cs-demo-manager repository with a Debug build of the
// plugin:
//
//	node out/cli.js watch demo.dem --tick 4000 --exec "96:demo_timescale 3" --exec 4024:csdm_dump_hitboxes \
//	  --exec 4048:csdm_dump_hitboxes ... --exec 22100:quit
//
// Then, with the captures of each demo next to the demo itself:
//
//	go run ./cmd/hitbox-templates [-report] capture1.jsonl demo1.dem [capture2.jsonl demo2.dem]...
//
// -report prints how far the points of templates fitted on the other demos are from the real body of the players of
// each demo.
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/format"
	"math"
	"os"
	"sort"
	"strconv"

	"github.com/akiver/cs-demo-analyzer/internal/visibility"
	"github.com/golang/geo/r3"
	dem "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

// A template is fitted on fewer samples than that only when no coarser template has more.
const minSamplesPerTemplate = 30

// Points closer than that to a previous point of the template are dropped.
const minPointsDistance = 3.0

// Hitboxes in the order their points are tested, from the most to the least likely visible part when peeking.
var hitboxOrder = []string{
	"head_0", "neck_0", "spine_3", "spine_2", "spine_1", "spine_0", "pelvis", "arm_upper_l", "arm_upper_r",
	"arm_lower_l", "arm_lower_r", "hand_l", "hand_r", "leg_upper_l", "leg_upper_r", "leg_lower_l", "leg_lower_r",
	"ankle_l", "ankle_r",
}

type hitboxDescription struct {
	Name   string  `json:"name"`
	Radius float64 `json:"radius"`
}

type captureLine struct {
	Type     string              `json:"type"`
	Tick     int                 `json:"tick"`
	SteamID  string              `json:"steamId"`
	Origin   [3]float64          `json:"origin"`
	Root     [3]float64          `json:"root"`
	Capsules [][7]float64        `json:"capsules"`
	Hitboxes []hitboxDescription `json:"hitboxes"`
}

type capsule struct {
	a, b   r3.Vector
	radius float64
}

type sample struct {
	demo       int
	template   int
	eyesHeight float64
	// Capsules in the frame of the player: forward, right, up from the feet.
	capsules []capsule
}

func main() {
	report := flag.Bool("report", false, "print the error of the templates on each demo, fitted on the other demos")
	output := flag.String("out", "internal/visibility/hitbox_templates.go", "generated file")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 || len(args)%2 != 0 {
		fmt.Println("usage: hitbox-templates [-report] [-out file] <capture.jsonl> <demo.dem> [<capture.jsonl> <demo.dem>]...")
		os.Exit(1)
	}

	var hitboxes []hitboxDescription
	var samples []sample
	for i := 0; i < len(args); i += 2 {
		captures, descriptions := readCaptures(args[i])
		if len(descriptions) == 0 {
			fmt.Printf("%s: no hitboxes description\n", args[i])
			os.Exit(1)
		}
		if hitboxes == nil {
			hitboxes = descriptions
		} else if !sameHitboxes(hitboxes, descriptions) {
			fmt.Printf("%s: the hitboxes differ from the ones of %s\n", args[i], args[0])
			os.Exit(1)
		}
		samples = append(samples, joinDemo(i/2, args[i+1], captures)...)
	}

	order := hitboxIndices(hitboxes)
	if *report {
		reportErrors(samples, hitboxes, order, len(args)/2)
	}

	templates := fitTemplates(samples, len(hitboxes), order)
	if err := os.WriteFile(*output, generate(templates, len(samples), len(args)/2), 0o644); err != nil {
		panic(err)
	}
	fmt.Printf("%d samples, %s written\n", len(samples), *output)
}

func readCaptures(path string) (map[string]captureLine, []hitboxDescription) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	captures := make(map[string]captureLine)
	var hitboxes []hitboxDescription
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1<<20), 1<<20)
	for scanner.Scan() {
		var line captureLine
		if err := json.Unmarshal(scanner.Bytes(), &line); err != nil {
			panic(err)
		}
		if line.Type == "model" {
			hitboxes = line.Hitboxes
			continue
		}
		// The root bone far from the origin means bones that weren't set up for this frame.
		if math.Hypot(line.Root[0]-line.Origin[0], line.Root[1]-line.Origin[1]) > 8 {
			continue
		}
		captures[captureKey(line.Tick, line.SteamID)] = line
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	// A pose without one capsule per hitbox, such as a truncated line, would shift the capsules or pull the means of the
	// missing ones toward the feet.
	skippedCount := 0
	for key, line := range captures {
		if len(line.Capsules) != len(hitboxes) {
			delete(captures, key)
			skippedCount++
		}
	}
	if skippedCount > 0 {
		fmt.Printf("%s: %d poses skipped, not %d capsules\n", path, skippedCount, len(hitboxes))
	}

	return captures, hitboxes
}

// Every agent model shares the same hitboxes, a different set would need its own templates.
func sameHitboxes(a []hitboxDescription, b []hitboxDescription) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func hitboxIndices(hitboxes []hitboxDescription) []int {
	var order []int
	for _, name := range hitboxOrder {
		for i, hitbox := range hitboxes {
			if hitbox.Name == name {
				order = append(order, i)
			}
		}
	}
	if len(order) != len(hitboxes) {
		panic(fmt.Sprintf("the hitboxes changed, update hitboxOrder: %v", hitboxes))
	}

	return order
}

func captureKey(tick int, steamID string) string {
	return strconv.Itoa(tick) + "/" + steamID
}

// Joins the captures with the state of the players in the demo, the plugin logs the in-game tick.
func joinDemo(demoIndex int, path string, captures map[string]captureLine) []sample {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	parser := dem.NewParser(file)
	defer parser.Close()

	var samples []sample
	captureCount := len(captures)
	var originDistanceSum float64
	parser.RegisterEventHandler(func(events.FrameDone) {
		tick := parser.GameState().IngameTick()
		for _, player := range parser.GameState().Participants().Playing() {
			key := captureKey(tick, strconv.FormatUint(player.SteamID64, 10))
			capture, ok := captures[key]
			pawn := player.PlayerPawnEntity()
			if !ok || !player.IsAlive() || pawn == nil {
				continue
			}
			delete(captures, key)

			// Same snapshot as the visibility check.
			origin := pawn.Position()
			eyesHeight := 64.0
			if viewOffset, ok := pawn.PropertyValue("m_vecViewOffset.m_vecZ"); ok && viewOffset.Any != nil && viewOffset.Float() > 0 {
				eyesHeight = float64(viewOffset.Float())
			}
			stance := visibility.StanceOf(eyesHeight, player.IsAirborne())
			isMoving := visibility.IsMoving(player.Velocity())
			weaponKind := visibility.WeaponKindOf(player.ActiveWeapon(), player.IsScoped())

			yaw := float64(player.ViewDirectionX()) * math.Pi / 180
			forward := r3.Vector{X: math.Cos(yaw), Y: math.Sin(yaw)}
			right := r3.Vector{X: math.Sin(yaw), Y: -math.Cos(yaw)}
			toLocal := func(x float64, y float64, z float64) r3.Vector {
				d := r3.Vector{X: x, Y: y, Z: z}.Sub(origin)
				return r3.Vector{X: d.Dot(forward), Y: d.Dot(right), Z: d.Z}
			}

			s := sample{
				demo:       demoIndex,
				template:   visibility.HitboxTemplateIndex(stance, isMoving, weaponKind),
				eyesHeight: eyesHeight,
			}
			for _, c := range capture.Capsules {
				s.capsules = append(s.capsules, capsule{
					a:      toLocal(c[0], c[1], c[2]),
					b:      toLocal(c[3], c[4], c[5]),
					radius: c[6],
				})
			}
			originDistanceSum += origin.Sub(r3.Vector{X: capture.Origin[0], Y: capture.Origin[1], Z: capture.Origin[2]}).Norm()
			samples = append(samples, s)
		}
	})

	if err := parser.ParseToEnd(); err != nil {
		fmt.Printf("%s: %v\n", path, err)
	}

	fmt.Printf("%s: %d/%d captures joined, mean origin distance %.2f\n", path, len(samples), captureCount,
		originDistanceSum/math.Max(1, float64(len(samples))))

	return samples
}

type template struct {
	sampleCount int
	eyesHeight  float64
	capsules    []capsule
	points      []visibility.HitboxPoint
}

type accumulator struct {
	count      int
	eyesHeight float64
	a, b       []r3.Vector
	radius     []float64
}

func (acc *accumulator) add(s sample) {
	if acc.a == nil {
		acc.a = make([]r3.Vector, len(s.capsules))
		acc.b = make([]r3.Vector, len(s.capsules))
		acc.radius = make([]float64, len(s.capsules))
	}
	acc.count++
	acc.eyesHeight += s.eyesHeight
	for i, c := range s.capsules {
		acc.a[i] = acc.a[i].Add(c.a)
		acc.b[i] = acc.b[i].Add(c.b)
		acc.radius[i] = c.radius
	}
}

// Fits each template with the mean capsules of its samples. A template with too few samples takes the ones of the
// same stance and weapon kind, then of the same stance and moving state, then of the same stance.
func fitTemplates(samples []sample, hitboxCount int, order []int) []template {
	type weaponKey struct {
		stance     visibility.Stance
		weaponKind visibility.WeaponKind
	}
	type movingKey struct {
		stance   visibility.Stance
		isMoving bool
	}
	exact := make([]accumulator, visibility.HitboxTemplateCount)
	byWeapon := make(map[weaponKey]*accumulator)
	byMoving := make(map[movingKey]*accumulator)
	byStance := make(map[visibility.Stance]*accumulator)
	for _, s := range samples {
		stance, isMoving, weaponKind := visibility.SplitHitboxTemplateIndex(s.template)
		exact[s.template].add(s)
		getAccumulator(byWeapon, weaponKey{stance, weaponKind}).add(s)
		getAccumulator(byMoving, movingKey{stance, isMoving}).add(s)
		getAccumulator(byStance, stance).add(s)
	}

	templates := make([]template, visibility.HitboxTemplateCount)
	for index := range templates {
		stance, isMoving, weaponKind := visibility.SplitHitboxTemplateIndex(index)
		candidates := []*accumulator{
			&exact[index], byWeapon[weaponKey{stance, weaponKind}], byMoving[movingKey{stance, isMoving}], byStance[stance],
		}
		var acc *accumulator
		for _, candidate := range candidates {
			if candidate != nil && candidate.count >= minSamplesPerTemplate {
				acc = candidate
				break
			}
		}
		if acc == nil {
			acc = byStance[stance]
		}
		if acc == nil {
			acc = byStance[visibility.StanceStanding]
		}
		if acc == nil {
			panic("no standing sample")
		}

		t := template{sampleCount: acc.count, eyesHeight: acc.eyesHeight / float64(acc.count)}
		for h := 0; h < hitboxCount; h++ {
			t.capsules = append(t.capsules, capsule{
				a:      acc.a[h].Mul(1 / float64(acc.count)),
				b:      acc.b[h].Mul(1 / float64(acc.count)),
				radius: acc.radius[h],
			})
		}
		t.points = templatePoints(t.capsules, order)
		templates[index] = t
	}

	return templates
}

func getAccumulator[K comparable](accumulators map[K]*accumulator, key K) *accumulator {
	if accumulators[key] == nil {
		accumulators[key] = &accumulator{}
	}

	return accumulators[key]
}

// Points along the axis of each capsule, the center of the short ones and the ends and middle of the long ones. The
// head also gets a point near its top to see players peeking over cover.
func templatePoints(capsules []capsule, order []int) []visibility.HitboxPoint {
	var points []r3.Vector
	addPoint := func(point r3.Vector) {
		for _, existing := range points {
			if existing.Sub(point).Norm() < minPointsDistance {
				return
			}
		}
		points = append(points, point)
	}

	for rank, h := range order {
		c := capsules[h]
		length := c.b.Sub(c.a).Norm()
		var ts []float64
		switch {
		case length < 6:
			ts = []float64{0.5}
		case length < 14:
			ts = []float64{0, 1}
		default:
			ts = []float64{0, 0.5, 1}
		}
		if rank == 0 {
			top := c.a
			if c.b.Z > top.Z {
				top = c.b
			}
			addPoint(top.Add(r3.Vector{Z: c.radius / 2}))
		}
		for _, t := range ts {
			addPoint(c.a.Add(c.b.Sub(c.a).Mul(t)))
		}
	}

	result := make([]visibility.HitboxPoint, len(points))
	for i, point := range points {
		result[i] = visibility.HitboxPoint{Forward: point.X, Right: point.Y, Up: point.Z}
	}

	return result
}

func generate(templates []template, sampleCount int, demoCount int) []byte {
	var buffer bytes.Buffer
	fmt.Fprintf(&buffer, "// Code generated by cmd/hitbox-templates from %d poses of %d demos, DO NOT EDIT.\n\n", sampleCount,
		demoCount)
	buffer.WriteString("package visibility\n\n")
	buffer.WriteString("var hitboxTemplates = [HitboxTemplateCount]HitboxTemplate{\n")
	for index, t := range templates {
		fmt.Fprintf(&buffer, "// %s, %d poses\n", visibility.HitboxTemplateName(index), t.sampleCount)
		fmt.Fprintf(&buffer, "{EyesHeight: %.2f, Points: []HitboxPoint{\n", t.eyesHeight)
		for _, point := range t.points {
			fmt.Fprintf(&buffer, "{%.1f, %.1f, %.1f},\n", point.Forward, point.Right, point.Up)
		}
		buffer.WriteString("}},\n")
	}
	buffer.WriteString("}\n")

	source, err := format.Source(buffer.Bytes())
	if err != nil {
		panic(err)
	}

	return source
}

// Position of a template point for the eyes height of the sample, like the visibility check.
func pointAt(t template, point visibility.HitboxPoint, eyesHeight float64) r3.Vector {
	return r3.Vector{X: point.Forward, Y: point.Right, Z: point.Up * eyesHeight / t.eyesHeight}
}

func segmentDistance(point r3.Vector, c capsule) float64 {
	segment := c.b.Sub(c.a)
	t := 0.0
	if lengthSquared := segment.Norm2(); lengthSquared > 0 {
		t = math.Max(0, math.Min(1, point.Sub(c.a).Dot(segment)/lengthSquared))
	}

	return point.Sub(c.a.Add(segment.Mul(t))).Norm()
}

func percentile(values []float64, p float64) float64 {
	if len(values) == 0 {
		return math.NaN()
	}
	sorted := append([]float64(nil), values...)
	sort.Float64s(sorted)

	return sorted[int(math.Min(float64(len(sorted)-1), p*float64(len(sorted))))]
}

// For each demo, fits the templates on the other demos and measures on the players of the demo how far the points are
// from the body (outside when more than 2 units from the closest capsule surface) and how far each capsule is from
// the closest point.
func reportErrors(samples []sample, hitboxes []hitboxDescription, order []int, demoCount int) {
	if demoCount < 2 {
		fmt.Println("-report needs at least 2 demos")
		return
	}

	var distances []float64
	var pointCounts []float64
	outside := 0
	coverage := make([][]float64, len(hitboxes))
	for demo := 0; demo < demoCount; demo++ {
		var training, test []sample
		for _, s := range samples {
			if s.demo == demo {
				test = append(test, s)
			} else {
				training = append(training, s)
			}
		}
		templates := fitTemplates(training, len(hitboxes), order)
		for _, s := range test {
			t := templates[s.template]
			var points []r3.Vector
			for _, point := range t.points {
				position := pointAt(t, point, s.eyesHeight)
				points = append(points, position)
				distance := math.Inf(1)
				for _, c := range s.capsules {
					distance = math.Min(distance, segmentDistance(position, c)-c.radius)
				}
				distances = append(distances, distance)
				if distance > 2 {
					outside++
				}
			}
			pointCounts = append(pointCounts, float64(len(points)))
			for h, c := range s.capsules {
				for step := 0.0; step <= 1; step += 0.25 {
					axisPoint := c.a.Add(c.b.Sub(c.a).Mul(step))
					closest := math.Inf(1)
					for _, point := range points {
						closest = math.Min(closest, point.Sub(axisPoint).Norm())
					}
					coverage[h] = append(coverage[h], closest-c.radius)
				}
			}
		}
	}

	fmt.Printf("Points per template: p50 %.0f, max %.0f\n", percentile(pointCounts, 0.5), percentile(pointCounts, 1))
	fmt.Printf("Signed distance of the points to the body: p10 %.1f p50 %.1f p90 %.1f, %.0f%% outside\n",
		percentile(distances, 0.1), percentile(distances, 0.5), percentile(distances, 0.9),
		100*float64(outside)/float64(len(distances)))
	fmt.Println("Distance from each capsule axis to the closest point, minus the radius (p50 p90 p99):")
	for h, hitbox := range hitboxes {
		fmt.Printf("  %-12s %5.1f %5.1f %5.1f\n", hitbox.Name, percentile(coverage[h], 0.5), percentile(coverage[h], 0.9),
			percentile(coverage[h], 0.99))
	}
}
