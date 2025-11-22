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
