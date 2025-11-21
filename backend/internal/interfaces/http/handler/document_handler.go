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

// Upload godoc
// @Summary Upload a document
// @Description Upload a document file (PDF, DOCX, PPTX, TXT) for processing
// @Tags documents
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Document file"
// @Param title formData string false "Document title"
// @Success 201 {object} dto.DocumentUploadResponse
// @Failure 400 {object} dto.ErrorResponse "Invalid input or unsupported file type"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /documents [post]
func (h *DocumentHandler) Upload(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	// Parse multipart form
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeInvalidInput, "no file provided"),
		)
	}

	// Validate file size
	if file.Size > h.maxFileSize {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeFileTooLarge, fmt.Sprintf("file size exceeds maximum allowed size of %d bytes", h.maxFileSize)),
		)
	}

	// Get file extension
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(file.Filename), "."))
	if ext == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeInvalidFileType, "file has no extension"),
		)
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
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeInvalidFileType, fmt.Sprintf("unsupported file type. Supported types: %v", supportedTypes)),
		)
	}

	// Generate unique filename
	uniqueID := uuid.New()
	savedFilename := fmt.Sprintf("%s%s", uniqueID.String(), filepath.Ext(file.Filename))
	savedPath := filepath.Join(h.uploadDir, savedFilename)

	// Save file
	if err := c.SaveFile(file, savedPath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeInternalError, "failed to save file"),
		)
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
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeDatabaseError, "failed to create document record"),
		)
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

// List godoc
// @Summary List user's documents
// @Description Get paginated list of documents uploaded by the current user
// @Tags documents
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} dto.DocumentListResponse
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /documents [get]
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
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeDatabaseError, "failed to fetch documents"),
		)
	}

	total, err := h.documentRepo.CountByUserID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeDatabaseError, "failed to count documents"),
		)
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

// GetByID godoc
// @Summary Get document by ID
// @Description Get details of a specific document by its ID
// @Tags documents
// @Produce json
// @Security BearerAuth
// @Param id path string true "Document ID"
// @Success 200 {object} dto.DocumentUploadResponse
// @Failure 400 {object} dto.ErrorResponse "Invalid document ID"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Access denied"
// @Failure 404 {object} dto.ErrorResponse "Document not found"
// @Router /documents/{id} [get]
func (h *DocumentHandler) GetByID(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	documentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeInvalidUUID, "invalid document ID"),
		)
	}

	document, err := h.documentRepo.FindByID(c.Context(), documentID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			dto.NewErrorResponse(dto.ErrCodeDocumentNotFound, "document not found"),
		)
	}

	// Check ownership
	if document.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(
			dto.NewErrorResponse(dto.ErrCodeForbidden, "access denied"),
		)
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

// Delete godoc
// @Summary Delete a document
// @Description Delete a document and its associated file
// @Tags documents
// @Produce json
// @Security BearerAuth
// @Param id path string true "Document ID"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse "Invalid document ID"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Access denied"
// @Failure 404 {object} dto.ErrorResponse "Document not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /documents/{id} [delete]
func (h *DocumentHandler) Delete(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	documentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeInvalidUUID, "invalid document ID"),
		)
	}

	document, err := h.documentRepo.FindByID(c.Context(), documentID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			dto.NewErrorResponse(dto.ErrCodeDocumentNotFound, "document not found"),
		)
	}

	// Check ownership
	if document.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(
			dto.NewErrorResponse(dto.ErrCodeForbidden, "access denied"),
		)
	}

	// Delete file
	os.Remove(document.FilePath)

	// Delete from database
	if err := h.documentRepo.Delete(c.Context(), documentID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeDatabaseError, "failed to delete document"),
		)
	}

	return c.JSON(dto.NewMessageResponse("document deleted successfully"))
}

// Parse godoc
// @Summary Parse a document
// @Description Extract text content from an uploaded document
// @Tags documents
// @Produce json
// @Security BearerAuth
// @Param id path string true "Document ID"
// @Success 200 {object} dto.ParseDocumentResponse
// @Failure 400 {object} dto.ErrorResponse "Invalid document ID or file type"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Access denied"
// @Failure 404 {object} dto.ErrorResponse "Document not found"
// @Failure 500 {object} dto.ErrorResponse "Parsing failed"
// @Router /documents/{id}/parse [post]
func (h *DocumentHandler) Parse(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	documentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeInvalidUUID, "invalid document ID"),
		)
	}

	document, err := h.documentRepo.FindByID(c.Context(), documentID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			dto.NewErrorResponse(dto.ErrCodeDocumentNotFound, "document not found"),
		)
	}

	// Check ownership
	if document.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(
			dto.NewErrorResponse(dto.ErrCodeForbidden, "access denied"),
		)
	}

	// Get parser
	docParser, err := h.parserFactory.CreateParser(string(document.FileType))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeInvalidFileType, err.Error()),
		)
	}

	// Open file
	file, err := os.Open(document.FilePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeInternalError, "failed to open file"),
		)
	}
	defer file.Close()

	// Parse document
	document.MarkAsParsing()
	h.documentRepo.Update(c.Context(), document)

	parsedText, err := docParser.Parse(file)
	if err != nil {
		document.MarkAsError()
		h.documentRepo.Update(c.Context(), document)
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeParsingFailed, "failed to parse document"),
		)
	}

	// Update document with parsed text
	document.MarkAsParsed(parsedText)
	if err := h.documentRepo.Update(c.Context(), document); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeDatabaseError, "failed to save parsed text"),
		)
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
