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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock repositories for test handler
type mockTestUpdateRepository struct {
	mock.Mock
}

func (m *mockTestUpdateRepository) Create(ctx context.Context, test *entity.Test) error {
	return nil
}
func (m *mockTestUpdateRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Test, error) {
	args := m.Called(ctx, id)
	if res := args.Get(0); res != nil {
		return res.(*entity.Test), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockTestUpdateRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Test, error) {
	return nil, nil
}
func (m *mockTestUpdateRepository) FindAll(ctx context.Context, limit, offset int) ([]*entity.Test, error) {
	return nil, nil
}
func (m *mockTestUpdateRepository) Update(ctx context.Context, test *entity.Test) error {
	args := m.Called(ctx, test)
	return args.Error(0)
}
func (m *mockTestUpdateRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *mockTestUpdateRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	return 0, nil
}
func (m *mockTestUpdateRepository) CountAll(ctx context.Context) (int64, error) {
	return 0, nil
}

type mockDocumentUpdateRepository struct {
	mock.Mock
}

func (m *mockDocumentUpdateRepository) Create(ctx context.Context, document *entity.Document) error {
	return nil
}
func (m *mockDocumentUpdateRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Document, error) {
	return nil, nil
}
func (m *mockDocumentUpdateRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Document, error) {
	return nil, nil
}
func (m *mockDocumentUpdateRepository) FindAll(ctx context.Context, limit, offset int) ([]*entity.Document, error) {
	return nil, nil
}
func (m *mockDocumentUpdateRepository) Update(ctx context.Context, document *entity.Document) error {
	return nil
}
func (m *mockDocumentUpdateRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *mockDocumentUpdateRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	return 0, nil
}
func (m *mockDocumentUpdateRepository) CountAll(ctx context.Context) (int64, error) {
	return 0, nil
}

type mockQuestionUpdateRepository struct {
	mock.Mock
}

func (m *mockQuestionUpdateRepository) Create(ctx context.Context, question *entity.Question) error {
	return nil
}
func (m *mockQuestionUpdateRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Question, error) {
	args := m.Called(ctx, id)
	if res := args.Get(0); res != nil {
		return res.(*entity.Question), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockQuestionUpdateRepository) FindByTestID(ctx context.Context, testID uuid.UUID) ([]*entity.Question, error) {
	return nil, nil
}
func (m *mockQuestionUpdateRepository) Update(ctx context.Context, question *entity.Question) error {
	args := m.Called(ctx, question)
	return args.Error(0)
}
func (m *mockQuestionUpdateRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *mockQuestionUpdateRepository) CountByTestID(ctx context.Context, testID uuid.UUID) (int, error) {
	return 0, nil
}
func (m *mockQuestionUpdateRepository) ReorderQuestions(ctx context.Context, testID uuid.UUID, questionIDs []uuid.UUID) error {
	return nil
}
func (m *mockQuestionUpdateRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	return 0, nil
}
func (m *mockQuestionUpdateRepository) CountAll(ctx context.Context) (int64, error) {
	return 0, nil
}

type mockAnswerUpdateRepository struct {
	mock.Mock
}

func (m *mockAnswerUpdateRepository) Create(ctx context.Context, answer *entity.Answer) error {
	args := m.Called(ctx, answer)
	return args.Error(0)
}
func (m *mockAnswerUpdateRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Answer, error) {
	return nil, nil
}
func (m *mockAnswerUpdateRepository) FindByQuestionID(ctx context.Context, questionID uuid.UUID) ([]*entity.Answer, error) {
	args := m.Called(ctx, questionID)
	if res := args.Get(0); res != nil {
		return res.([]*entity.Answer), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockAnswerUpdateRepository) Update(ctx context.Context, answer *entity.Answer) error {
	args := m.Called(ctx, answer)
	return args.Error(0)
}
func (m *mockAnswerUpdateRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *mockAnswerUpdateRepository) DeleteByQuestionID(ctx context.Context, questionID uuid.UUID) error {
	return nil
}

type mockUserUpdateRepository struct {
	mock.Mock
}

func (m *mockUserUpdateRepository) Create(ctx context.Context, user *entity.User) error {
	return nil
}
func (m *mockUserUpdateRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	args := m.Called(ctx, id)
	if res := args.Get(0); res != nil {
		return res.(*entity.User), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockUserUpdateRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	return nil, nil
}
func (m *mockUserUpdateRepository) Update(ctx context.Context, user *entity.User) error {
	return nil
}
func (m *mockUserUpdateRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *mockUserUpdateRepository) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	return nil, nil
}
func (m *mockUserUpdateRepository) Count(ctx context.Context) (int64, error) {
	return 0, nil
}

// TestUpdate_Success tests successful test update
func TestUpdate_Success(t *testing.T) {
	userID := uuid.New()
	testID := uuid.New()

	testRepo := new(mockTestUpdateRepository)
	documentRepo := new(mockDocumentUpdateRepository)
	questionRepo := new(mockQuestionUpdateRepository)
	answerRepo := new(mockAnswerUpdateRepository)
	userRepo := new(mockUserUpdateRepository)

	existingTest := &entity.Test{
		ID:          testID,
		UserID:      userID,
		Title:       "Old Title",
		Description: "Old Description",
	}

	testRepo.On("FindByID", mock.Anything, testID).Return(existingTest, nil)
	testRepo.On("Update", mock.Anything, mock.MatchedBy(func(t *entity.Test) bool {
		return t.Title == "New Title" && t.Description == "New Description"
	})).Return(nil)

	handler := NewTestHandler(testRepo, documentRepo, questionRepo, answerRepo, userRepo, nil, nil)
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		return c.Next()
	})
	app.Put("/tests/:id", handler.Update)

	updateReq := dto.UpdateTestRequest{
		Title:       "New Title",
		Description: "New Description",
	}
	body, _ := json.Marshal(updateReq)

	req := httptest.NewRequest(http.MethodPut, "/tests/"+testID.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response dto.TestResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "New Title", response.Title)
	assert.Equal(t, "New Description", response.Description)

	testRepo.AssertExpectations(t)
}

// TestUpdate_NotFound tests updating non-existent test
func TestUpdate_NotFound(t *testing.T) {
	userID := uuid.New()
	testID := uuid.New()

	testRepo := new(mockTestUpdateRepository)
	documentRepo := new(mockDocumentUpdateRepository)
	questionRepo := new(mockQuestionUpdateRepository)
	answerRepo := new(mockAnswerUpdateRepository)
	userRepo := new(mockUserUpdateRepository)

	testRepo.On("FindByID", mock.Anything, testID).Return(nil, assert.AnError)

	handler := NewTestHandler(testRepo, documentRepo, questionRepo, answerRepo, userRepo, nil, nil)
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		return c.Next()
	})
	app.Put("/tests/:id", handler.Update)

	updateReq := dto.UpdateTestRequest{
		Title:       "New Title",
		Description: "New Description",
	}
	body, _ := json.Marshal(updateReq)

	req := httptest.NewRequest(http.MethodPut, "/tests/"+testID.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	testRepo.AssertExpectations(t)
}

// TestUpdate_Unauthorized tests updating test by non-owner
func TestUpdate_Unauthorized(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()
	testID := uuid.New()

	testRepo := new(mockTestUpdateRepository)
	documentRepo := new(mockDocumentUpdateRepository)
	questionRepo := new(mockQuestionUpdateRepository)
	answerRepo := new(mockAnswerUpdateRepository)
	userRepo := new(mockUserUpdateRepository)

	existingTest := &entity.Test{
		ID:     testID,
		UserID: otherUserID, // Different user
		Title:  "Old Title",
	}

	testRepo.On("FindByID", mock.Anything, testID).Return(existingTest, nil)

	handler := NewTestHandler(testRepo, documentRepo, questionRepo, answerRepo, userRepo, nil, nil)
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		return c.Next()
	})
	app.Put("/tests/:id", handler.Update)

	updateReq := dto.UpdateTestRequest{
		Title:       "New Title",
		Description: "New Description",
	}
	body, _ := json.Marshal(updateReq)

	req := httptest.NewRequest(http.MethodPut, "/tests/"+testID.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)

	testRepo.AssertExpectations(t)
}

// TestUpdateQuestion_Success tests successful question update
func TestUpdateQuestion_Success(t *testing.T) {
	userID := uuid.New()
	testID := uuid.New()
	questionID := uuid.New()
	answerID1 := uuid.New()
	answerID2 := uuid.New()

	testRepo := new(mockTestUpdateRepository)
	documentRepo := new(mockDocumentUpdateRepository)
	questionRepo := new(mockQuestionUpdateRepository)
	answerRepo := new(mockAnswerUpdateRepository)
	userRepo := new(mockUserUpdateRepository)

	existingTest := &entity.Test{
		ID:     testID,
		UserID: userID,
		Title:  "Test",
	}

	existingQuestion := &entity.Question{
		ID:           questionID,
		TestID:       testID,
		QuestionText: "Old Question",
		QuestionType: "single_choice",
		Difficulty:   "medium",
		Points:       1.0,
	}

	existingAnswers := []*entity.Answer{
		{
			ID:         answerID1,
			QuestionID: questionID,
			AnswerText: "Old Answer 1",
			IsCorrect:  true,
		},
		{
			ID:         answerID2,
			QuestionID: questionID,
			AnswerText: "Old Answer 2",
			IsCorrect:  false,
		},
	}

	testRepo.On("FindByID", mock.Anything, testID).Return(existingTest, nil)
	questionRepo.On("FindByID", mock.Anything, questionID).Return(existingQuestion, nil)
	answerRepo.On("FindByQuestionID", mock.Anything, questionID).Return(existingAnswers, nil)
	questionRepo.On("Update", mock.Anything, mock.MatchedBy(func(q *entity.Question) bool {
		return q.QuestionText == "New Question"
	})).Return(nil)
	// Mock Delete for old answers
	answerRepo.On("Delete", mock.Anything, answerID1).Return(nil)
	answerRepo.On("Delete", mock.Anything, answerID2).Return(nil)
	// Mock Create for new answers
	answerRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	handler := NewTestHandler(testRepo, documentRepo, questionRepo, answerRepo, userRepo, nil, nil)
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		return c.Next()
	})
	app.Put("/tests/:testId/questions/:questionId", handler.UpdateQuestion)

	updateReq := dto.UpdateQuestionRequest{
		QuestionText: "New Question",
		QuestionType: "single_choice",
		Difficulty:   "hard",
		Points:       ptrFloat64(2.0),
		Answers: []dto.UpdateAnswerRequest{
			{
				ID:         ptrString(answerID1.String()),
				AnswerText: "Updated Answer 1",
				IsCorrect:  true,
				OrderNum:   0,
			},
			{
				ID:         ptrString(answerID2.String()),
				AnswerText: "Updated Answer 2",
				IsCorrect:  false,
				OrderNum:   1,
			},
		},
	}
	body, _ := json.Marshal(updateReq)

	req := httptest.NewRequest(
		http.MethodPut,
		"/tests/"+testID.String()+"/questions/"+questionID.String(),
		bytes.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response dto.QuestionDTO
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "New Question", response.QuestionText)
	assert.Equal(t, "hard", response.Difficulty)
	assert.Equal(t, 2.0, response.Points)

	testRepo.AssertExpectations(t)
	questionRepo.AssertExpectations(t)
	answerRepo.AssertExpectations(t)
}

// TestUpdateQuestion_NotFound tests updating non-existent question
func TestUpdateQuestion_NotFound(t *testing.T) {
	userID := uuid.New()
	testID := uuid.New()
	questionID := uuid.New()

	testRepo := new(mockTestUpdateRepository)
	documentRepo := new(mockDocumentUpdateRepository)
	questionRepo := new(mockQuestionUpdateRepository)
	answerRepo := new(mockAnswerUpdateRepository)
	userRepo := new(mockUserUpdateRepository)

	existingTest := &entity.Test{
		ID:     testID,
		UserID: userID,
		Title:  "Test",
	}

	testRepo.On("FindByID", mock.Anything, testID).Return(existingTest, nil)
	questionRepo.On("FindByID", mock.Anything, questionID).Return(nil, assert.AnError)

	handler := NewTestHandler(testRepo, documentRepo, questionRepo, answerRepo, userRepo, nil, nil)
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		return c.Next()
	})
	app.Put("/tests/:testId/questions/:questionId", handler.UpdateQuestion)

	updateReq := dto.UpdateQuestionRequest{
		QuestionText: "New Question",
	}
	body, _ := json.Marshal(updateReq)

	req := httptest.NewRequest(
		http.MethodPut,
		"/tests/"+testID.String()+"/questions/"+questionID.String(),
		bytes.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	testRepo.AssertExpectations(t)
	questionRepo.AssertExpectations(t)
}

// TestUpdateQuestion_AddNewAnswer tests adding new answers
func TestUpdateQuestion_AddNewAnswer(t *testing.T) {
	userID := uuid.New()
	testID := uuid.New()
	questionID := uuid.New()

	testRepo := new(mockTestUpdateRepository)
	documentRepo := new(mockDocumentUpdateRepository)
	questionRepo := new(mockQuestionUpdateRepository)
	answerRepo := new(mockAnswerUpdateRepository)
	userRepo := new(mockUserUpdateRepository)

	existingTest := &entity.Test{
		ID:     testID,
		UserID: userID,
		Title:  "Test",
	}

	existingQuestion := &entity.Question{
		ID:           questionID,
		TestID:       testID,
		QuestionText: "Question",
		QuestionType: "single_choice",
		Difficulty:   "medium",
		Points:       1.0,
	}

	existingAnswers := []*entity.Answer{}

	testRepo.On("FindByID", mock.Anything, testID).Return(existingTest, nil)
	questionRepo.On("FindByID", mock.Anything, questionID).Return(existingQuestion, nil)
	answerRepo.On("FindByQuestionID", mock.Anything, questionID).Return(existingAnswers, nil)
	questionRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
	answerRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	handler := NewTestHandler(testRepo, documentRepo, questionRepo, answerRepo, userRepo, nil, nil)
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		return c.Next()
	})
	app.Put("/tests/:testId/questions/:questionId", handler.UpdateQuestion)

	updateReq := dto.UpdateQuestionRequest{
		QuestionText: "Question",
		Answers: []dto.UpdateAnswerRequest{
			{
				ID:         nil, // New answer
				AnswerText: "New Answer 1",
				IsCorrect:  true,
				OrderNum:   0,
			},
			{
				ID:         nil, // New answer
				AnswerText: "New Answer 2",
				IsCorrect:  false,
				OrderNum:   1,
			},
		},
	}
	body, _ := json.Marshal(updateReq)

	req := httptest.NewRequest(
		http.MethodPut,
		"/tests/"+testID.String()+"/questions/"+questionID.String(),
		bytes.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	answerRepo.AssertCalled(t, "Create", mock.Anything, mock.Anything)
	testRepo.AssertExpectations(t)
	questionRepo.AssertExpectations(t)
	answerRepo.AssertExpectations(t)
}

// Helper functions
func ptrString(s string) *string {
	return &s
}

func ptrFloat64(f float64) *float64 {
	return &f
}
