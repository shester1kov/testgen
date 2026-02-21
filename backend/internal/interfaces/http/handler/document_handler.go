package handler

import (
	"github.com/shester1kov/testgen-backend/internal/domain"

	"errors"

	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/application/dto"
	"github.com/shester1kov/testgen-backend/internal/application/usecase/document"
	"github.com/shester1kov/testgen-backend/pkg/security"
)

// DocumentHandler handles document operations
type DocumentHandler struct {
	uploadUseCase *document.UploadUseCase
	deleteUseCase *document.DeleteUseCase
	parseUseCase  *document.ParseUseCase
	listUseCase   *document.ListUseCase
	getUseCase    *document.GetUseCase
}

// NewDocumentHandler creates a new document handler
func NewDocumentHandler(
	uploadUseCase *document.UploadUseCase,
	deleteUseCase *document.DeleteUseCase,
	parseUseCase *document.ParseUseCase,
	listUseCase *document.ListUseCase,
	getUseCase *document.GetUseCase,
) *DocumentHandler {

	return &DocumentHandler{
		uploadUseCase: uploadUseCase,
		deleteUseCase: deleteUseCase,
		parseUseCase:  parseUseCase,
		listUseCase:   listUseCase,
		getUseCase:    getUseCase,
	}
}

// Upload godoc
// @Summary Upload a document
// @Description Upload a document file (PDF, DOCX, PPTX, TXT, MD) for processing
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
	// 1. Извлечь userID из контекста (JWT middleware)
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(
			dto.NewErrorResponse(dto.ErrCodeUnauthorized, "Unauthorized"),
		)
	}

	// 2. Получить файл из multipart form
	// Parse multipart form
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeInvalidInput, "no file provided"),
		)
	}

	// 3. Открыть файл для чтения (io.Reader для MinIO)
	// Fiber возвращает multipart.FileHeader, нужно открыть его
	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeInternalError, "failed to open uploaded file"),
		)
	}
	defer file.Close()

	// 4. Извлечь расширение файла (например, ".pdf" -> "pdf")
	// Get file extension
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(fileHeader.Filename), "."))
	if ext == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeInvalidFileType, "file has no extension"),
		)
	}

	// 5. Получить title из формы (или использовать имя файла)
	title := c.FormValue("title")
	if title == "" {
		title = fileHeader.Filename
	}

	// Sanitize title to prevent XSS
	sanitizedTitle := security.SanitizeInput(title)

	// 6. Подготовить параметры для use case
	params := document.UploadParams{
		UserID:      userID,
		Title:       sanitizedTitle,
		FileName:    fileHeader.Filename,
		FileType:    ext,
		FileSize:    fileHeader.Size,
		FileContent: file,
	}

	// 7. Вызвать use case (вся бизнес-логика внутри)
	response, err := h.uploadUseCase.Execute(c.Context(), params)
	if err != nil {
		// Use case вернет ошибку с описанием (размер, тип, БД, MinIO)
		// Определяем тип ошибки для правильного HTTP статуса
		if errors.Is(err, domain.ErrFileSizeExceeds) {
			return c.Status(fiber.StatusBadRequest).JSON(
				dto.NewErrorResponse(dto.ErrCodeFileTooLarge, err.Error()),
			)
		}

		if errors.Is(err, domain.ErrUnsupportedFileType) {
			return c.Status(fiber.StatusBadRequest).JSON(
				dto.NewErrorResponse(dto.ErrCodeInvalidFileType, err.Error()),
			)

		}

		// Все остальные ошибки (MinIO, БД) - 500
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeInternalError, "failed to upload document"),
		)
	}

	// 8. Преобразовать entity.Document в DTO и вернуть
	// Create document record
	return c.Status(fiber.StatusCreated).JSON(response)
}

// List godoc
// @Summary List documents
// @Description Get paginated list of documents. Admin sees all documents with user info, others see only their own
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
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(
			dto.NewErrorResponse(dto.ErrCodeUnauthorized, "Unauthorized"),
		)
	}

	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 20)

	response, err := h.listUseCase.Execute(c.Context(), userID, page, pageSize)

	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return c.Status(fiber.StatusInternalServerError).JSON(
				dto.NewErrorResponse(dto.ErrCodeDatabaseError, "failed to fetch user"),
			)
		}

		if errors.Is(err, domain.ErrDocumentNotFound) {
			return c.Status(fiber.StatusInternalServerError).JSON(
				dto.NewErrorResponse(dto.ErrCodeDatabaseError, "failed to fetch documents"),
			)
		}

		if errors.Is(err, domain.ErrDocumentNotFound) {
			return c.Status(fiber.StatusInternalServerError).JSON(
				dto.NewErrorResponse(dto.ErrCodeDatabaseError, "failed to count documents"),
			)
		}

		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeDatabaseError, "internal server error"),
		)

	}

	return c.JSON(response)
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
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(
			dto.NewErrorResponse(dto.ErrCodeUnauthorized, "Unauthorized"),
		)
	}
	documentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeInvalidUUID, "invalid document ID"),
		)
	}

	// Get document WITHOUT parsed text (avoid large response)
	response, err := h.getUseCase.Execute(c.Context(), documentID, userID, false)
	if err != nil {
		if errors.Is(err, domain.ErrAccessDenied) {
			return c.Status(fiber.StatusForbidden).JSON(
				dto.NewErrorResponse(dto.ErrCodeForbidden, "access denied"),
			)
		}

		if errors.Is(err, domain.ErrDocumentNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(
				dto.NewErrorResponse(dto.ErrCodeDocumentNotFound, "document not found"),
			)
		}

		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeInternalError, "internal server error"),
		)
	}

	return c.JSON(response)
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
	// 1. Извлечь userI
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(
			dto.NewErrorResponse(dto.ErrCodeUnauthorized, "Unauthorized"),
		)
	}

	// 2. Извлечь documentID из URL параметра
	documentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeInvalidUUID, "invalid document ID"),
		)
	}

	// 3. Вызвать use case (вся логика внутри)
	if err := h.deleteUseCase.Execute(c.Context(), documentID, userID); err != nil {
		// Определяем тип ошибки
		if errors.Is(err, domain.ErrDocumentNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(
				dto.NewErrorResponse(dto.ErrCodeDocumentNotFound, "document not found"),
			)
		}

		if errors.Is(err, domain.ErrAccessDenied) {
			return c.Status(fiber.StatusForbidden).JSON(
				dto.NewErrorResponse(dto.ErrCodeForbidden, "access denied"),
			)
		}

		// Ошибки MinIO или БД
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeInternalError, "failed to delete document"),
		)

	}

	// 4. Успешный ответ
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
	// 1. Извлечь userID
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(
			dto.NewErrorResponse(dto.ErrCodeUnauthorized, "Unauthorized"),
		)
	}

	// 2. Извлечь documentID
	documentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeInvalidUUID, "invalid document ID"),
		)
	}

	// 3. Вызвать use case (вся логика внутри)
	if err := h.parseUseCase.Execute(c.Context(), documentID, userID); err != nil {
		// Определяем тип ошибки

		if errors.Is(err, domain.ErrDocumentNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(
				dto.NewErrorResponse(dto.ErrCodeDocumentNotFound, "document not found"),
			)
		}

		if errors.Is(err, domain.ErrUnauthorized) {
			return c.Status(fiber.StatusForbidden).JSON(
				dto.NewErrorResponse(dto.ErrCodeForbidden, "access denied"),
			)
		}

		if errors.Is(err, domain.ErrAlreadyParsed) {
			return c.Status(fiber.StatusBadRequest).JSON(
				dto.NewErrorResponse(dto.ErrCodeInvalidInput, "document already parsed"),
			)
		}

		// Ошибки парсинга
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeParsingFailed, "failed to parse document"),
		)
	}

	// 4. После успешного парсинга получить обновленный документ для ответа
	// Include parsed text for Parse endpoint (true)
	document, err := h.getUseCase.Execute(c.Context(), documentID, userID, true)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeDatabaseError, "failed to fetch parsed document"),
		)
	}

	// 5. Подготовить preview и полный текст (работаем с указателем!)
	preview := ""
	fullText := ""

	// Проверяем, что ParsedText не nil и распаковываем указатель
	if document.ParsedText != nil {
		fullText = *document.ParsedText // Распаковать *string → string
		preview = fullText

		// Обрезать preview до 500 символов
		if len(preview) > 500 {
			preview = preview[:500] + "..."
		}
	}

	// 6. Вернуть ответ с правильными типами
	return c.JSON(dto.ParseDocumentResponse{
		ID:          document.ID,     // Уже string
		ParsedText:  fullText,        // string (распакован из *string)
		Status:      document.Status, // Уже string
		TextPreview: preview,         // string
	})
}
