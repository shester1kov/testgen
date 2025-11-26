package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/application/dto"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/llm"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/moodle"
)

type TestHandler struct {
	testRepo     repository.TestRepository
	documentRepo repository.DocumentRepository
	questionRepo repository.QuestionRepository
	answerRepo   repository.AnswerRepository
	userRepo     repository.UserRepository
	llmFactory   *llm.LLMFactory
	xmlExporter  *moodle.MoodleXMLExporter
}

func NewTestHandler(
	testRepo repository.TestRepository,
	documentRepo repository.DocumentRepository,
	questionRepo repository.QuestionRepository,
	answerRepo repository.AnswerRepository,
	userRepo repository.UserRepository,
	llmFactory *llm.LLMFactory,
	xmlExporter *moodle.MoodleXMLExporter,
) *TestHandler {
	return &TestHandler{
		testRepo:     testRepo,
		documentRepo: documentRepo,
		questionRepo: questionRepo,
		answerRepo:   answerRepo,
		userRepo:     userRepo,
		llmFactory:   llmFactory,
		xmlExporter:  xmlExporter,
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
		ID: test.ID.String(), UserID: test.UserID.String(), Title: test.Title, Description: test.Description,
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
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			dto.NewErrorResponse(dto.ErrCodeDocumentNotFound, "document not found"),
		)
	}

	// Check if user has access to document (admin sees all, others see only their own)
	user, err := h.userRepo.FindByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeDatabaseError, "failed to fetch user"),
		)
	}

	// Non-admin users can only generate tests from their own documents
	if !user.IsAdmin() && document.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(
			dto.NewErrorResponse(dto.ErrCodeForbidden, "access denied to this document"),
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
		UserID:         test.UserID.String(),
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

	// Get user to check role
	user, err := h.userRepo.FindByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeDatabaseError, "failed to fetch user"),
		)
	}

	page, pageSize := c.QueryInt("page", 1), c.QueryInt("page_size", 20)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var tests []*entity.Test
	var total int64

	// Admin sees all tests, others see only their own
	if user.IsAdmin() {
		tests, err = h.testRepo.FindAll(c.Context(), pageSize, offset)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				dto.NewErrorResponse(dto.ErrCodeDatabaseError, "failed to fetch tests"),
			)
		}

		total, err = h.testRepo.CountAll(c.Context())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				dto.NewErrorResponse(dto.ErrCodeDatabaseError, "failed to count tests"),
			)
		}
	} else {
		tests, err = h.testRepo.FindByUserID(c.Context(), userID, pageSize, offset)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				dto.NewErrorResponse(dto.ErrCodeDatabaseError, "failed to fetch tests"),
			)
		}

		total, err = h.testRepo.CountByUserID(c.Context(), userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				dto.NewErrorResponse(dto.ErrCodeDatabaseError, "failed to count tests"),
			)
		}
	}

	result := make([]dto.TestResponse, len(tests))
	for i, t := range tests {
		// Include user info for admin
		var userName *string
		var userEmail *string
		if user.IsAdmin() && t.User.ID != uuid.Nil {
			userName = &t.User.FullName
			userEmail = &t.User.Email
		}

		result[i] = dto.TestResponse{
			ID:             t.ID.String(),
			UserID:         t.UserID.String(),
			UserName:       userName,
			UserEmail:      userEmail,
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

	// Check if test exists and belongs to user
	test, err := h.testRepo.FindByID(c.Context(), testID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			dto.NewErrorResponse(dto.ErrCodeTestNotFound, "test not found"),
		)
	}

	if test.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(
			dto.NewErrorResponse(dto.ErrCodeForbidden, "access denied"),
		)
	}

	// Parse request
	var req dto.UpdateTestRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeInvalidInput, "invalid request body"),
		)
	}

	// Validate title if provided
	if req.Title != "" && len(req.Title) < 3 {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeInvalidInput, "title must be at least 3 characters"),
		)
	}

	// Update fields if provided
	if req.Title != "" {
		test.Title = req.Title
	}
	test.Description = req.Description // Allow empty description
	test.UpdatedAt = time.Now()

	// Save to database
	if err := h.testRepo.Update(c.Context(), test); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeInternalError, "failed to update test"),
		)
	}

	return c.JSON(dto.TestResponse{
		ID:             test.ID.String(),
		UserID:         test.UserID.String(),
		Title:          test.Title,
		Description:    test.Description,
		TotalQuestions: test.TotalQuestions,
		Status:         string(test.Status),
		MoodleSynced:   test.MoodleSynced,
		CreatedAt:      test.CreatedAt.Format(time.RFC3339),
	})
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

	// Check if test exists and belongs to user
	test, err := h.testRepo.FindByID(c.Context(), testID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			dto.NewErrorResponse(dto.ErrCodeTestNotFound, "test not found"),
		)
	}

	if test.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(
			dto.NewErrorResponse(dto.ErrCodeForbidden, "access denied"),
		)
	}

	// Check if question exists and belongs to test
	question, err := h.questionRepo.FindByID(c.Context(), questionID)
	if err != nil || question.TestID != testID {
		return c.Status(fiber.StatusNotFound).JSON(
			dto.NewErrorResponse(dto.ErrCodeNotFound, "question not found"),
		)
	}

	// Parse request
	var req dto.UpdateQuestionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeInvalidInput, "invalid request body"),
		)
	}

	// Validate question text if provided
	if req.QuestionText != "" && len(req.QuestionText) < 3 {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeInvalidInput, "question text must be at least 3 characters"),
		)
	}

	// Update fields if provided
	if req.QuestionText != "" {
		question.QuestionText = req.QuestionText
	}
	if req.QuestionType != "" {
		question.QuestionType = entity.QuestionType(req.QuestionType)
	}
	if req.Difficulty != "" {
		question.Difficulty = entity.Difficulty(req.Difficulty)
	}
	if req.Points != nil {
		question.Points = *req.Points
	}
	question.UpdatedAt = time.Now()

	// Save question
	if err := h.questionRepo.Update(c.Context(), question); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeInternalError, "failed to update question"),
		)
	}

	// Update answers if provided
	if len(req.Answers) > 0 {
		// Delete old answers
		oldAnswers, _ := h.answerRepo.FindByQuestionID(c.Context(), questionID)
		for _, oldAnswer := range oldAnswers {
			h.answerRepo.Delete(c.Context(), oldAnswer.ID)
		}

		// Create new answers
		for _, answerReq := range req.Answers {
			answer := &entity.Answer{
				ID:         uuid.New(),
				QuestionID: questionID,
				AnswerText: answerReq.AnswerText,
				IsCorrect:  answerReq.IsCorrect,
				OrderNum:   answerReq.OrderNum,
				CreatedAt:  time.Now(),
			}
			if err := h.answerRepo.Create(c.Context(), answer); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(
					dto.NewErrorResponse(dto.ErrCodeInternalError, "failed to create answer"),
				)
			}
		}
	}

	// Load updated answers
	answers, err := h.answerRepo.FindByQuestionID(c.Context(), questionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeDatabaseError, "failed to load answers"),
		)
	}

	// Build answers DTO
	answersDTO := make([]dto.AnswerDTO, len(answers))
	for i, a := range answers {
		answersDTO[i] = dto.AnswerDTO{
			ID:         a.ID.String(),
			AnswerText: a.AnswerText,
			IsCorrect:  a.IsCorrect,
			OrderNum:   a.OrderNum,
		}
	}

	return c.JSON(dto.QuestionDTO{
		ID:           question.ID.String(),
		QuestionText: question.QuestionText,
		QuestionType: string(question.QuestionType),
		Difficulty:   string(question.Difficulty),
		Points:       question.Points,
		OrderNum:     question.OrderNum,
		Answers:      answersDTO,
	})
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

	// Check if test exists and belongs to user
	test, err := h.testRepo.FindByID(c.Context(), testID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			dto.NewErrorResponse(dto.ErrCodeTestNotFound, "test not found"),
		)
	}

	if test.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(
			dto.NewErrorResponse(dto.ErrCodeForbidden, "access denied"),
		)
	}

	// Load questions with answers
	questions, err := h.questionRepo.FindByTestID(c.Context(), testID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeDatabaseError, "failed to load questions"),
		)
	}

	if len(questions) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeTestHasNoQuestions, "test has no questions"),
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

	testResponse := dto.TestResponse{
		ID:             test.ID.String(),
		UserID:         test.UserID.String(),
		Title:          test.Title,
		Description:    test.Description,
		TotalQuestions: test.TotalQuestions,
		Status:         string(test.Status),
		MoodleSynced:   test.MoodleSynced,
		CreatedAt:      test.CreatedAt.Format(time.RFC3339),
		Questions:      questionsDTO,
	}

	// Set content disposition header for file download
	filename := "test_" + test.ID.String() + ".json"
	c.Set("Content-Disposition", "attachment; filename="+filename)
	c.Set("Content-Type", "application/json")

	return c.JSON(testResponse)
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

	// Check if test exists and belongs to user
	test, err := h.testRepo.FindByID(c.Context(), testID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			dto.NewErrorResponse(dto.ErrCodeTestNotFound, "test not found"),
		)
	}

	if test.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(
			dto.NewErrorResponse(dto.ErrCodeForbidden, "access denied"),
		)
	}

	// Get questions for the test
	questions, err := h.questionRepo.FindByTestID(c.Context(), testID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeDatabaseError, "failed to retrieve questions"),
		)
	}

	if len(questions) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeTestHasNoQuestions, "test has no questions"),
		)
	}

	// Get answers for each question
	answersMap := make(map[string][]*entity.Answer)
	for _, q := range questions {
		answers, err := h.answerRepo.FindByQuestionID(c.Context(), q.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				dto.NewErrorResponse(dto.ErrCodeDatabaseError, "failed to retrieve answers"),
			)
		}
		answersMap[q.ID.String()] = answers
	}

	// Export to XML
	xmlContent, err := h.xmlExporter.Export(test, questions, answersMap)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeExportFailed, "failed to export XML"),
		)
	}

	// Set headers for file download
	filename := test.Title + ".xml"
	c.Set("Content-Type", "application/xml")
	c.Set("Content-Disposition", "attachment; filename="+filename)

	return c.SendString(xmlContent)
}
