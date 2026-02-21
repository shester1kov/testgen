package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDocument_IsParsed(t *testing.T) {
	doc := &Document{
		Status: StatusParsed,
	}
	assert.True(t, doc.IsParsed())

	doc.Status = StatusUploaded
	assert.False(t, doc.IsParsed())
}

func TestDocument_MarkAsParsing(t *testing.T) {
	doc := &Document{
		Status: StatusUploaded,
	}
	doc.MarkAsParsing()
	assert.Equal(t, StatusParsing, doc.Status)
}

func TestDocument_MarkAsParsed(t *testing.T) {
	doc := &Document{
		Status: StatusParsing,
	}
	parsedText := "This is parsed text"
	doc.MarkAsParsed(parsedText)

	assert.Equal(t, StatusParsed, doc.Status)
	assert.Equal(t, parsedText, doc.ParsedText)
}

func TestDocument_MarkAsError(t *testing.T) {
	doc := &Document{
		Status: StatusParsing,
	}
	errorMsg := "test error message"
	doc.MarkAsError(errorMsg)
	assert.Equal(t, StatusError, doc.Status)
	assert.Equal(t, errorMsg, doc.ErrorMsg)
}

func TestDocument_TableName(t *testing.T) {
	doc := Document{}
	assert.Equal(t, "documents", doc.TableName())
}

func TestDocument_GetContentType(t *testing.T) {
	t.Run("Returns content type if set", func(t *testing.T) {
		doc := &Document{
			FileType:    FileTypePDF,
			ContentType: "custom/type",
		}
		assert.Equal(t, "custom/type", doc.GetContentType())
	})

	t.Run("Determines content type by file type if not set", func(t *testing.T) {
		tests := []struct {
			fileType   FileType
			expectedCT string
		}{
			{FileTypePDF, "application/pdf"},
			{FileTypeDOCX, "application/vnd.openxmlformats-officedocument.wordprocessingml.document"},
			{FileTypePPTX, "application/vnd.openxmlformats-officedocument.presentationml.presentation"},
			{FileTypeTXT, "text/plain"},
			{FileTypeMD, "text/markdown"},
		}

		for _, tt := range tests {
			doc := &Document{FileType: tt.fileType}
			assert.Equal(t, tt.expectedCT, doc.GetContentType())
		}
	})
}

func TestGetContentTypeByExtension(t *testing.T) {
	tests := []struct {
		name     string
		fileType FileType
		expected string
	}{
		{"PDF extension", FileTypePDF, "application/pdf"},
		{"DOCX extension", FileTypeDOCX, "application/vnd.openxmlformats-officedocument.wordprocessingml.document"},
		{"PPTX extension", FileTypePPTX, "application/vnd.openxmlformats-officedocument.presentationml.presentation"},
		{"TXT extension", FileTypeTXT, "text/plain"},
		{"MD extension", FileTypeMD, "text/markdown"},
		{"Unknown extension", FileType("unknown"), "application/octet-stream"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetContentTypeByExtension(tt.fileType)
			assert.Equal(t, tt.expected, result)
		})
	}
}
