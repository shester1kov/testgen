package handler

import (
	"time"

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
	questionRepo repository.QuestionRepository
	answerRepo   repository.AnswerRepository
	llmFactory   *llm.LLMFactory
}

func NewTestHandler(
	testRepo repository.TestRepository,
	documentRepo repository.DocumentRepository,
	questionRepo repository.QuestionRepository,
	answerRepo repository.AnswerRepository,
	llmFactory *llm.LLMFactory,
) *TestHandler {
	return &TestHandler{
		testRepo:     testRepo,
		documentRepo: documentRepo,
		questionRepo: questionRepo,
		answerRepo:   answerRepo,
		llmFactory:   llmFactory,
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

	// Create Test entity
	test := &entity.Test{
		ID:             uuid.New(),
		UserID:         userID,
		DocumentID:     &docID,
		Title:          req.Title,
		TotalQuestions: len(questions),
		Status:         entity.TestStatusDraft,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Save test to database
	if err := h.testRepo.Create(c.Context(), test); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeInternalError, "failed to save test"),
		)
	}

	// Save questions and answers to database
	for i, q := range questions {
		question := &entity.Question{
			ID:           uuid.New(),
			TestID:       test.ID,
			QuestionText: q.QuestionText,
			QuestionType: entity.QuestionType(q.QuestionType),
			Difficulty:   entity.Difficulty(q.Difficulty),
			Points:       1.0, // Default points
			OrderNum:     i + 1,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		if err := h.questionRepo.Create(c.Context(), question); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				dto.NewErrorResponse(dto.ErrCodeInternalError, "failed to save question"),
			)
		}

		// Save answers
		for j, a := range q.Answers {
			answer := &entity.Answer{
				ID:         uuid.New(),
				QuestionID: question.ID,
				AnswerText: a.Text,
				IsCorrect:  a.IsCorrect,
				OrderNum:   j + 1,
				CreatedAt:  time.Now(),
			}

			if err := h.answerRepo.Create(c.Context(), answer); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(
					dto.NewErrorResponse(dto.ErrCodeInternalError, "failed to save answer"),
				)
			}
		}
	}

	// Return test response with ID
	return c.Status(fiber.StatusCreated).JSON(dto.TestResponse{
		ID:             test.ID.String(),
		Title:          test.Title,
		Description:    "", // Empty for generated tests
		TotalQuestions: test.TotalQuestions,
		Status:         string(test.Status),
		MoodleSynced:   false,
		CreatedAt:      test.CreatedAt.Format(time.RFC3339),
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
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(
			dto.NewErrorResponse(dto.ErrCodeUnauthorized, "Unauthorized"),
		)
	}
	page, pageSize := c.QueryInt("page", 1), c.QueryInt("page_size", 20)
	offset := (page - 1) * pageSize

	tests, _ := h.testRepo.FindByUserID(c.Context(), userID, pageSize, offset)
	total, _ := h.testRepo.CountByUserID(c.Context(), userID)

	result := make([]dto.TestResponse, len(tests))
	for i, t := range tests {
		result[i] = dto.TestResponse{
			ID:             t.ID.String(),
			Title:          t.Title,
			Description:    t.Description,
			TotalQuestions: t.TotalQuestions,
			Status:         string(t.Status),
			MoodleSynced:   t.MoodleSynced,
			CreatedAt:      t.CreatedAt.Format(time.RFC3339),
		}
	}

	return c.JSON(dto.TestListResponse{Tests: result, Total: total, Page: page, PageSize: pageSize})
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
	test, err := h.testRepo.FindByID(c.Context(), testID)
	if err != nil || test.UserID != userID {
		return c.Status(fiber.StatusNotFound).JSON(
			dto.NewErrorResponse(dto.ErrCodeTestNotFound, "test not found"),
		)
	}

	// Load questions with answers
	questions, err := h.questionRepo.FindByTestID(c.Context(), testID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeDatabaseError, "failed to load questions"),
		)
	}

	// Build questions DTO with answers
	questionsDTO := make([]dto.QuestionDTO, len(questions))
	for i, q := range questions {
		// Load answers for each question
		answers, err := h.answerRepo.FindByQuestionID(c.Context(), q.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				dto.NewErrorResponse(dto.ErrCodeDatabaseError, "failed to load answers"),
			)
		}

		// Build answers DTO
		answersDTO := make([]dto.AnswerDTO, len(answers))
		for j, a := range answers {
			answersDTO[j] = dto.AnswerDTO{
				ID:         a.ID.String(),
				AnswerText: a.AnswerText,
				IsCorrect:  a.IsCorrect,
				OrderNum:   a.OrderNum,
			}
		}

		questionsDTO[i] = dto.QuestionDTO{
			ID:           q.ID.String(),
			QuestionText: q.QuestionText,
			QuestionType: string(q.QuestionType),
			Difficulty:   string(q.Difficulty),
			Points:       q.Points,
			OrderNum:     q.OrderNum,
			Answers:      answersDTO,
		}
	}

	return c.JSON(dto.TestResponse{
		ID:             test.ID.String(),
		Title:          test.Title,
		Description:    test.Description,
		TotalQuestions: test.TotalQuestions,
		Status:         string(test.Status),
		MoodleSynced:   test.MoodleSynced,
		CreatedAt:      test.CreatedAt.Format(time.RFC3339),
		Questions:      questionsDTO,
	})
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
	test, err := h.testRepo.FindByID(c.Context(), testID)
	if err != nil || test.UserID != userID {
		return c.Status(fiber.StatusNotFound).JSON(
			dto.NewErrorResponse(dto.ErrCodeTestNotFound, "test not found"),
		)
	}

	h.testRepo.Delete(c.Context(), testID)
	return c.JSON(dto.NewMessageResponse("test deleted"))
}
