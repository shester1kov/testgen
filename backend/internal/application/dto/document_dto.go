package dto

// DocumentUploadResponse represents document upload response
type DocumentUploadResponse struct {
	ID         string  `json:"id"`
	UserID     string  `json:"user_id"`
	UserName   *string `json:"user_name,omitempty"`   // Only for admin
	UserEmail  *string `json:"user_email,omitempty"`  // Only for admin
	Title      string  `json:"title"`
	FileName   string  `json:"file_name"`
	FileType   string  `json:"file_type"`
	FileSize   int64   `json:"file_size"`
	ParsedText *string `json:"parsed_text,omitempty"` // Pointer to omit if null
	Status     string  `json:"status"`
	ErrorMsg   *string `json:"error_msg,omitempty"` // Pointer to omit if null
	CreatedAt  string  `json:"created_at"`
}

// DocumentListResponse represents list of documents
type DocumentListResponse struct {
	Documents []DocumentUploadResponse `json:"documents"`
	Total     int64                    `json:"total"`
	Page      int                      `json:"page"`
	PageSize  int                      `json:"page_size"`
}

// ParseDocumentRequest represents document parsing request
type ParseDocumentRequest struct {
	DocumentID string `json:"document_id" validate:"required,uuid"`
}

// ParseDocumentResponse represents document parsing response
type ParseDocumentResponse struct {
	ID          string `json:"id"`
	ParsedText  string `json:"parsed_text"`
	Status      string `json:"status"`
	TextPreview string `json:"text_preview"`
}
