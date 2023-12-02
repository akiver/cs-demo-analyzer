package api

import (
	"fmt"
	"strings"

	"github.com/akiver/cs-demo-analyzer/internal/slice"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
)

func FormatValidDemoSources() string {
	var sources []string
	for _, source := range constants.SupportedDemoSources {
		sources = append(sources, string(source))
	}

	return "[" + strings.Join(sources, ",") + "]"
}

func ValidateDemoSource(source constants.DemoSource) error {
	isValid := slice.Contains(constants.SupportedDemoSources, source)
	if isValid {
		return nil
	}

	return fmt.Errorf("invalid source provided, valid sources: %s", FormatValidDemoSources())
}
