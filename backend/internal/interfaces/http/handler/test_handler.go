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

func (h *TestHandler) Create(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	var req dto.CreateTestRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	test := &entity.Test{UserID: userID, Title: req.Title, Description: req.Description}
	if req.DocumentID != nil {
		docID, _ := uuid.Parse(*req.DocumentID)
		test.DocumentID = &docID
	}

	if err := h.testRepo.Create(c.Context(), test); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create test"})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.TestResponse{
		ID: test.ID.String(), Title: test.Title, Description: test.Description,
		TotalQuestions: 0, Status: string(test.Status), CreatedAt: test.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

func (h *TestHandler) Generate(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	var req dto.GenerateTestRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	docID, _ := uuid.Parse(req.DocumentID)
	document, err := h.documentRepo.FindByID(c.Context(), docID)
	if err != nil || document.UserID != userID {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "document not found"})
	}

	if !document.IsParsed() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "document not parsed yet"})
	}

	provider := req.LLMProvider
	if provider == "" {
		provider = "perplexity"
	}

	strategy, err := h.llmFactory.CreateStrategy(provider)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	llmContext := llm.NewLLMContext(strategy)
	questions, err := llmContext.GenerateQuestions(c.Context(), llm.GenerationParams{
		Text: document.ParsedText, NumQuestions: req.NumQuestions, Difficulty: req.Difficulty,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate questions"})
	}

	return c.JSON(fiber.Map{"message": "questions generated", "count": len(questions), "questions": questions})
}

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

func (h *TestHandler) GetByID(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	testID, _ := uuid.Parse(c.Params("id"))
	test, err := h.testRepo.FindByID(c.Context(), testID)
	if err != nil || test.UserID != userID {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "test not found"})
	}

	return c.JSON(dto.TestResponse{ID: test.ID.String(), Title: test.Title, Description: test.Description, TotalQuestions: test.TotalQuestions})
}

func (h *TestHandler) Delete(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	testID, _ := uuid.Parse(c.Params("id"))
	test, err := h.testRepo.FindByID(c.Context(), testID)
	if err != nil || test.UserID != userID {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "test not found"})
	}

	h.testRepo.Delete(c.Context(), testID)
	return c.JSON(fiber.Map{"message": "test deleted"})
}
