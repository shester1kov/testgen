package exporter

import (
	"fmt"

	"github.com/shester1kov/testgen-backend/internal/domain/entity"
)

// ExportFormat represents supported export formats
type ExportFormat string

const (
	FormatJSON       ExportFormat = "json"
	FormatMoodleXML  ExportFormat = "moodle_xml"
	FormatStepikCSV  ExportFormat = "stepik_csv"
	FormatGIFT       ExportFormat = "gift"
	FormatAiken      ExportFormat = "aiken"
	FormatBlackboard ExportFormat = "blackboard"
	FormatQTI        ExportFormat = "qti"
)

// ExportResult contains the exported content and metadata
type ExportResult struct {
	Content     string
	ContentType string
	FileExt     string
	Format      ExportFormat
}

// Exporter interface for all export formats
type Exporter interface {
	// Export converts test data to the specific format
	Export(test *entity.Test, questions []*entity.Question, answers map[string][]*entity.Answer) (*ExportResult, error)
	// Format returns the export format identifier
	Format() ExportFormat
	// Name returns human-readable format name
	Name() string
	// Description returns format description
	Description() string
	// SupportedQuestionTypes returns list of supported question types
	SupportedQuestionTypes() []entity.QuestionType
}

// ExporterFactory creates exporters for different formats
type ExporterFactory struct {
	exporters map[ExportFormat]Exporter
}

// NewExporterFactory creates a new exporter factory with all available exporters
func NewExporterFactory() *ExporterFactory {
	factory := &ExporterFactory{
		exporters: make(map[ExportFormat]Exporter),
	}

	// Register all exporters
	factory.Register(NewJSONExporter())
	factory.Register(NewMoodleXMLExporter())
	factory.Register(NewStepikCSVExporter())
	factory.Register(NewGIFTExporter())
	factory.Register(NewAikenExporter())
	factory.Register(NewBlackboardExporter())
	factory.Register(NewQTIExporter())

	return factory
}

// Register adds an exporter to the factory
func (f *ExporterFactory) Register(exporter Exporter) {
	f.exporters[exporter.Format()] = exporter
}

// GetExporter returns an exporter for the specified format
func (f *ExporterFactory) GetExporter(format ExportFormat) (Exporter, error) {
	exporter, ok := f.exporters[format]
	if !ok {
		return nil, fmt.Errorf("unsupported export format: %s", format)
	}
	return exporter, nil
}

// GetAvailableFormats returns list of all available export formats
func (f *ExporterFactory) GetAvailableFormats() []FormatInfo {
	formats := make([]FormatInfo, 0, len(f.exporters))
	for _, exporter := range f.exporters {
		formats = append(formats, FormatInfo{
			Format:      exporter.Format(),
			Name:        exporter.Name(),
			Description: exporter.Description(),
		})
	}
	return formats
}

// FormatInfo contains information about an export format
type FormatInfo struct {
	Format      ExportFormat `json:"format"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
}
