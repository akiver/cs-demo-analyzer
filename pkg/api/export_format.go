package api

import (
	"fmt"
	"strings"

	"github.com/akiver/cs-demo-analyzer/internal/slice"
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
)

func FormatValidExportFormats() string {
	var formats []string
	for _, format := range constants.ExportFormats {
		formats = append(formats, string(format))
	}

	return "[" + strings.Join(formats, ",") + "]"
}

func ValidateExportFormat(format constants.ExportFormat) error {
	isValid := slice.Contains(constants.ExportFormats, constants.ExportFormat(format))
	if isValid {
		return nil
	}

	return fmt.Errorf("invalid format provided, valid formats: %s", FormatValidExportFormats())
}
