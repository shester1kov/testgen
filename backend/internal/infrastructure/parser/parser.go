package parser

import (
	"fmt"
	"io"
)

// DocumentParser defines the interface for document parsing
type DocumentParser interface {
	Parse(reader io.Reader) (string, error)
	SupportedType() string
}

// DocumentParserFactory creates parsers based on file type (Factory Pattern)
type DocumentParserFactory struct {
	parsers map[string]DocumentParser
}

// NewDocumentParserFactory creates a new parser factory
func NewDocumentParserFactory() *DocumentParserFactory {
	factory := &DocumentParserFactory{
		parsers: make(map[string]DocumentParser),
	}

	// Register all available parsers
	factory.Register(NewPDFParser())
	factory.Register(NewDOCXParser())
	factory.Register(NewPPTXParser())
	factory.Register(NewTXTParser())
	factory.Register(NewMDParser())

	return factory
}

// Register adds a parser to the factory
func (f *DocumentParserFactory) Register(parser DocumentParser) {
	f.parsers[parser.SupportedType()] = parser
}

// CreateParser returns appropriate parser for the file type
func (f *DocumentParserFactory) CreateParser(fileType string) (DocumentParser, error) {
	parser, exists := f.parsers[fileType]
	if !exists {
		return nil, fmt.Errorf("unsupported file type: %s", fileType)
	}
	return parser, nil
}

// GetSupportedTypes returns list of supported file types
func (f *DocumentParserFactory) GetSupportedTypes() []string {
	types := make([]string, 0, len(f.parsers))
	for fileType := range f.parsers {
		types = append(types, fileType)
	}
	return types
}
