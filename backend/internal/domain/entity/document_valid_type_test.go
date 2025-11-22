package entity

import "testing"

func TestDocument_IsValidType(t *testing.T) {
	tests := []struct {
		name     string
		fileType FileType
		valid    bool
	}{
		{"pdf valid", FileTypePDF, true},
		{"docx valid", FileTypeDOCX, true},
		{"pptx valid", FileTypePPTX, true},
		{"txt valid", FileTypeTXT, true},
		{"md valid", FileTypeMD, true},
		{"empty invalid", FileType(""), false},
		{"unknown invalid", FileType("exe"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &Document{FileType: tt.fileType}
			if got := doc.IsValidType(); got != tt.valid {
				t.Fatalf("expected %v, got %v", tt.valid, got)
			}
		})
	}
}
