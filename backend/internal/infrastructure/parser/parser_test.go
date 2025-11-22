package parser

import (
	"bytes"
	"strings"
	"testing"
)

// POSITIVE TEST: Factory creates parser for supported types
func TestDocumentParserFactory_CreateParser_Success(t *testing.T) {
	factory := NewDocumentParserFactory()

	supportedTypes := []string{"pdf", "docx", "pptx", "txt", "md"}
	for _, fileType := range supportedTypes {
		t.Run(fileType, func(t *testing.T) {
			parser, err := factory.CreateParser(fileType)
			if err != nil {
				t.Errorf("Failed to create parser for %s: %v", fileType, err)
			}
			if parser == nil {
				t.Errorf("Parser for %s is nil", fileType)
			}
			if parser.SupportedType() != fileType {
				t.Errorf("Expected parser type %s, got %s", fileType, parser.SupportedType())
			}
		})
	}
}

// NEGATIVE TEST: Factory fails for unsupported file type
func TestDocumentParserFactory_CreateParser_UnsupportedType(t *testing.T) {
	factory := NewDocumentParserFactory()

	unsupportedTypes := []string{
		"exe",
		"mp3",
		"jpg",
		"png",
		"zip",
		"unknown",
		"",
		"PDF",     // case sensitive
		" pdf",    // with whitespace
		"pdf ",    // trailing whitespace
		"doc",     // old format
		"xls",     // Excel
		"xlsx",    // Excel
		"csv",
	}

	for _, fileType := range unsupportedTypes {
		t.Run("unsupported_"+fileType, func(t *testing.T) {
			parser, err := factory.CreateParser(fileType)
			if err == nil {
				t.Errorf("Expected error for unsupported file type '%s', got nil", fileType)
			}
			if parser != nil {
				t.Errorf("Expected nil parser for unsupported file type '%s', got %v", fileType, parser)
			}
			if err != nil && !strings.Contains(err.Error(), "unsupported file type") {
				t.Errorf("Expected error message to contain 'unsupported file type', got: %s", err.Error())
			}
		})
	}
}

// NEGATIVE TEST: TXT Parser with empty content
func TestTXTParser_Parse_EmptyContent(t *testing.T) {
	parser := NewTXTParser()
	reader := bytes.NewReader([]byte(""))

	text, err := parser.Parse(reader)
	if err != nil {
		t.Errorf("Expected no error for empty content, got: %v", err)
	}
	if text != "" {
		t.Errorf("Expected empty string for empty content, got: %s", text)
	}
}

// NEGATIVE TEST: TXT Parser with nil reader
func TestTXTParser_Parse_NilReader(t *testing.T) {
	parser := NewTXTParser()

	// This should panic or error gracefully
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil reader, got none")
		}
	}()

	parser.Parse(nil)
}

// NEGATIVE TEST: TXT Parser with very large file
func TestTXTParser_Parse_VeryLargeContent(t *testing.T) {
	parser := NewTXTParser()

	// Create 10MB of text
	largeContent := strings.Repeat("A", 10*1024*1024)
	reader := bytes.NewReader([]byte(largeContent))

	text, err := parser.Parse(reader)
	if err != nil {
		t.Errorf("Failed to parse large content: %v", err)
	}
	if len(text) != len(largeContent) {
		t.Errorf("Expected text length %d, got %d", len(largeContent), len(text))
	}
}

// NEGATIVE TEST: TXT Parser with special characters
func TestTXTParser_Parse_SpecialCharacters(t *testing.T) {
	parser := NewTXTParser()

	testCases := []struct {
		name    string
		content string
	}{
		{"null bytes", "Hello\x00World"},
		{"unicode", "Hello 世界 Мир"},
		{"emojis", "Hello 😀 👍 🎉"},
		{"control characters", "Hello\r\n\tWorld"},
		{"mixed newlines", "Hello\nWorld\r\nTest\rEnd"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := bytes.NewReader([]byte(tc.content))
			text, err := parser.Parse(reader)
			if err != nil {
				t.Errorf("Failed to parse %s: %v", tc.name, err)
			}
			if text == "" {
				t.Error("Expected non-empty text")
			}
		})
	}
}

// NEGATIVE TEST: PDF Parser with corrupted data
// Note: Current PDF parser is a placeholder - returns success with TODO message
// TODO: In production, implement real PDF parsing that validates format
func TestPDFParser_Parse_CorruptedData(t *testing.T) {
	parser := NewPDFParser()

	corruptedData := []byte("This is not a valid PDF file")
	reader := bytes.NewReader(corruptedData)

	text, err := parser.Parse(reader)
	// Current implementation returns placeholder message, no error
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	// Verify it returns placeholder message
	if !strings.Contains(text, "TODO") {
		t.Error("Expected placeholder message containing 'TODO'")
	}
}

// NEGATIVE TEST: PDF Parser with empty data
// Note: Placeholder returns success - real implementation should validate
func TestPDFParser_Parse_EmptyData(t *testing.T) {
	parser := NewPDFParser()

	reader := bytes.NewReader([]byte(""))

	text, err := parser.Parse(reader)
	// Placeholder returns success with TODO message
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if text == "" {
		t.Error("Expected non-empty placeholder message")
	}
}

// NEGATIVE TEST: PDF Parser with partial PDF header
// Note: Future implementation should validate PDF format
func TestPDFParser_Parse_PartialHeader(t *testing.T) {
	parser := NewPDFParser()

	// Valid PDF starts with "%PDF-", but this is incomplete
	partialData := []byte("%PDF")
	reader := bytes.NewReader(partialData)

	text, err := parser.Parse(reader)
	// Placeholder returns success
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if text == "" {
		t.Error("Expected non-empty placeholder message")
	}
}

// NEGATIVE TEST: DOCX Parser with non-docx data
// Note: Current DOCX parser is also a placeholder
func TestDOCXParser_Parse_InvalidData(t *testing.T) {
	parser := NewDOCXParser()

	invalidData := []byte("This is not a DOCX file")
	reader := bytes.NewReader(invalidData)

	text, err := parser.Parse(reader)
	// Placeholder returns success with TODO message
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !strings.Contains(text, "TODO") {
		t.Error("Expected placeholder message containing 'TODO'")
	}
}

// NEGATIVE TEST: DOCX Parser with corrupted ZIP structure
// Note: Real implementation should validate ZIP/DOCX structure
func TestDOCXParser_Parse_CorruptedZip(t *testing.T) {
	parser := NewDOCXParser()

	// DOCX is a ZIP file, but this is corrupted
	corruptedZip := []byte("PK\x03\x04CORRUPTED")
	reader := bytes.NewReader(corruptedZip)

	text, err := parser.Parse(reader)
	// Placeholder returns success
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if text == "" {
		t.Error("Expected non-empty placeholder message")
	}
}

// NEGATIVE TEST: PPTX Parser with non-pptx data
// Note: Current PPTX parser is also a placeholder
func TestPPTXParser_Parse_InvalidData(t *testing.T) {
	parser := NewPPTXParser()

	invalidData := []byte("This is not a PPTX file")
	reader := bytes.NewReader(invalidData)

	text, err := parser.Parse(reader)
	// Placeholder returns success with TODO message
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !strings.Contains(text, "TODO") {
		t.Error("Expected placeholder message containing 'TODO'")
	}
}

// NEGATIVE TEST: Parser Factory GetSupportedTypes
func TestDocumentParserFactory_GetSupportedTypes(t *testing.T) {
	factory := NewDocumentParserFactory()

	supportedTypes := factory.GetSupportedTypes()
	if len(supportedTypes) == 0 {
		t.Error("Expected at least one supported type, got none")
	}

	// Verify expected types are present
	expectedTypes := map[string]bool{
		"pdf":  false,
		"docx": false,
		"pptx": false,
		"txt":  false,
	}

	for _, fileType := range supportedTypes {
		if _, exists := expectedTypes[fileType]; exists {
			expectedTypes[fileType] = true
		}
	}

	for fileType, found := range expectedTypes {
		if !found {
			t.Errorf("Expected file type '%s' not found in supported types", fileType)
		}
	}
}

// NEGATIVE TEST: Register duplicate parser
func TestDocumentParserFactory_Register_Duplicate(t *testing.T) {
	factory := NewDocumentParserFactory()

	// Register a parser twice
	txtParser1 := NewTXTParser()
	txtParser2 := NewTXTParser()

	factory.Register(txtParser1)
	initialCount := len(factory.GetSupportedTypes())

	// Registering again should overwrite, not add
	factory.Register(txtParser2)
	finalCount := len(factory.GetSupportedTypes())

	if initialCount != finalCount {
		t.Error("Expected same number of parsers after duplicate registration")
	}
}

// NEGATIVE TEST: Binary data that looks like text
func TestTXTParser_Parse_BinaryData(t *testing.T) {
	parser := NewTXTParser()

	// Random binary data
	binaryData := []byte{0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE, 0xFD, 0xFC}
	reader := bytes.NewReader(binaryData)

	// TXT parser should accept any data, even binary
	text, err := parser.Parse(reader)
	if err != nil {
		t.Errorf("TXT parser should accept binary data, got error: %v", err)
	}
	if len(text) == 0 {
		t.Error("Expected non-empty text from binary data")
	}
}
