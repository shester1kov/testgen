package parser

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDOCXParser_SupportedType(t *testing.T) {
	parser := NewDOCXParser()
	assert.Equal(t, "docx", parser.SupportedType())
}

func TestDOCXParser_Parse_EmptyReader(t *testing.T) {
	parser := NewDOCXParser()

	// Test with empty reader
	emptyReader := bytes.NewReader([]byte{})
	_, err := parser.Parse(emptyReader)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "DOCX file is empty")
}

func TestDOCXParser_Parse_InvalidDOCX(t *testing.T) {
	parser := NewDOCXParser()

	// Test with invalid DOCX data
	invalidData := []byte("This is not a DOCX file")
	invalidReader := bytes.NewReader(invalidData)
	_, err := parser.Parse(invalidReader)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to open DOCX document")
}

func TestDOCXParser_ParseFromBytes(t *testing.T) {
	parser := NewDOCXParser()

	// Test the helper method with empty data
	_, err := parser.ParseFromBytes([]byte{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "DOCX file is empty")
}

// createMinimalDOCX creates a minimal valid DOCX file for testing
func createMinimalDOCX(content string) ([]byte, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Create [Content_Types].xml
	contentTypes := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
	<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
	<Default Extension="xml" ContentType="application/xml"/>
	<Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
</Types>`

	w, err := zipWriter.Create("[Content_Types].xml")
	if err != nil {
		return nil, err
	}
	if _, err := w.Write([]byte(contentTypes)); err != nil {
		return nil, err
	}

	// Create _rels/.rels
	rels := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
	<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>
</Relationships>`

	w, err = zipWriter.Create("_rels/.rels")
	if err != nil {
		return nil, err
	}
	if _, err := w.Write([]byte(rels)); err != nil {
		return nil, err
	}

	// Create word/_rels/document.xml.rels (required by the library)
	docRels := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
</Relationships>`

	w, err = zipWriter.Create("word/_rels/document.xml.rels")
	if err != nil {
		return nil, err
	}
	if _, err := w.Write([]byte(docRels)); err != nil {
		return nil, err
	}

	// Create word/document.xml with the content
	documentXML := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
	<w:body>
		<w:p>
			<w:r>
				<w:t>` + content + `</w:t>
			</w:r>
		</w:p>
	</w:body>
</w:document>`

	w, err = zipWriter.Create("word/document.xml")
	if err != nil {
		return nil, err
	}
	if _, err := w.Write([]byte(documentXML)); err != nil {
		return nil, err
	}

	if err := zipWriter.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func TestDOCXParser_Parse_ValidDOCX(t *testing.T) {
	parser := NewDOCXParser()

	// Create a minimal valid DOCX
	testContent := "Hello, World! This is a test document."
	docxData, err := createMinimalDOCX(testContent)
	require.NoError(t, err)

	// Parse the DOCX
	text, err := parser.Parse(bytes.NewReader(docxData))
	require.NoError(t, err)
	assert.NotEmpty(t, text)
	assert.Contains(t, text, "Hello, World")
}

func TestDOCXParser_Parse_EmptyContent(t *testing.T) {
	parser := NewDOCXParser()

	// Create a DOCX with empty content
	docxData, err := createMinimalDOCX("")
	require.NoError(t, err)

	// Parse the DOCX - should return error for empty content
	_, err = parser.Parse(bytes.NewReader(docxData))

	// Now that we extract text properly, empty DOCX should return error
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no text content found")
}

func TestDOCXParser_Parse_UnicodeContent(t *testing.T) {
	parser := NewDOCXParser()

	// Create a DOCX with unicode content
	unicodeContent := "Привет мир! 你好世界! مرحبا بالعالم!"
	docxData, err := createMinimalDOCX(unicodeContent)
	require.NoError(t, err)

	// Parse the DOCX
	text, err := parser.Parse(bytes.NewReader(docxData))
	require.NoError(t, err)
	assert.NotEmpty(t, text)
	// Note: The actual content may vary depending on the library's parsing
}

func TestDOCXParser_Parse_SpecialCharacters(t *testing.T) {
	parser := NewDOCXParser()

	// Create a DOCX with special characters
	specialContent := "Special chars: <>&\"' \t\n"
	docxData, err := createMinimalDOCX(specialContent)
	require.NoError(t, err)

	// Parse the DOCX
	text, err := parser.Parse(bytes.NewReader(docxData))
	require.NoError(t, err)
	assert.NotEmpty(t, text)
}

func TestDOCXParser_ErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		expectError bool
		errorMsg    string
	}{
		{
			name:        "empty input",
			input:       []byte{},
			expectError: true,
			errorMsg:    "DOCX file is empty",
		},
		{
			name:        "invalid DOCX header",
			input:       []byte("not a docx"),
			expectError: true,
			errorMsg:    "failed to open DOCX document",
		},
		{
			name:        "corrupted zip",
			input:       []byte("PK\x03\x04corrupted"),
			expectError: true,
			errorMsg:    "failed to open DOCX document",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewDOCXParser()
			reader := bytes.NewReader(tt.input)

			_, err := parser.Parse(reader)

			if tt.expectError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDOCXParser_Parse_RealFile(t *testing.T) {
	// This test creates a real DOCX file and tests parsing
	parser := NewDOCXParser()

	// Create a test DOCX file
	testContent := "This is a test document created for unit testing."
	docxData, err := createMinimalDOCX(testContent)
	require.NoError(t, err)

	// Save to temporary file for verification
	tmpFile, err := os.CreateTemp("", "test-*.docx")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write(docxData)
	require.NoError(t, err)
	tmpFile.Close()

	// Read and parse the file
	fileData, err := os.ReadFile(tmpFile.Name())
	require.NoError(t, err)

	text, err := parser.Parse(bytes.NewReader(fileData))
	require.NoError(t, err)
	assert.NotEmpty(t, text)
}

func TestDOCXParser_Parse_LargeContent(t *testing.T) {
	parser := NewDOCXParser()

	// Create a DOCX with large content (using printable characters)
	largeContent := ""
	for i := 0; i < 1000; i++ {
		largeContent += fmt.Sprintf("This is line %d. ", i)
	}

	docxData, err := createMinimalDOCX(largeContent)
	require.NoError(t, err)

	// Parse the DOCX
	text, err := parser.Parse(bytes.NewReader(docxData))
	require.NoError(t, err)
	assert.NotEmpty(t, text)
	assert.Contains(t, text, "This is line")
}

// Benchmark test
func BenchmarkDOCXParser_Parse(b *testing.B) {
	parser := NewDOCXParser()
	testContent := "Benchmark test content for DOCX parser"
	docxData, err := createMinimalDOCX(testContent)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.Parse(bytes.NewReader(docxData))
		if err != nil {
			b.Fatal(err)
		}
	}
}
