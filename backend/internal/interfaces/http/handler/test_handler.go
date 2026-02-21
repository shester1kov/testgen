package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/application/dto"
	testuc "github.com/shester1kov/testgen-backend/internal/application/usecase/test"
)

type TestHandler struct {
	createUseCase         *testuc.CreateUseCase
	generateUseCase       *testuc.GenerateUseCase
	listUseCase           *testuc.ListUseCase
	getUseCase            *testuc.GetUseCase
	updateUseCase         *testuc.UpdateUseCase
	deleteUseCase         *testuc.DeleteUseCase
	updateQuestionUseCase *testuc.UpdateQuestionUseCase
	exportUseCase         *testuc.ExportUseCase
}

func NewTestHandler(
	createUseCase *testuc.CreateUseCase,
	generateUseCase *testuc.GenerateUseCase,
	listUseCase *testuc.ListUseCase,
	getUseCase *testuc.GetUseCase,
	updateUseCase *testuc.UpdateUseCase,
	deleteUseCase *testuc.DeleteUseCase,
	updateQuestionUseCase *testuc.UpdateQuestionUseCase,
	exportUseCase *testuc.ExportUseCase,
) *TestHandler {
	return &TestHandler{
		createUseCase:         createUseCase,
		generateUseCase:       generateUseCase,
		listUseCase:           listUseCase,
		getUseCase:            getUseCase,
		updateUseCase:         updateUseCase,
		deleteUseCase:         deleteUseCase,
		updateQuestionUseCase: updateQuestionUseCase,
		exportUseCase:         exportUseCase,
	}
}

// Create godoc
// @Summary Create a new test
// @Description Create a new test with optional document association
// @Tags tests
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateTestRequest true "Create test request"
// @Success 201 {object} dto.TestResponse
// @Failure 400 {object} dto.ErrorResponse "Invalid input"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /tests [post]
func (h *TestHandler) Create(c *fiber.Ctx) error {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(
			dto.NewErrorResponse(dto.ErrCodeUnauthorized, "Unauthorized"),
		)
	}
	var req dto.CreateTestRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeInvalidInput, "invalid request"),
		)
	}

	result, err := h.createUseCase.Execute(c.Context(), testuc.CreateParams{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		DocumentID:  req.DocumentID,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeDatabaseError, "failed to create test"),
		)
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

// Generate godoc
// @Summary Generate test questions
// @Description Generate test questions from a document using LLM
// @Tags tests
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.GenerateTestRequest true "Generate test request"
// @Success 201 {object} dto.TestResponse "Test created with generated questions"
// @Failure 400 {object} dto.ErrorResponse "Invalid input or document not parsed"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 404 {object} dto.ErrorResponse "Document not found"
// @Failure 500 {object} dto.ErrorResponse "Generation failed or database error"
// @Router /tests/generate [post]
func (h *TestHandler) Generate(c *fiber.Ctx) error {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(
			dto.NewErrorResponse(dto.ErrCodeUnauthorized, "Unauthorized"),
		)
	}
	var req dto.GenerateTestRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeInvalidInput, "invalid request"),
		)
	}

	docID, _ := uuid.Parse(req.DocumentID)

	result, err := h.generateUseCase.Execute(c.Context(), testuc.GenerateParams{
		UserID:      userID,
		DocumentID:  docID,
		Title:       req.Title,
		NumQuestions: req.NumQuestions,
		Difficulty:  req.Difficulty,
		LLMProvider: req.LLMProvider,
	})
	if err != nil {
		if strings.Contains(err.Error(), "document not found") {
			return c.Status(fiber.StatusNotFound).JSON(
				dto.NewErrorResponse(dto.ErrCodeDocumentNotFound, "document not found"),
			)
		}
		if strings.Contains(err.Error(), "access denied") {
			return c.Status(fiber.StatusForbidden).JSON(
				dto.NewErrorResponse(dto.ErrCodeForbidden, "access denied to this document"),
			)
		}
		if strings.Contains(err.Error(), "not parsed yet") {
			return c.Status(fiber.StatusBadRequest).JSON(
				dto.NewErrorResponse(dto.ErrCodeDocumentNotParsed, "document not parsed yet"),
			)
		}
		if strings.Contains(err.Error(), "invalid LLM provider") {
			return c.Status(fiber.StatusBadRequest).JSON(
				dto.NewErrorResponse(dto.ErrCodeInvalidProvider, err.Error()),
			)
		}
		if strings.Contains(err.Error(), "failed to generate") {
			return c.Status(fiber.StatusInternalServerError).JSON(
				dto.NewErrorResponse(dto.ErrCodeGenerationFailed, "failed to generate questions"),
			)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeInternalError, "failed to generate test"),
		)
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

// List godoc
// @Summary List user's tests
// @Description Get paginated list of tests created by the current user. Admin sees all tests with user info, others see only their own
// @Tags tests
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} dto.TestListResponse
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /tests [get]
func (h *TestHandler) List(c *fiber.Ctx) error {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(
			dto.NewErrorResponse(dto.ErrCodeUnauthorized, "Unauthorized"),
		)
	}

	page, pageSize := c.QueryInt("page", 1), c.QueryInt("page_size", 20)

	result, err := h.listUseCase.Execute(c.Context(), userID, page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeDatabaseError, "failed to fetch tests"),
		)
	}

	return c.JSON(result)
}

// GetByID godoc
// @Summary Get test by ID
// @Description Get details of a specific test by its ID with questions and answers
// @Tags tests
// @Produce json
// @Security BearerAuth
// @Param id path string true "Test ID"
// @Success 200 {object} dto.TestResponse
// @Failure 400 {object} dto.ErrorResponse "Invalid test ID"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Access denied"
// @Failure 404 {object} dto.ErrorResponse "Test not found"
// @Router /tests/{id} [get]
func (h *TestHandler) GetByID(c *fiber.Ctx) error {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(
			dto.NewErrorResponse(dto.ErrCodeUnauthorized, "Unauthorized"),
		)
	}
	testID, _ := uuid.Parse(c.Params("id"))

	result, err := h.getUseCase.Execute(c.Context(), testID, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			dto.NewErrorResponse(dto.ErrCodeTestNotFound, "test not found"),
		)
	}

	return c.JSON(result)
}

// Delete godoc
// @Summary Delete a test
// @Description Delete a test and all its associated questions
// @Tags tests
// @Produce json
// @Security BearerAuth
// @Param id path string true "Test ID"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse "Invalid test ID"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Access denied"
// @Failure 404 {object} dto.ErrorResponse "Test not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /tests/{id} [delete]
func (h *TestHandler) Delete(c *fiber.Ctx) error {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(
			dto.NewErrorResponse(dto.ErrCodeUnauthorized, "Unauthorized"),
		)
	}
	testID, _ := uuid.Parse(c.Params("id"))

	if err := h.deleteUseCase.Execute(c.Context(), testID, userID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			dto.NewErrorResponse(dto.ErrCodeTestNotFound, "test not found"),
		)
	}

	return c.JSON(dto.NewMessageResponse("test deleted"))
}

// Update godoc
// @Summary Update a test
// @Description Update test title and description
// @Tags tests
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Test ID"
// @Param request body dto.UpdateTestRequest true "Update test request"
// @Success 200 {object} dto.TestResponse
// @Failure 400 {object} dto.ErrorResponse "Invalid request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Access denied"
// @Failure 404 {object} dto.ErrorResponse "Test not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /tests/{id} [put]
func (h *TestHandler) Update(c *fiber.Ctx) error {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(
			dto.NewErrorResponse(dto.ErrCodeUnauthorized, "Unauthorized"),
		)
	}

	testID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeInvalidInput, "invalid test ID"),
		)
	}

	var req dto.UpdateTestRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeInvalidInput, "invalid request body"),
		)
	}

	result, err := h.updateUseCase.Execute(c.Context(), testuc.UpdateParams{
		TestID:      testID,
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		if strings.Contains(err.Error(), "test not found") {
			return c.Status(fiber.StatusNotFound).JSON(
				dto.NewErrorResponse(dto.ErrCodeTestNotFound, "test not found"),
			)
		}
		if strings.Contains(err.Error(), "access denied") {
			return c.Status(fiber.StatusForbidden).JSON(
				dto.NewErrorResponse(dto.ErrCodeForbidden, "access denied"),
			)
		}
		if strings.Contains(err.Error(), "title must be at least") {
			return c.Status(fiber.StatusBadRequest).JSON(
				dto.NewErrorResponse(dto.ErrCodeInvalidInput, "title must be at least 3 characters"),
			)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeInternalError, "failed to update test"),
		)
	}

	return c.JSON(result)
}

// UpdateQuestion godoc
// @Summary Update a question
// @Description Update question text, type, difficulty, points and answers
// @Tags tests
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param testId path string true "Test ID"
// @Param questionId path string true "Question ID"
// @Param request body dto.UpdateQuestionRequest true "Update question request"
// @Success 200 {object} dto.QuestionDTO
// @Failure 400 {object} dto.ErrorResponse "Invalid request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Access denied"
// @Failure 404 {object} dto.ErrorResponse "Question not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /tests/{testId}/questions/{questionId} [put]
func (h *TestHandler) UpdateQuestion(c *fiber.Ctx) error {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(
			dto.NewErrorResponse(dto.ErrCodeUnauthorized, "Unauthorized"),
		)
	}

	testID, err := uuid.Parse(c.Params("testId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeInvalidInput, "invalid test ID"),
		)
	}

	questionID, err := uuid.Parse(c.Params("questionId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeInvalidInput, "invalid question ID"),
		)
	}

	var req dto.UpdateQuestionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeInvalidInput, "invalid request body"),
		)
	}

	result, err := h.updateQuestionUseCase.Execute(c.Context(), testuc.UpdateQuestionParams{
		TestID:       testID,
		QuestionID:   questionID,
		UserID:       userID,
		QuestionText: req.QuestionText,
		QuestionType: req.QuestionType,
		Difficulty:   req.Difficulty,
		Points:       req.Points,
		Answers:      req.Answers,
	})
	if err != nil {
		if strings.Contains(err.Error(), "test not found") {
			return c.Status(fiber.StatusNotFound).JSON(
				dto.NewErrorResponse(dto.ErrCodeTestNotFound, "test not found"),
			)
		}
		if strings.Contains(err.Error(), "access denied") {
			return c.Status(fiber.StatusForbidden).JSON(
				dto.NewErrorResponse(dto.ErrCodeForbidden, "access denied"),
			)
		}
		if strings.Contains(err.Error(), "question not found") {
			return c.Status(fiber.StatusNotFound).JSON(
				dto.NewErrorResponse(dto.ErrCodeNotFound, "question not found"),
			)
		}
		if strings.Contains(err.Error(), "question text must be at least") {
			return c.Status(fiber.StatusBadRequest).JSON(
				dto.NewErrorResponse(dto.ErrCodeInvalidInput, "question text must be at least 3 characters"),
			)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeInternalError, "failed to update question"),
		)
	}

	return c.JSON(result)
}

// ExportToJSON godoc
// @Summary Export test to JSON format
// @Description Export a test and its questions to JSON format for download
// @Tags tests
// @Produce json
// @Security BearerAuth
// @Param id path string true "Test ID"
// @Success 200 {object} dto.TestResponse
// @Failure 400 {object} dto.ErrorResponse "Invalid test ID"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Access denied"
// @Failure 404 {object} dto.ErrorResponse "Test not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /tests/{id}/export/json [get]
func (h *TestHandler) ExportToJSON(c *fiber.Ctx) error {
	return h.handleExport(c, "json")
}

// ExportToXML godoc
// @Summary Export test to Moodle XML format
// @Description Export a test and its questions to Moodle XML format for download
// @Tags tests
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "Test ID"
// @Success 200 {string} string "XML file"
// @Failure 400 {object} dto.ErrorResponse "Invalid test ID or test has no questions"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Access denied"
// @Failure 404 {object} dto.ErrorResponse "Test not found"
// @Failure 500 {object} dto.ErrorResponse "Export failed"
// @Router /tests/{id}/export/xml [get]
func (h *TestHandler) ExportToXML(c *fiber.Ctx) error {
	return h.handleExport(c, "moodle_xml")
}

// Export godoc
// @Summary Export test to various LMS formats
// @Description Export a test to different LMS formats (moodle_xml, stepik_csv, gift, aiken, blackboard, qti)
// @Tags tests
// @Produce application/octet-stream
// @Security BearerAuth
// @Param id path string true "Test ID"
// @Param format query string true "Export format" Enums(moodle_xml, stepik_csv, gift, aiken, blackboard, qti)
// @Success 200 {string} string "Exported file"
// @Failure 400 {object} dto.ErrorResponse "Invalid test ID, format, or test has no questions"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 403 {object} dto.ErrorResponse "Access denied"
// @Failure 404 {object} dto.ErrorResponse "Test not found"
// @Failure 500 {object} dto.ErrorResponse "Export failed"
// @Router /tests/{id}/export [get]
func (h *TestHandler) Export(c *fiber.Ctx) error {
	format := c.Query("format", "moodle_xml")
	return h.handleExport(c, format)
}

// GetExportFormats godoc
// @Summary Get available export formats
// @Description Get list of all available export formats with descriptions
// @Tags tests
// @Produce json
// @Security BearerAuth
// @Success 200 {array} exporter.FormatInfo
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Router /tests/export/formats [get]
func (h *TestHandler) GetExportFormats(c *fiber.Ctx) error {
	_, ok := getUserIDFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(
			dto.NewErrorResponse(dto.ErrCodeUnauthorized, "Unauthorized"),
		)
	}

	formats := h.exportUseCase.GetAvailableFormats()
	return c.JSON(formats)
}

// handleExport is a shared helper for all export endpoints
func (h *TestHandler) handleExport(c *fiber.Ctx, format string) error {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(
			dto.NewErrorResponse(dto.ErrCodeUnauthorized, "Unauthorized"),
		)
	}

	testID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeInvalidTestID, "invalid test ID"),
		)
	}

	result, err := h.exportUseCase.Execute(c.Context(), testID, userID, format)
	if err != nil {
		if strings.Contains(err.Error(), "test not found") {
			return c.Status(fiber.StatusNotFound).JSON(
				dto.NewErrorResponse(dto.ErrCodeTestNotFound, "test not found"),
			)
		}
		if strings.Contains(err.Error(), "access denied") {
			return c.Status(fiber.StatusForbidden).JSON(
				dto.NewErrorResponse(dto.ErrCodeForbidden, "access denied"),
			)
		}
		if strings.Contains(err.Error(), "no questions") {
			return c.Status(fiber.StatusBadRequest).JSON(
				dto.NewErrorResponse(dto.ErrCodeTestHasNoQuestions, "test has no questions"),
			)
		}
		if strings.Contains(err.Error(), "unsupported export format") {
			return c.Status(fiber.StatusBadRequest).JSON(
				dto.NewErrorResponse(dto.ErrCodeInvalidInput, err.Error()),
			)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeExportFailed, "failed to export"),
		)
	}

	filename := result.Title + "." + result.FileExt
	c.Set("Content-Type", result.ContentType)
	c.Set("Content-Disposition", "attachment; filename="+filename)

	return c.SendString(result.Content)
}
