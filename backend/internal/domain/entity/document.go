package entity

import (
	"time"

	"github.com/google/uuid"
)

type FileType string

const (
	FileTypePDF  FileType = "pdf"
	FileTypeDOCX FileType = "docx"
	FileTypePPTX FileType = "pptx"
	FileTypeTXT  FileType = "txt"
	FileTypeMD   FileType = "md"
)

type DocumentStatus string

const (
	StatusUploaded DocumentStatus = "uploaded"
	StatusParsing  DocumentStatus = "parsing"
	StatusParsed   DocumentStatus = "parsed"
	StatusError    DocumentStatus = "error"
)

type Document struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID   uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	Title    string    `json:"title" gorm:"type:varchar(500);not null"`
	FileName string    `json:"file_name" gorm:"type:varchar(500);not null"`

	// ObjectKey - MinIO object key (путь к объекту в bucket)
	// Формат: "documents/{user_id}/{uuid}.{ext}"
	// Пример: "documents/550e8400-e29b-41d4-a716-446655440000/file.pdf"
	ObjectKey string `json:"object_key" gorm:"column:object_key;type:varchar(1000);not null"`

	FileType FileType `json:"file_type" gorm:"type:varchar(50);not null"`
	FileSize int64    `json:"file_size" gorm:"not null"`

	// ContentType - MIME-тип файла
	// Используется для установки Content-Type заголовка при скачивании
	ContentType string         `json:"content_type" gorm:"type:varchar(100)"`
	ParsedText  string         `json:"parsed_text,omitempty" gorm:"type:text"`
	Status      DocumentStatus `json:"status" gorm:"type:varchar(50);default:'uploaded';index"`
	ErrorMsg    string         `json:"error_msg,omitempty" gorm:"type:text"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   *time.Time     `json:"deleted_at,omitempty" gorm:"index"`

	// Relations
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for GORM
func (Document) TableName() string {
	return "documents"
}

// IsParsed checks if document is successfully parsed
func (d *Document) IsParsed() bool {
	return d.Status == StatusParsed
}

// MarkAsParsing sets document status to parsing
func (d *Document) MarkAsParsing() {
	d.Status = StatusParsing
}

// MarkAsParsed sets document status to parsed
func (d *Document) MarkAsParsed(parsedText string) {
	d.Status = StatusParsed
	d.ParsedText = parsedText
}

// MarkAsError sets document status to error
func (d *Document) MarkAsError(errMsg string) {
	d.Status = StatusError
	d.ErrorMsg = errMsg
}

// IsValidType checks if the document type is supported
func (d *Document) IsValidType() bool {
	switch d.FileType {
	case FileTypePDF, FileTypeDOCX, FileTypePPTX, FileTypeTXT, FileTypeMD:
		return true
	default:
		return false
	}

}

// GetContentType возвращает MIME-тип документа
// Если content_type заполнен - возвращает его
// Иначе определяет по file_type
func (d *Document) GetContentType() string {
	if d.ContentType != "" {
		return d.ContentType
	}

	return GetContentTypeByExtension(d.FileType)
}

// GetContentType возвращает MIME-тип на основе file_type
// Используется если content_type не заполнен
func GetContentTypeByExtension(ext FileType) string {

	switch ext {
	case FileTypePDF:
		return "application/pdf"
	case FileTypeDOCX:
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case FileTypePPTX:
		return "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	case FileTypeTXT:
		return "text/plain"
	case FileTypeMD:
		return "text/markdown"
	default:
		return "application/octet-stream"
	}
}
