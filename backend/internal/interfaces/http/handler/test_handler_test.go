package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/application/dto"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/llm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockTestRepository struct{ mock.Mock }

func (m *mockTestRepository) Create(ctx context.Context, test *entity.Test) error {
	args := m.Called(ctx, test)
	return args.Error(0)
}

func (m *mockTestRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Test, error) {
	args := m.Called(ctx, id)
	if res := args.Get(0); res != nil {
		return res.(*entity.Test), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockTestRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Test, error) {
	args := m.Called(ctx, userID, limit, offset)
	if res := args.Get(0); res != nil {
		return res.([]*entity.Test), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockTestRepository) Update(ctx context.Context, test *entity.Test) error { return nil }

func (m *mockTestRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockTestRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockTestRepository) FindAll(ctx context.Context, limit, offset int) ([]*entity.Test, error) {
	return nil, nil
}

func (m *mockTestRepository) CountAll(ctx context.Context) (int64, error) {
	return 0, nil
}

type mockTestDocRepository struct{ mock.Mock }

func (m *mockTestDocRepository) Create(ctx context.Context, doc *entity.Document) error { return nil }
func (m *mockTestDocRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Document, error) {
	args := m.Called(ctx, id)
	if res := args.Get(0); res != nil {
		return res.(*entity.Document), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockTestDocRepository) List(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Document, error) {
	return nil, nil
}
func (m *mockTestDocRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Document, error) {
	args := m.Called(ctx, userID, limit, offset)
	if res := args.Get(0); res != nil {
		return res.([]*entity.Document), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockTestDocRepository) Update(ctx context.Context, doc *entity.Document) error { return nil }
func (m *mockTestDocRepository) Delete(ctx context.Context, id uuid.UUID) error         { return nil }
func (m *mockTestDocRepository) MarkAsParsed(ctx context.Context, id uuid.UUID, parsedText string) error {
	return nil
}
func (m *mockTestDocRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockTestDocRepository) FindAll(ctx context.Context, limit, offset int) ([]*entity.Document, error) {
	return nil, nil
}

func (m *mockTestDocRepository) CountAll(ctx context.Context) (int64, error) {
	return 0, nil
}

type mockQuestionRepository struct{ mock.Mock }

func (m *mockQuestionRepository) Create(ctx context.Context, question *entity.Question) error {
	args := m.Called(ctx, question)
	return args.Error(0)
}
func (m *mockQuestionRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Question, error) {
	return nil, nil
}
func (m *mockQuestionRepository) FindByTestID(ctx context.Context, testID uuid.UUID) ([]*entity.Question, error) {
	args := m.Called(ctx, testID)
	if res := args.Get(0); res != nil {
		return res.([]*entity.Question), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockQuestionRepository) Update(ctx context.Context, question *entity.Question) error { return nil }
func (m *mockQuestionRepository) Delete(ctx context.Context, id uuid.UUID) error              { return nil }
func (m *mockQuestionRepository) CountByTestID(ctx context.Context, testID uuid.UUID) (int, error) {
	return 0, nil
}
func (m *mockQuestionRepository) ReorderQuestions(ctx context.Context, testID uuid.UUID, questionIDs []uuid.UUID) error {
	return nil
}

func (m *mockQuestionRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	return 0, nil
}

func (m *mockQuestionRepository) CountAll(ctx context.Context) (int64, error) {
	return 0, nil
}

type mockAnswerRepository struct{ mock.Mock }

func (m *mockAnswerRepository) Create(ctx context.Context, answer *entity.Answer) error {
	args := m.Called(ctx, answer)
	return args.Error(0)
}
func (m *mockAnswerRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Answer, error) {
	return nil, nil
}
func (m *mockAnswerRepository) FindByQuestionID(ctx context.Context, questionID uuid.UUID) ([]*entity.Answer, error) {
	args := m.Called(ctx, questionID)
	if res := args.Get(0); res != nil {
		return res.([]*entity.Answer), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockAnswerRepository) Update(ctx context.Context, answer *entity.Answer) error { return nil }
func (m *mockAnswerRepository) Delete(ctx context.Context, id uuid.UUID) error          { return nil }
func (m *mockAnswerRepository) DeleteByQuestionID(ctx context.Context, questionID uuid.UUID) error {
	return nil
}

type mockTestUserRepository struct{ mock.Mock }

func (m *mockTestUserRepository) Create(ctx context.Context, user *entity.User) error { return nil }
func (m *mockTestUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	args := m.Called(ctx, id)
	if res := args.Get(0); res != nil {
		return res.(*entity.User), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockTestUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	return nil, nil
}
func (m *mockTestUserRepository) Update(ctx context.Context, user *entity.User) error { return nil }
func (m *mockTestUserRepository) Delete(ctx context.Context, id uuid.UUID) error      { return nil }
func (m *mockTestUserRepository) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	return nil, nil
}
func (m *mockTestUserRepository) Count(ctx context.Context) (int64, error) { return 0, nil }

func TestCreateTest_Success(t *testing.T) {
	userID := uuid.New()
	testRepo := new(mockTestRepository)
	docRepo := new(mockTestDocRepository)

	testRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Test")).Run(func(args mock.Arguments) {
		test := args.Get(1).(*entity.Test)
		test.ID = uuid.New()
	}).Return(nil)

	handler := NewTestHandler(testRepo, docRepo, new(mockQuestionRepository), new(mockAnswerRepository), new(mockTestUserRepository), nil, nil)
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error { c.Locals("userID", userID); return c.Next() })
	app.Post("/tests", handler.Create)

	body, _ := json.Marshal(dto.CreateTestRequest{Title: "Sample", Description: "Desc"})
	req := httptest.NewRequest(http.MethodPost, "/tests", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	testRepo.AssertExpectations(t)
}

func TestCreateTest_InvalidBody(t *testing.T) {
	handler := NewTestHandler(new(mockTestRepository), new(mockTestDocRepository), new(mockQuestionRepository), new(mockAnswerRepository), new(mockTestUserRepository), nil, nil)
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error { c.Locals("userID", uuid.New()); return c.Next() })
	app.Post("/tests", handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/tests", bytes.NewBufferString("{invalid"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestGenerate_Success(t *testing.T) {
	userID := uuid.New()
	docID := uuid.New()

	testRepo := new(mockTestRepository)
	docRepo := new(mockTestDocRepository)
	questionRepo := new(mockQuestionRepository)
	answerRepo := new(mockAnswerRepository)

	// Mock document
	document := &entity.Document{
		ID:         docID,
		UserID:     userID,
		Title:      "Test Document",
		ParsedText: "Sample text for testing",
		Status:     entity.StatusParsed,
	}
	docRepo.On("FindByID", mock.Anything, docID).Return(document, nil)

	// Mock user repository to return a non-admin user who owns the document
	userRepo := new(mockTestUserRepository)
	user := &entity.User{
		ID:    userID,
		Email: "test@example.com",
		Role: &entity.Role{
			ID:   uuid.New(),
			Name: entity.RoleNameTeacher, // Non-admin role
		},
	}
	userRepo.On("FindByID", mock.Anything, userID).Return(user, nil)

	// Mock test creation
	testRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Test")).Return(nil)

	// Mock question creation
	questionRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Question")).Return(nil)

	// Mock answer creation
	answerRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Answer")).Return(nil)

	// Use perplexity provider which returns mock data in tests
	mockFactory := llm.NewLLMFactory("test-key", "", "", "", "")

	handler := NewTestHandler(testRepo, docRepo, questionRepo, answerRepo, userRepo, mockFactory, nil)
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error { c.Locals("userID", userID); return c.Next() })
	app.Post("/tests/generate", handler.Generate)

	body, _ := json.Marshal(dto.GenerateTestRequest{
		DocumentID:   docID.String(),
		Title:        "Generated Test",
		NumQuestions: 2,
		Difficulty:   "medium",
		LLMProvider:  "perplexity", // Use perplexity which exists in factory
	})
	req := httptest.NewRequest(http.MethodPost, "/tests/generate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	// Verify response body
	var response dto.TestResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.NotEmpty(t, response.ID)
	assert.Equal(t, "Generated Test", response.Title)
	assert.Equal(t, "draft", response.Status)
}

func TestGenerate_DocumentNotFound(t *testing.T) {
	userID := uuid.New()
	docRepo := new(mockTestDocRepository)
	docRepo.On("FindByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil, assert.AnError)

	handler := NewTestHandler(new(mockTestRepository), docRepo, new(mockQuestionRepository), new(mockAnswerRepository), new(mockTestUserRepository), nil, nil)
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error { c.Locals("userID", userID); return c.Next() })
	app.Post("/tests/generate", handler.Generate)

	body, _ := json.Marshal(dto.GenerateTestRequest{
		DocumentID:   uuid.New().String(),
		Title:        "Test",
		NumQuestions: 1,
		Difficulty:   "easy",
	})
	req := httptest.NewRequest(http.MethodPost, "/tests/generate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestGenerate_DocumentNotParsed(t *testing.T) {
	userID := uuid.New()
	docID := uuid.New()

	docRepo := new(mockTestDocRepository)
	document := &entity.Document{
		ID:     docID,
		UserID: userID,
		Status: entity.StatusUploaded, // Not parsed yet
	}
	docRepo.On("FindByID", mock.Anything, docID).Return(document, nil)

	userRepo := new(mockTestUserRepository)
	user := &entity.User{
		ID:    userID,
		Email: "test@example.com",
		Role: &entity.Role{
			ID:   uuid.New(),
			Name: entity.RoleNameTeacher,
		},
	}
	userRepo.On("FindByID", mock.Anything, userID).Return(user, nil)

	handler := NewTestHandler(new(mockTestRepository), docRepo, new(mockQuestionRepository), new(mockAnswerRepository), userRepo, nil, nil)
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error { c.Locals("userID", userID); return c.Next() })
	app.Post("/tests/generate", handler.Generate)

	body, _ := json.Marshal(dto.GenerateTestRequest{
		DocumentID:   docID.String(),
		Title:        "Test",
		NumQuestions: 1,
		Difficulty:   "easy",
	})
	req := httptest.NewRequest(http.MethodPost, "/tests/generate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestGenerate_InvalidProvider(t *testing.T) {
	userID := uuid.New()
	docID := uuid.New()

	docRepo := new(mockTestDocRepository)
	document := &entity.Document{
		ID:         docID,
		UserID:     userID,
		ParsedText: "text",
		Status:     entity.StatusParsed,
	}
	docRepo.On("FindByID", mock.Anything, docID).Return(document, nil)

	userRepo := new(mockTestUserRepository)
	user := &entity.User{
		ID:    userID,
		Email: "test@example.com",
		Role: &entity.Role{
			ID:   uuid.New(),
			Name: entity.RoleNameTeacher,
		},
	}
	userRepo.On("FindByID", mock.Anything, userID).Return(user, nil)

	// Factory will return error for invalid provider (empty factory)
	mockFactory := llm.NewLLMFactory("", "", "", "", "")

	handler := NewTestHandler(new(mockTestRepository), docRepo, new(mockQuestionRepository), new(mockAnswerRepository), userRepo, mockFactory, nil)
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error { c.Locals("userID", userID); return c.Next() })
	app.Post("/tests/generate", handler.Generate)

	body, _ := json.Marshal(dto.GenerateTestRequest{
		DocumentID:   docID.String(),
		Title:        "Test",
		NumQuestions: 1,
		Difficulty:   "easy",
		LLMProvider:  "invalid_provider",
	})
	req := httptest.NewRequest(http.MethodPost, "/tests/generate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestListTests_Success(t *testing.T) {
	userID := uuid.New()
	testRepo := new(mockTestRepository)
	docRepo := new(mockTestDocRepository)
	userRepo := new(mockTestUserRepository)

	// Mock user with teacher role
	teacherRole := &entity.Role{ID: uuid.New(), Name: "teacher"}
	teacherUser := &entity.User{ID: userID, Email: "teacher@test.com", RoleID: teacherRole.ID, Role: teacherRole}
	userRepo.On("FindByID", mock.Anything, userID).Return(teacherUser, nil)

	testRepo.On("FindByUserID", mock.Anything, userID, 20, 0).Return([]*entity.Test{{ID: uuid.New(), Title: "T1", UserID: userID}}, nil)
	testRepo.On("CountByUserID", mock.Anything, userID).Return(int64(1), nil)

	handler := NewTestHandler(testRepo, docRepo, new(mockQuestionRepository), new(mockAnswerRepository), userRepo, nil, nil)
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error { c.Locals("userID", userID); return c.Next() })
	app.Get("/tests", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/tests", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestGetByID_NotFound(t *testing.T) {
	userID := uuid.New()
	testRepo := new(mockTestRepository)
	testRepo.On("FindByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil, assert.AnError)

	handler := NewTestHandler(testRepo, new(mockTestDocRepository), new(mockQuestionRepository), new(mockAnswerRepository), new(mockTestUserRepository), nil, nil)
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error { c.Locals("userID", userID); return c.Next() })
	app.Get("/tests/:id", handler.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/tests/"+uuid.New().String(), nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestDeleteTest_Success(t *testing.T) {
	userID := uuid.New()
	testID := uuid.New()
	testRepo := new(mockTestRepository)

	testRepo.On("FindByID", mock.Anything, testID).Return(&entity.Test{ID: testID, UserID: userID}, nil)
	testRepo.On("Delete", mock.Anything, testID).Return(nil)

	handler := NewTestHandler(testRepo, new(mockTestDocRepository), new(mockQuestionRepository), new(mockAnswerRepository), new(mockTestUserRepository), nil, nil)
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error { c.Locals("userID", userID); return c.Next() })
	app.Delete("/tests/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/tests/"+testID.String(), nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	testRepo.AssertExpectations(t)
}

// NEW TESTS FOR UPDATED FUNCTIONALITY

// TestGetByID_Success_WithQuestionsAndAnswers tests the updated GetByID handler
// that returns test with questions and answers
func TestGetByID_Success_WithQuestionsAndAnswers(t *testing.T) {
	userID := uuid.New()
	testID := uuid.New()
	questionID1 := uuid.New()
	questionID2 := uuid.New()

	testRepo := new(mockTestRepository)
	questionRepo := new(mockQuestionRepository)
	answerRepo := new(mockAnswerRepository)

	// Mock test
	test := &entity.Test{
		ID:             testID,
		UserID:         userID,
		Title:          "Sample Test",
		Description:    "Test Description",
		TotalQuestions: 2,
		Status:         entity.TestStatusDraft,
		MoodleSynced:   false,
	}
	testRepo.On("FindByID", mock.Anything, testID).Return(test, nil)

	// Mock questions
	questions := []*entity.Question{
		{
			ID:           questionID1,
			TestID:       testID,
			QuestionText: "What is Go?",
			QuestionType: entity.QuestionTypeSingleChoice,
			Difficulty:   entity.DifficultyEasy,
			Points:       1.0,
			OrderNum:     1,
		},
		{
			ID:           questionID2,
			TestID:       testID,
			QuestionText: "Is Go statically typed?",
			QuestionType: entity.QuestionTypeTrueFalse,
			Difficulty:   entity.DifficultyMedium,
			Points:       1.0,
			OrderNum:     2,
		},
	}
	questionRepo.On("FindByTestID", mock.Anything, testID).Return(questions, nil)

	// Mock answers for question 1
	answers1 := []*entity.Answer{
		{ID: uuid.New(), QuestionID: questionID1, AnswerText: "A programming language", IsCorrect: true, OrderNum: 1},
		{ID: uuid.New(), QuestionID: questionID1, AnswerText: "A database", IsCorrect: false, OrderNum: 2},
		{ID: uuid.New(), QuestionID: questionID1, AnswerText: "A framework", IsCorrect: false, OrderNum: 3},
		{ID: uuid.New(), QuestionID: questionID1, AnswerText: "An IDE", IsCorrect: false, OrderNum: 4},
	}
	answerRepo.On("FindByQuestionID", mock.Anything, questionID1).Return(answers1, nil)

	// Mock answers for question 2
	answers2 := []*entity.Answer{
		{ID: uuid.New(), QuestionID: questionID2, AnswerText: "True", IsCorrect: true, OrderNum: 1},
		{ID: uuid.New(), QuestionID: questionID2, AnswerText: "False", IsCorrect: false, OrderNum: 2},
	}
	answerRepo.On("FindByQuestionID", mock.Anything, questionID2).Return(answers2, nil)

	handler := NewTestHandler(testRepo, new(mockTestDocRepository), questionRepo, answerRepo, new(mockTestUserRepository), nil, nil)
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error { c.Locals("userID", userID); return c.Next() })
	app.Get("/tests/:id", handler.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/tests/"+testID.String(), nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Verify response structure
	var response dto.TestResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	// Assert test fields
	assert.Equal(t, testID.String(), response.ID)
	assert.Equal(t, "Sample Test", response.Title)
	assert.Equal(t, "Test Description", response.Description)
	assert.Equal(t, 2, response.TotalQuestions)
	assert.Equal(t, "draft", response.Status)
	assert.False(t, response.MoodleSynced)
	assert.NotEmpty(t, response.CreatedAt)

	// Assert questions
	require.Len(t, response.Questions, 2)

	// First question
	q1 := response.Questions[0]
	assert.Equal(t, questionID1.String(), q1.ID)
	assert.Equal(t, "What is Go?", q1.QuestionText)
	assert.Equal(t, "single_choice", q1.QuestionType)
	assert.Equal(t, "easy", q1.Difficulty)
	assert.Equal(t, 1.0, q1.Points)
	assert.Equal(t, 1, q1.OrderNum)
	require.Len(t, q1.Answers, 4)
	assert.True(t, q1.Answers[0].IsCorrect)
	assert.Equal(t, "A programming language", q1.Answers[0].AnswerText)

	// Second question
	q2 := response.Questions[1]
	assert.Equal(t, questionID2.String(), q2.ID)
	assert.Equal(t, "Is Go statically typed?", q2.QuestionText)
	assert.Equal(t, "true_false", q2.QuestionType)
	assert.Equal(t, "medium", q2.Difficulty)
	require.Len(t, q2.Answers, 2)

	testRepo.AssertExpectations(t)
	questionRepo.AssertExpectations(t)
	answerRepo.AssertExpectations(t)
}

// TestGetByID_QuestionsLoadError tests error handling when loading questions fails
func TestGetByID_QuestionsLoadError(t *testing.T) {
	userID := uuid.New()
	testID := uuid.New()

	testRepo := new(mockTestRepository)
	questionRepo := new(mockQuestionRepository)

	test := &entity.Test{ID: testID, UserID: userID, Title: "Test"}
	testRepo.On("FindByID", mock.Anything, testID).Return(test, nil)
	questionRepo.On("FindByTestID", mock.Anything, testID).Return(nil, assert.AnError)

	handler := NewTestHandler(testRepo, new(mockTestDocRepository), questionRepo, new(mockAnswerRepository), new(mockTestUserRepository), nil, nil)
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error { c.Locals("userID", userID); return c.Next() })
	app.Get("/tests/:id", handler.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/tests/"+testID.String(), nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response dto.ErrorResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response.Error.Message, "failed to load questions")
}

// TestGetByID_AnswersLoadError tests error handling when loading answers fails
func TestGetByID_AnswersLoadError(t *testing.T) {
	userID := uuid.New()
	testID := uuid.New()
	questionID := uuid.New()

	testRepo := new(mockTestRepository)
	questionRepo := new(mockQuestionRepository)
	answerRepo := new(mockAnswerRepository)

	test := &entity.Test{ID: testID, UserID: userID, Title: "Test"}
	testRepo.On("FindByID", mock.Anything, testID).Return(test, nil)

	questions := []*entity.Question{
		{ID: questionID, TestID: testID, QuestionText: "Question 1", OrderNum: 1},
	}
	questionRepo.On("FindByTestID", mock.Anything, testID).Return(questions, nil)
	answerRepo.On("FindByQuestionID", mock.Anything, questionID).Return(nil, assert.AnError)

	handler := NewTestHandler(testRepo, new(mockTestDocRepository), questionRepo, answerRepo, new(mockTestUserRepository), nil, nil)
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error { c.Locals("userID", userID); return c.Next() })
	app.Get("/tests/:id", handler.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/tests/"+testID.String(), nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response dto.ErrorResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response.Error.Message, "failed to load answers")
}

// TestListTests_ReturnsCompleteData tests the updated List handler
// that returns complete test data with all fields
func TestListTests_ReturnsCompleteData(t *testing.T) {
	userID := uuid.New()
	testID1 := uuid.New()
	testID2 := uuid.New()

	testRepo := new(mockTestRepository)
	userRepo := new(mockTestUserRepository)

	// Mock user with teacher role
	teacherRole := &entity.Role{ID: uuid.New(), Name: "teacher"}
	teacherUser := &entity.User{ID: userID, Email: "teacher@test.com", RoleID: teacherRole.ID, Role: teacherRole}
	userRepo.On("FindByID", mock.Anything, userID).Return(teacherUser, nil)

	tests := []*entity.Test{
		{
			ID:             testID1,
			UserID:         userID,
			Title:          "Test 1",
			Description:    "Description 1",
			TotalQuestions: 10,
			Status:         entity.TestStatusDraft,
			MoodleSynced:   false,
		},
		{
			ID:             testID2,
			UserID:         userID,
			Title:          "Test 2",
			Description:    "Description 2",
			TotalQuestions: 5,
			Status:         entity.TestStatusPublished,
			MoodleSynced:   true,
		},
	}

	testRepo.On("FindByUserID", mock.Anything, userID, 20, 0).Return(tests, nil)
	testRepo.On("CountByUserID", mock.Anything, userID).Return(int64(2), nil)

	handler := NewTestHandler(testRepo, new(mockTestDocRepository), new(mockQuestionRepository), new(mockAnswerRepository), userRepo, nil, nil)
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error { c.Locals("userID", userID); return c.Next() })
	app.Get("/tests", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/tests", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Verify response structure
	var response dto.TestListResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, int64(2), response.Total)
	assert.Equal(t, 1, response.Page)
	assert.Equal(t, 20, response.PageSize)
	require.Len(t, response.Tests, 2)

	// Verify first test
	test1 := response.Tests[0]
	assert.Equal(t, testID1.String(), test1.ID)
	assert.Equal(t, "Test 1", test1.Title)
	assert.Equal(t, "Description 1", test1.Description)
	assert.Equal(t, 10, test1.TotalQuestions)
	assert.Equal(t, "draft", test1.Status)
	assert.False(t, test1.MoodleSynced)
	assert.NotEmpty(t, test1.CreatedAt)

	// Verify second test
	test2 := response.Tests[1]
	assert.Equal(t, testID2.String(), test2.ID)
	assert.Equal(t, "Test 2", test2.Title)
	assert.Equal(t, "Description 2", test2.Description)
	assert.Equal(t, 5, test2.TotalQuestions)
	assert.Equal(t, "published", test2.Status)
	assert.True(t, test2.MoodleSynced)
	assert.NotEmpty(t, test2.CreatedAt)

	testRepo.AssertExpectations(t)
}

// TestListTests_WithPagination tests pagination parameters
func TestListTests_WithPagination(t *testing.T) {
	userID := uuid.New()
	testRepo := new(mockTestRepository)
	userRepo := new(mockTestUserRepository)

	// Mock user with teacher role
	teacherRole := &entity.Role{ID: uuid.New(), Name: "teacher"}
	teacherUser := &entity.User{ID: userID, Email: "teacher@test.com", RoleID: teacherRole.ID, Role: teacherRole}
	userRepo.On("FindByID", mock.Anything, userID).Return(teacherUser, nil)

	// Request page 2 with page_size 10
	testRepo.On("FindByUserID", mock.Anything, userID, 10, 10).Return([]*entity.Test{}, nil)
	testRepo.On("CountByUserID", mock.Anything, userID).Return(int64(25), nil)

	handler := NewTestHandler(testRepo, new(mockTestDocRepository), new(mockQuestionRepository), new(mockAnswerRepository), userRepo, nil, nil)
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error { c.Locals("userID", userID); return c.Next() })
	app.Get("/tests", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/tests?page=2&page_size=10", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response dto.TestListResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, int64(25), response.Total)
	assert.Equal(t, 2, response.Page)
	assert.Equal(t, 10, response.PageSize)

	testRepo.AssertExpectations(t)
}
