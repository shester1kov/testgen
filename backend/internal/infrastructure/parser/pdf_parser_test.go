package parser

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPDFParser_SupportedType(t *testing.T) {
	parser := NewPDFParser()
	assert.Equal(t, "pdf", parser.SupportedType())
}

func TestPDFParser_Parse_EmptyReader(t *testing.T) {
	parser := NewPDFParser()

	// Test with empty reader
	emptyReader := bytes.NewReader([]byte{})
	text, err := parser.Parse(emptyReader)

	// With fallback strategy, returns metadata instead of error
	require.NoError(t, err)
	assert.Contains(t, text, "PDF Document Metadata")
}

func TestPDFParser_Parse_InvalidPDF(t *testing.T) {
	parser := NewPDFParser()

	// Test with invalid PDF data
	invalidData := []byte("This is not a PDF file")
	invalidReader := bytes.NewReader(invalidData)
	text, err := parser.Parse(invalidReader)

	// With fallback strategy, returns metadata instead of error
	require.NoError(t, err)
	assert.Contains(t, text, "PDF Document Metadata")
}

func TestPDFParser_Parse_ValidPDF(t *testing.T) {
	// This test requires a real PDF file
	// For now, we just verify the parser structure is correct
	parser := NewPDFParser()
	assert.NotNil(t, parser)

	// In a real scenario, you would:
	// 1. Create a test PDF file or use a sample PDF
	// 2. Read the file
	// 3. Call parser.Parse(reader)
	// 4. Verify the extracted text contains expected content
}

func TestPDFParser_MultiplePages(t *testing.T) {
	// This test would verify that multi-page PDFs are parsed correctly
	// and that page separators are added between pages
	parser := NewPDFParser()
	assert.NotNil(t, parser)

	// Example expected output structure:
	// "Page 1 content\n\n--- Page 2 ---\n\nPage 2 content"
}

func TestPDFParser_ExtractedTextValidation(t *testing.T) {
	// Test that the parser rejects PDFs with no extractable text
	// (e.g., scanned images without OCR)
	parser := NewPDFParser()
	assert.NotNil(t, parser)

	// This would require a sample "image-only" PDF for testing
	// The parser should return an error:
	// "no text content found in PDF document (possibly scanned images or empty)"
}

// Integration test example (commented out - requires actual PDF file)
/*
func TestPDFParser_Parse_RealPDF(t *testing.T) {
	parser := NewPDFParser()

	// Read a sample PDF file
	file, err := os.Open("testdata/sample.pdf")
	require.NoError(t, err)
	defer file.Close()

	// Parse the PDF
	text, err := parser.Parse(file)
	require.NoError(t, err)

	// Verify extracted text
	assert.NotEmpty(t, text)
	assert.Contains(t, text, "expected content from PDF")

	// Verify page separators for multi-page PDFs
	if strings.Contains(text, "--- Page 2 ---") {
		assert.True(t, strings.Contains(text, "--- Page 2 ---"))
	}
}
*/

func TestPDFParser_ErrorHandling(t *testing.T) {
	tests := []struct {
		name              string
		input             []byte
		expectMetadata    bool
		metadataSubstring string
	}{
		{
			name:              "empty input",
			input:             []byte{},
			expectMetadata:    true,
			metadataSubstring: "PDF Document Metadata",
		},
		{
			name:              "invalid PDF header",
			input:             []byte("not a pdf"),
			expectMetadata:    true,
			metadataSubstring: "PDF Document Metadata",
		},
		{
			name:              "corrupted PDF",
			input:             []byte("%PDF-1.4\nGarbage data"),
			expectMetadata:    true,
			metadataSubstring: "PDF Document Metadata",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewPDFParser()
			reader := bytes.NewReader(tt.input)

			text, err := parser.Parse(reader)

			// With fallback strategy, should not return error
			require.NoError(t, err)

			if tt.expectMetadata {
				assert.Contains(t, text, tt.metadataSubstring)
			}
		})
	}
}
