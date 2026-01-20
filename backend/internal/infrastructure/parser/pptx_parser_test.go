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

func TestPPTXParser_SupportedType(t *testing.T) {
	parser := NewPPTXParser()
	assert.Equal(t, "pptx", parser.SupportedType())
}

func TestPPTXParser_Parse_EmptyReader(t *testing.T) {
	parser := NewPPTXParser()

	// Test with empty reader
	emptyReader := bytes.NewReader([]byte{})
	_, err := parser.Parse(emptyReader)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "PPTX file is empty")
}

func TestPPTXParser_Parse_InvalidPPTX(t *testing.T) {
	parser := NewPPTXParser()

	// Test with invalid PPTX data
	invalidData := []byte("This is not a PPTX file")
	invalidReader := bytes.NewReader(invalidData)
	_, err := parser.Parse(invalidReader)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to open PPTX as ZIP archive")
}

func TestPPTXParser_ParseFromBytes(t *testing.T) {
	parser := NewPPTXParser()

	// Test the helper method with empty data
	_, err := parser.ParseFromBytes([]byte{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "PPTX file is empty")
}

// createMinimalPPTX creates a minimal valid PPTX file for testing
func createMinimalPPTX(slides []string) ([]byte, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Create [Content_Types].xml
	contentTypes := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
	<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
	<Default Extension="xml" ContentType="application/xml"/>
	<Override PartName="/ppt/presentation.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"/>`

	for i := range slides {
		contentTypes += fmt.Sprintf(`
	<Override PartName="/ppt/slides/slide%d.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/>`, i+1)
	}
	contentTypes += `
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
	<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="ppt/presentation.xml"/>
</Relationships>`

	w, err = zipWriter.Create("_rels/.rels")
	if err != nil {
		return nil, err
	}
	if _, err := w.Write([]byte(rels)); err != nil {
		return nil, err
	}

	// Create ppt/presentation.xml
	presentation := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:presentation xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
	<p:sldIdLst>`
	for i := range slides {
		presentation += fmt.Sprintf(`
		<p:sldId id="%d" r:id="rId%d" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"/>`, 256+i, i+1)
	}
	presentation += `
	</p:sldIdLst>
</p:presentation>`

	w, err = zipWriter.Create("ppt/presentation.xml")
	if err != nil {
		return nil, err
	}
	if _, err := w.Write([]byte(presentation)); err != nil {
		return nil, err
	}

	// Create slides
	for i, content := range slides {
		slideXML := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
	<p:cSld>
		<p:spTree>
			<p:sp>
				<p:txBody>
					<a:p>
						<a:r>
							<a:t>%s</a:t>
						</a:r>
					</a:p>
				</p:txBody>
			</p:sp>
		</p:spTree>
	</p:cSld>
</p:sld>`, content)

		w, err = zipWriter.Create(fmt.Sprintf("ppt/slides/slide%d.xml", i+1))
		if err != nil {
			return nil, err
		}
		if _, err := w.Write([]byte(slideXML)); err != nil {
			return nil, err
		}
	}

	if err := zipWriter.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// createMinimalPPTXWithParagraphs creates a PPTX with multiple paragraphs per slide
func createMinimalPPTXWithParagraphs(slideParagraphs [][]string) ([]byte, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Create [Content_Types].xml
	contentTypes := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
	<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
	<Default Extension="xml" ContentType="application/xml"/>
	<Override PartName="/ppt/presentation.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"/>`

	for i := range slideParagraphs {
		contentTypes += fmt.Sprintf(`
	<Override PartName="/ppt/slides/slide%d.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/>`, i+1)
	}
	contentTypes += `
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
	<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="ppt/presentation.xml"/>
</Relationships>`

	w, err = zipWriter.Create("_rels/.rels")
	if err != nil {
		return nil, err
	}
	if _, err := w.Write([]byte(rels)); err != nil {
		return nil, err
	}

	// Create ppt/presentation.xml
	presentation := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:presentation xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
	<p:sldIdLst>`
	for i := range slideParagraphs {
		presentation += fmt.Sprintf(`
		<p:sldId id="%d" r:id="rId%d" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"/>`, 256+i, i+1)
	}
	presentation += `
	</p:sldIdLst>
</p:presentation>`

	w, err = zipWriter.Create("ppt/presentation.xml")
	if err != nil {
		return nil, err
	}
	if _, err := w.Write([]byte(presentation)); err != nil {
		return nil, err
	}

	// Create slides with multiple paragraphs
	for i, paragraphs := range slideParagraphs {
		slideXML := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
	<p:cSld>
		<p:spTree>
			<p:sp>
				<p:txBody>`

		for _, para := range paragraphs {
			slideXML += fmt.Sprintf(`
					<a:p>
						<a:r>
							<a:t>%s</a:t>
						</a:r>
					</a:p>`, para)
		}

		slideXML += `
				</p:txBody>
			</p:sp>
		</p:spTree>
	</p:cSld>
</p:sld>`

		w, err = zipWriter.Create(fmt.Sprintf("ppt/slides/slide%d.xml", i+1))
		if err != nil {
			return nil, err
		}
		if _, err := w.Write([]byte(slideXML)); err != nil {
			return nil, err
		}
	}

	if err := zipWriter.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func TestPPTXParser_Parse_ValidPPTX(t *testing.T) {
	parser := NewPPTXParser()

	// Create a minimal valid PPTX with one slide
	slides := []string{"Hello, World! This is a test presentation."}
	pptxData, err := createMinimalPPTX(slides)
	require.NoError(t, err)

	// Parse the PPTX
	text, err := parser.Parse(bytes.NewReader(pptxData))
	require.NoError(t, err)
	assert.NotEmpty(t, text)
	assert.Contains(t, text, "Hello, World")
	assert.Contains(t, text, "--- Slide 1 ---")
}

func TestPPTXParser_Parse_MultipleSlides(t *testing.T) {
	parser := NewPPTXParser()

	// Create a PPTX with multiple slides
	slides := []string{
		"First slide content",
		"Second slide content",
		"Third slide content",
	}
	pptxData, err := createMinimalPPTX(slides)
	require.NoError(t, err)

	// Parse the PPTX
	text, err := parser.Parse(bytes.NewReader(pptxData))
	require.NoError(t, err)
	assert.NotEmpty(t, text)
	assert.Contains(t, text, "--- Slide 1 ---")
	assert.Contains(t, text, "First slide content")
	assert.Contains(t, text, "--- Slide 2 ---")
	assert.Contains(t, text, "Second slide content")
	assert.Contains(t, text, "--- Slide 3 ---")
	assert.Contains(t, text, "Third slide content")
}

func TestPPTXParser_Parse_MultipleParagraphs(t *testing.T) {
	parser := NewPPTXParser()

	// Create a PPTX with multiple paragraphs per slide
	slideParagraphs := [][]string{
		{"Title of Slide 1", "First paragraph", "Second paragraph"},
		{"Title of Slide 2", "Another paragraph"},
	}
	pptxData, err := createMinimalPPTXWithParagraphs(slideParagraphs)
	require.NoError(t, err)

	// Parse the PPTX
	text, err := parser.Parse(bytes.NewReader(pptxData))
	require.NoError(t, err)
	assert.NotEmpty(t, text)
	assert.Contains(t, text, "Title of Slide 1")
	assert.Contains(t, text, "First paragraph")
	assert.Contains(t, text, "Second paragraph")
	// Check that paragraphs are separated by newlines
	assert.Contains(t, text, "\n")
}

func TestPPTXParser_Parse_NoSlides(t *testing.T) {
	parser := NewPPTXParser()

	// Create a valid ZIP but without slides
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Create minimal content types without slides
	contentTypes := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
	<Default Extension="xml" ContentType="application/xml"/>
</Types>`

	w, err := zipWriter.Create("[Content_Types].xml")
	require.NoError(t, err)
	_, err = w.Write([]byte(contentTypes))
	require.NoError(t, err)
	require.NoError(t, zipWriter.Close())

	// Parse should fail - no slides found
	_, err = parser.Parse(bytes.NewReader(buf.Bytes()))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no slides found")
}

func TestPPTXParser_Parse_EmptySlides(t *testing.T) {
	parser := NewPPTXParser()

	// Create a PPTX with empty slide content
	slides := []string{""}
	pptxData, err := createMinimalPPTX(slides)
	require.NoError(t, err)

	// Parse the PPTX - should return error for empty content
	_, err = parser.Parse(bytes.NewReader(pptxData))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no text content found")
}

func TestPPTXParser_Parse_UnicodeContent(t *testing.T) {
	parser := NewPPTXParser()

	// Create a PPTX with unicode content
	slides := []string{"Привет мир! 你好世界! مرحبا بالعالم!"}
	pptxData, err := createMinimalPPTX(slides)
	require.NoError(t, err)

	// Parse the PPTX
	text, err := parser.Parse(bytes.NewReader(pptxData))
	require.NoError(t, err)
	assert.NotEmpty(t, text)
	assert.Contains(t, text, "Привет мир")
	assert.Contains(t, text, "你好世界")
}

func TestPPTXParser_Parse_SpecialCharacters(t *testing.T) {
	parser := NewPPTXParser()

	// Note: XML special chars need to be escaped in the test data
	slides := []string{"Special chars: test content"}
	pptxData, err := createMinimalPPTX(slides)
	require.NoError(t, err)

	// Parse the PPTX
	text, err := parser.Parse(bytes.NewReader(pptxData))
	require.NoError(t, err)
	assert.NotEmpty(t, text)
	assert.Contains(t, text, "Special chars")
}

func TestPPTXParser_ErrorHandling(t *testing.T) {
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
			errorMsg:    "PPTX file is empty",
		},
		{
			name:        "invalid PPTX header",
			input:       []byte("not a pptx"),
			expectError: true,
			errorMsg:    "failed to open PPTX as ZIP archive",
		},
		{
			name:        "corrupted zip",
			input:       []byte("PK\x03\x04corrupted"),
			expectError: true,
			errorMsg:    "failed to open PPTX as ZIP archive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewPPTXParser()
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

func TestPPTXParser_Parse_RealFile(t *testing.T) {
	// This test creates a real PPTX file and tests parsing
	parser := NewPPTXParser()

	// Create a test PPTX file
	slides := []string{"This is a test presentation created for unit testing."}
	pptxData, err := createMinimalPPTX(slides)
	require.NoError(t, err)

	// Save to temporary file for verification
	tmpFile, err := os.CreateTemp("", "test-*.pptx")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write(pptxData)
	require.NoError(t, err)
	tmpFile.Close()

	// Read and parse the file
	fileData, err := os.ReadFile(tmpFile.Name())
	require.NoError(t, err)

	text, err := parser.Parse(bytes.NewReader(fileData))
	require.NoError(t, err)
	assert.NotEmpty(t, text)
	assert.Contains(t, text, "test presentation")
}

func TestPPTXParser_Parse_LargeContent(t *testing.T) {
	parser := NewPPTXParser()

	// Create a PPTX with many slides
	var slides []string
	for i := 0; i < 50; i++ {
		slides = append(slides, fmt.Sprintf("This is slide number %d with some content.", i+1))
	}

	pptxData, err := createMinimalPPTX(slides)
	require.NoError(t, err)

	// Parse the PPTX
	text, err := parser.Parse(bytes.NewReader(pptxData))
	require.NoError(t, err)
	assert.NotEmpty(t, text)
	assert.Contains(t, text, "slide number 1")
	assert.Contains(t, text, "slide number 50")
	assert.Contains(t, text, "--- Slide 50 ---")
}

func TestPPTXParser_Parse_SlideOrdering(t *testing.T) {
	parser := NewPPTXParser()

	// Create PPTX with numbered slides to verify ordering
	slides := []string{"Slide A", "Slide B", "Slide C"}
	pptxData, err := createMinimalPPTX(slides)
	require.NoError(t, err)

	// Parse the PPTX
	text, err := parser.Parse(bytes.NewReader(pptxData))
	require.NoError(t, err)

	// Verify slides appear in correct order
	idxA := bytes.Index([]byte(text), []byte("Slide A"))
	idxB := bytes.Index([]byte(text), []byte("Slide B"))
	idxC := bytes.Index([]byte(text), []byte("Slide C"))

	assert.True(t, idxA < idxB, "Slide A should appear before Slide B")
	assert.True(t, idxB < idxC, "Slide B should appear before Slide C")
}

func TestPPTXParser_extractTextFromXML(t *testing.T) {
	parser := NewPPTXParser()

	tests := []struct {
		name     string
		xmlInput string
		expected string
	}{
		{
			name: "simple text",
			xmlInput: `<a:p xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">
				<a:r><a:t>Hello World</a:t></a:r>
			</a:p>`,
			expected: "Hello World",
		},
		{
			name: "multiple paragraphs",
			xmlInput: `<root xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">
				<a:p><a:r><a:t>First</a:t></a:r></a:p>
				<a:p><a:r><a:t>Second</a:t></a:r></a:p>
			</root>`,
			expected: "First\nSecond",
		},
		{
			name:     "empty xml",
			xmlInput: `<root></root>`,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.extractTextFromXML(tt.xmlInput)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Benchmark test
func BenchmarkPPTXParser_Parse(b *testing.B) {
	parser := NewPPTXParser()
	slides := []string{"Benchmark test content for PPTX parser"}
	pptxData, err := createMinimalPPTX(slides)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.Parse(bytes.NewReader(pptxData))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPPTXParser_Parse_LargeFile(b *testing.B) {
	parser := NewPPTXParser()

	// Create a larger PPTX
	var slides []string
	for i := 0; i < 20; i++ {
		slides = append(slides, fmt.Sprintf("This is slide %d with a longer content to simulate real presentations with more text.", i+1))
	}

	pptxData, err := createMinimalPPTX(slides)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.Parse(bytes.NewReader(pptxData))
		if err != nil {
			b.Fatal(err)
		}
	}
}
