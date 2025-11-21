package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/application/dto"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/llm"
)

type TestHandler struct {
	testRepo     repository.TestRepository
	documentRepo repository.DocumentRepository
	llmFactory   *llm.LLMFactory
}

func NewTestHandler(testRepo repository.TestRepository, documentRepo repository.DocumentRepository, llmFactory *llm.LLMFactory) *TestHandler {
	return &TestHandler{testRepo: testRepo, documentRepo: documentRepo, llmFactory: llmFactory}
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
	userID := c.Locals("userID").(uuid.UUID)
	var req dto.CreateTestRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeInvalidInput, "invalid request"),
		)
	}

	test := &entity.Test{UserID: userID, Title: req.Title, Description: req.Description}
	if req.DocumentID != nil {
		docID, _ := uuid.Parse(*req.DocumentID)
		test.DocumentID = &docID
	}

	if err := h.testRepo.Create(c.Context(), test); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeDatabaseError, "failed to create test"),
		)
	}

	return c.Status(fiber.StatusCreated).JSON(dto.TestResponse{
		ID: test.ID.String(), Title: test.Title, Description: test.Description,
		TotalQuestions: 0, Status: string(test.Status), CreatedAt: test.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// Generate godoc
// @Summary Generate test questions
// @Description Generate test questions from a document using LLM
// @Tags tests
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.GenerateTestRequest true "Generate test request"
// @Success 200 {object} dto.GenerateQuestionsResponse
// @Failure 400 {object} dto.ErrorResponse "Invalid input or document not parsed"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 404 {object} dto.ErrorResponse "Document not found"
// @Failure 500 {object} dto.ErrorResponse "Generation failed"
// @Router /tests/generate [post]
func (h *TestHandler) Generate(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	var req dto.GenerateTestRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeInvalidInput, "invalid request"),
		)
	}

	docID, _ := uuid.Parse(req.DocumentID)
	document, err := h.documentRepo.FindByID(c.Context(), docID)
	if err != nil || document.UserID != userID {
		return c.Status(fiber.StatusNotFound).JSON(
			dto.NewErrorResponse(dto.ErrCodeDocumentNotFound, "document not found"),
		)
	}

	if !document.IsParsed() {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeDocumentNotParsed, "document not parsed yet"),
		)
	}

	provider := req.LLMProvider
	if provider == "" {
		provider = "perplexity"
	}

	strategy, err := h.llmFactory.CreateStrategy(provider)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeInvalidProvider, err.Error()),
		)
	}

	llmContext := llm.NewLLMContext(strategy)
	questions, err := llmContext.GenerateQuestions(c.Context(), llm.GenerationParams{
		Text: document.ParsedText, NumQuestions: req.NumQuestions, Difficulty: req.Difficulty,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeGenerationFailed, "failed to generate questions"),
		)
	}

	// Convert llm.GeneratedQuestion to dto.QuestionDTO
	questionDTOs := make([]dto.QuestionDTO, len(questions))
	for i, q := range questions {
		answers := make([]dto.AnswerDTO, len(q.Answers))
		for j, a := range q.Answers {
			answers[j] = dto.AnswerDTO{
				AnswerText: a.Text,
				IsCorrect:  a.IsCorrect,
				OrderNum:   j + 1,
			}
		}
		questionDTOs[i] = dto.QuestionDTO{
			QuestionText: q.QuestionText,
			QuestionType: string(q.QuestionType),
			Difficulty:   q.Difficulty,
			OrderNum:     i + 1,
			Answers:      answers,
		}
	}

	return c.JSON(dto.GenerateQuestionsResponse{
		Message:   "questions generated",
		Count:     len(questions),
		Questions: questionDTOs,
	})
}

// List godoc
// @Summary List user's tests
// @Description Get paginated list of tests created by the current user
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
	userID := c.Locals("userID").(uuid.UUID)
	page, pageSize := c.QueryInt("page", 1), c.QueryInt("page_size", 20)
	offset := (page - 1) * pageSize

	tests, _ := h.testRepo.FindByUserID(c.Context(), userID, pageSize, offset)
	total, _ := h.testRepo.CountByUserID(c.Context(), userID)

	result := make([]dto.TestResponse, len(tests))
	for i, t := range tests {
		result[i] = dto.TestResponse{ID: t.ID.String(), Title: t.Title, TotalQuestions: t.TotalQuestions, Status: string(t.Status)}
	}

	return c.JSON(dto.TestListResponse{Tests: result, Total: total, Page: page, PageSize: pageSize})
}

// GetByID godoc
// @Summary Get test by ID
// @Description Get details of a specific test by its ID
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
	userID := c.Locals("userID").(uuid.UUID)
	testID, _ := uuid.Parse(c.Params("id"))
	test, err := h.testRepo.FindByID(c.Context(), testID)
	if err != nil || test.UserID != userID {
		return c.Status(fiber.StatusNotFound).JSON(
			dto.NewErrorResponse(dto.ErrCodeTestNotFound, "test not found"),
		)
	}

	return c.JSON(dto.TestResponse{ID: test.ID.String(), Title: test.Title, Description: test.Description, TotalQuestions: test.TotalQuestions})
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
	userID := c.Locals("userID").(uuid.UUID)
	testID, _ := uuid.Parse(c.Params("id"))
	test, err := h.testRepo.FindByID(c.Context(), testID)
	if err != nil || test.UserID != userID {
		return c.Status(fiber.StatusNotFound).JSON(
			dto.NewErrorResponse(dto.ErrCodeTestNotFound, "test not found"),
		)
	}

	h.testRepo.Delete(c.Context(), testID)
	return c.JSON(dto.NewMessageResponse("test deleted"))
}
