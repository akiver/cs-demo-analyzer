package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/akiver/cs-demo-analyzer/pkg/api"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
)

type cliArgs struct {
	demoPath         string
	includePositions bool
	source           string
	outputPath       string
	format           string
	minifyJSON       bool
}

func (cli *cliArgs) validateArgs() error {
	if cli.demoPath == "" {
		return errors.New("demo file path required, example: -demo-path path/to/demo.dem")
	}

	if cli.outputPath == "" {
		return errors.New("output path required, example: -output ./output")
	}

	if cli.format != "" {
		err := api.ValidateExportFormat(constants.ExportFormat(cli.format))
		if err != nil {
			return err
		}
	}

	if cli.source != "" {
		err := api.ValidateDemoSource(constants.DemoSource(cli.source))
		if err != nil {
			return err
		}
	}

	return nil
}

func (cli *cliArgs) fromArgs(args []string) error {
	fs := flag.NewFlagSet("csda", flag.ContinueOnError)
	fs.StringVar(&cli.demoPath, "demo-path", "", "Demo file path (mandatory)")
	fs.StringVar(&cli.outputPath, "output", "", "Output folder or file path, must be a folder when exporting to CSV (mandatory)")
	fs.StringVar(&cli.format, "format", "csv", "Export format, valid values: "+api.FormatValidExportFormats())
	fs.StringVar(&cli.source, "source", "", "Force demo's source, valid values: "+api.FormatValidDemoSources())
	fs.BoolVar(&cli.includePositions, "positions", false, "Include entities (players, grenades...) positions (default false)")
	fs.BoolVar(&cli.minifyJSON, "minify", false, "Minify JSON file, it has effect only when -format is set to json")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if err := cli.validateArgs(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		fs.Usage()
		return err
	}

	return nil
}

func Run(args []string) int {
	var cli cliArgs
	err := cli.fromArgs(args)
	if err != nil {
		return 2
	}

	err = api.AnalyzeAndExportDemo(cli.demoPath, cli.outputPath, api.AnalyzeAndExportDemoOptions{
		IncludePositions: cli.includePositions,
		Source:           constants.DemoSource(cli.source),
		Format:           constants.ExportFormat(cli.format),
		MinifyJSON:       cli.minifyJSON,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}

	return 0
}
