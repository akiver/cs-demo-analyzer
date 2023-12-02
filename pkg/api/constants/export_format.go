package constants

type ExportFormat string

const (
	ExportFormatCSV  ExportFormat = "csv"
	ExportFormatJSON ExportFormat = "json"
	ExportFormatCSDM ExportFormat = "csdm" // Special CSV export dedicated to the application CS Demo Manager
)

var ExportFormats = []ExportFormat{
	ExportFormatCSV,
	ExportFormatJSON,
	ExportFormatCSDM,
}
