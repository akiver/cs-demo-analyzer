package api

import (
	"errors"
	"fmt"
	"os"

	"github.com/akiver/cs-demo-analyzer/internal/analyzer/cs2"
	"github.com/akiver/cs-demo-analyzer/internal/analyzer/csgo"
	d "github.com/akiver/cs-demo-analyzer/internal/demo"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
)

type AnalyzeDemoOptions struct {
	IncludePositions bool
	Source           constants.DemoSource
}

type AnalyzeAndExportDemoOptions struct {
	IncludePositions bool
	Source           constants.DemoSource
	Format           constants.ExportFormat
	MinifyJSON       bool
}

func analyzeDemo(demoPath string, options AnalyzeDemoOptions) (*Match, error) {
	if options.Source != "" {
		err := ValidateDemoSource(options.Source)
		if err != nil {
			return nil, err
		}
	}

	demo, err := d.GetDemoFromPath(demoPath)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(demoPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("demo file %q not found", demoPath)
		}
		return nil, err
	}
	defer file.Close()

	source := options.Source
	if source == "" {
		source = d.GetDemoSource(demo)
	}

	// CS2 demos are parsed with demoinfocs-golang v6, CS:GO (Source 1) demos with v4.
	if demo.IsSource2() {
		return cs2.Analyze(file, demo, source, options.IncludePositions)
	}

	return csgo.Analyze(file, demo, source, options.IncludePositions)
}

func AnalyzeDemo(demoPath string, options AnalyzeDemoOptions) (*Match, error) {
	return analyzeDemo(demoPath, options)
}

func AnalyzeAndExportDemo(demoPath string, outputPath string, options AnalyzeAndExportDemoOptions) error {
	var err error
	if options.Format != "" {
		err = ValidateExportFormat(options.Format)
		if err != nil {
			return err
		}
	}

	match, err := analyzeDemo(demoPath, AnalyzeDemoOptions{
		IncludePositions: options.IncludePositions,
		Source:           options.Source,
	})
	if err != nil {
		return err
	}

	switch options.Format {
	case "csv":
		err = exportMatchToCSV(match, outputPath)
	case "json":
		err = exportMatchToJSON(match, outputPath, options.MinifyJSON)
	case "csdm":
		err = exportMatchForCSDM(match, outputPath)
	}

	return err
}
