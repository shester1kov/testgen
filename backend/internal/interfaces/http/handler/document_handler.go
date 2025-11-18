package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/application/dto"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/parser"
)

// DocumentHandler handles document operations
type DocumentHandler struct {
	documentRepo  repository.DocumentRepository
	parserFactory *parser.DocumentParserFactory
	uploadDir     string
	maxFileSize   int64
}

// NewDocumentHandler creates a new document handler
func NewDocumentHandler(
	documentRepo repository.DocumentRepository,
	parserFactory *parser.DocumentParserFactory,
	uploadDir string,
	maxFileSize int64,
) *DocumentHandler {
	// Ensure upload directory exists
	os.MkdirAll(uploadDir, 0755)

	return &DocumentHandler{
		documentRepo:  documentRepo,
		parserFactory: parserFactory,
		uploadDir:     uploadDir,
		maxFileSize:   maxFileSize,
	}
}

// Upload handles document upload
func (h *DocumentHandler) Upload(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	// Parse multipart form
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "no file provided",
		})
	}

	// Validate file size
	if file.Size > h.maxFileSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("file size exceeds maximum allowed size of %d bytes", h.maxFileSize),
		})
	}

	// Get file extension
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(file.Filename), "."))
	if ext == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "file has no extension",
		})
	}

	// Validate file type
	supportedTypes := h.parserFactory.GetSupportedTypes()
	isSupported := false
	for _, supportedType := range supportedTypes {
		if ext == supportedType {
			isSupported = true
			break
		}
	}

	if !isSupported {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("unsupported file type. Supported types: %v", supportedTypes),
		})
	}

	// Generate unique filename
	uniqueID := uuid.New()
	savedFilename := fmt.Sprintf("%s%s", uniqueID.String(), filepath.Ext(file.Filename))
	savedPath := filepath.Join(h.uploadDir, savedFilename)

	// Save file
	if err := c.SaveFile(file, savedPath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to save file",
		})
	}

	// Get title from form or use filename
	title := c.FormValue("title")
	if title == "" {
		title = file.Filename
	}

	// Create document record
	document := &entity.Document{
		UserID:   userID,
		Title:    title,
		FileName: file.Filename,
		FilePath: savedPath,
		FileType: entity.FileType(ext),
		FileSize: file.Size,
		Status:   entity.StatusUploaded,
	}

	if err := h.documentRepo.Create(c.Context(), document); err != nil {
		// Clean up file if database insert fails
		os.Remove(savedPath)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create document record",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.DocumentUploadResponse{
		ID:        document.ID.String(),
		Title:     document.Title,
		FileName:  document.FileName,
		FileType:  string(document.FileType),
		FileSize:  document.FileSize,
		Status:    string(document.Status),
		CreatedAt: document.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// List returns list of user's documents
func (h *DocumentHandler) List(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 20)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	documents, err := h.documentRepo.FindByUserID(c.Context(), userID, pageSize, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch documents",
		})
	}

	total, err := h.documentRepo.CountByUserID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to count documents",
		})
	}

	result := make([]dto.DocumentUploadResponse, len(documents))
	for i, doc := range documents {
		result[i] = dto.DocumentUploadResponse{
			ID:        doc.ID.String(),
			Title:     doc.Title,
			FileName:  doc.FileName,
			FileType:  string(doc.FileType),
			FileSize:  doc.FileSize,
			Status:    string(doc.Status),
			CreatedAt: doc.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return c.JSON(dto.DocumentListResponse{
		Documents: result,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
	})
}

// GetByID returns document by ID
func (h *DocumentHandler) GetByID(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	documentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid document ID",
		})
	}

	document, err := h.documentRepo.FindByID(c.Context(), documentID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "document not found",
		})
	}

	// Check ownership
	if document.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "access denied",
		})
	}

	return c.JSON(dto.DocumentUploadResponse{
		ID:        document.ID.String(),
		Title:     document.Title,
		FileName:  document.FileName,
		FileType:  string(document.FileType),
		FileSize:  document.FileSize,
		Status:    string(document.Status),
		CreatedAt: document.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// Delete removes a document
func (h *DocumentHandler) Delete(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	documentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid document ID",
		})
	}

	document, err := h.documentRepo.FindByID(c.Context(), documentID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "document not found",
		})
	}

	// Check ownership
	if document.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "access denied",
		})
	}

	// Delete file
	os.Remove(document.FilePath)

	// Delete from database
	if err := h.documentRepo.Delete(c.Context(), documentID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to delete document",
		})
	}

	return c.JSON(fiber.Map{
		"message": "document deleted successfully",
	})
}

// Parse parses document and extracts text
func (h *DocumentHandler) Parse(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	documentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid document ID",
		})
	}

	document, err := h.documentRepo.FindByID(c.Context(), documentID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "document not found",
		})
	}

	// Check ownership
	if document.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "access denied",
		})
	}

	// Get parser
	docParser, err := h.parserFactory.CreateParser(string(document.FileType))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Open file
	file, err := os.Open(document.FilePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to open file",
		})
	}
	defer file.Close()

	// Parse document
	document.MarkAsParsing()
	h.documentRepo.Update(c.Context(), document)

	parsedText, err := docParser.Parse(file)
	if err != nil {
		document.MarkAsError()
		h.documentRepo.Update(c.Context(), document)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to parse document",
		})
	}

	// Update document with parsed text
	document.MarkAsParsed(parsedText)
	if err := h.documentRepo.Update(c.Context(), document); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to save parsed text",
		})
	}

	// Return preview (first 500 chars)
	preview := parsedText
	if len(preview) > 500 {
		preview = preview[:500] + "..."
	}

	return c.JSON(dto.ParseDocumentResponse{
		ID:          document.ID.String(),
		ParsedText:  parsedText,
		Status:      string(document.Status),
		TextPreview: preview,
	})
}
