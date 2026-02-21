package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shester1kov/testgen-backend/internal/application/dto"
	"github.com/shester1kov/testgen-backend/internal/application/usecase/stats"
)

type StatsHandler struct {
	dashboardUseCase *stats.DashboardUseCase
}

func NewStatsHandler(dashboardUseCase *stats.DashboardUseCase) *StatsHandler {
	return &StatsHandler{
		dashboardUseCase: dashboardUseCase,
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
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(
			dto.NewErrorResponse(dto.ErrCodeUnauthorized, "user not authenticated"),
		)
	}

	result, err := h.dashboardUseCase.Execute(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeDatabaseError, "failed to fetch statistics"),
		)
	}

	return c.JSON(result)
}
