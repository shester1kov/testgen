package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/application/dto"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
)

type StatsHandler struct {
	testRepo     repository.TestRepository
	documentRepo repository.DocumentRepository
	questionRepo repository.QuestionRepository
	userRepo     repository.UserRepository
}

func NewStatsHandler(
	testRepo repository.TestRepository,
	documentRepo repository.DocumentRepository,
	questionRepo repository.QuestionRepository,
	userRepo repository.UserRepository,
) *StatsHandler {
	return &StatsHandler{
		testRepo:     testRepo,
		documentRepo: documentRepo,
		questionRepo: questionRepo,
		userRepo:     userRepo,
	}
}

// GetDashboardStats godoc
// @Summary Get dashboard statistics
// @Description Get statistics for the dashboard (documents, tests, questions count)
// @Tags stats
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.DashboardStatsResponse
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /stats/dashboard [get]
func (h *StatsHandler) GetDashboardStats(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(
			dto.NewErrorResponse(dto.ErrCodeUnauthorized, "user not authenticated"),
		)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeInvalidInput, "invalid user ID"),
		)
	}

	// Get user to check role
	user, err := h.userRepo.FindByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeDatabaseError, "failed to fetch user"),
		)
	}

	var documentsCount, testsCount, questionsCount int64

	// Admin sees all data, teacher/student sees only their own
	if user.IsAdmin() {
		// Get total counts for all users
		documentsCount, _ = h.documentRepo.CountAll(c.Context())
		testsCount, _ = h.testRepo.CountAll(c.Context())
		questionsCount, _ = h.questionRepo.CountAll(c.Context())
	} else {
		// Get counts for current user only
		documentsCount, _ = h.documentRepo.CountByUserID(c.Context(), userID)
		testsCount, _ = h.testRepo.CountByUserID(c.Context(), userID)
		questionsCount, _ = h.questionRepo.CountByUserID(c.Context(), userID)
	}

	return c.JSON(dto.DashboardStatsResponse{
		DocumentsCount: documentsCount,
		TestsCount:     testsCount,
		QuestionsCount: questionsCount,
	})
}
