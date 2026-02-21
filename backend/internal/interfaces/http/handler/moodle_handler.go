package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/application/dto"
	moodleuc "github.com/shester1kov/testgen-backend/internal/application/usecase/moodle"
)

// MoodleHandler handles Moodle integration operations
type MoodleHandler struct {
	exportXMLUseCase          *moodleuc.ExportXMLUseCase
	syncUseCase               *moodleuc.SyncUseCase
	getCoursesUseCase         *moodleuc.GetCoursesUseCase
	validateConnectionUseCase *moodleuc.ValidateConnectionUseCase
}

// NewMoodleHandler creates a new Moodle handler
func NewMoodleHandler(
	exportXMLUseCase *moodleuc.ExportXMLUseCase,
	syncUseCase *moodleuc.SyncUseCase,
	getCoursesUseCase *moodleuc.GetCoursesUseCase,
	validateConnectionUseCase *moodleuc.ValidateConnectionUseCase,
) *MoodleHandler {
	return &MoodleHandler{
		exportXMLUseCase:          exportXMLUseCase,
		syncUseCase:               syncUseCase,
		getCoursesUseCase:         getCoursesUseCase,
		validateConnectionUseCase: validateConnectionUseCase,
	}
}

// ExportToXML godoc
// @Summary Export test to Moodle XML format
// @Description Export a test and its questions to Moodle XML format for download
// @Tags moodle
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "Test ID"
// @Success 200 {string} string "XML file"
// @Failure 400 {object} dto.ErrorResponse "Invalid test ID or test has no questions"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 404 {object} dto.ErrorResponse "Test not found"
// @Failure 500 {object} dto.ErrorResponse "Export failed"
// @Router /moodle/tests/{id}/export [get]
func (h *MoodleHandler) ExportToXML(c *fiber.Ctx) error {
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

	result, err := h.exportXMLUseCase.Execute(c.Context(), testID, userID)
	if err != nil {
		if strings.Contains(err.Error(), "test not found") {
			return c.Status(fiber.StatusNotFound).JSON(
				dto.NewErrorResponse(dto.ErrCodeTestNotFound, "test not found"),
			)
		}
		if strings.Contains(err.Error(), "no questions") {
			return c.Status(fiber.StatusBadRequest).JSON(
				dto.NewErrorResponse(dto.ErrCodeTestHasNoQuestions, "test has no questions"),
			)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeExportFailed, "failed to export XML"),
		)
	}

	c.Set("Content-Type", result.ContentType)
	c.Set("Content-Disposition", "attachment; filename="+result.Title+"."+result.FileExt)

	return c.SendString(result.Content)
}

// SyncToMoodle godoc
// @Summary Sync test to Moodle
// @Description Synchronize a test with Moodle by uploading it as a quiz to a specified course
// @Tags moodle
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Test ID"
// @Param request body dto.SyncMoodleRequest true "Sync request with course name"
// @Success 200 {object} dto.MoodleSyncResponse
// @Failure 400 {object} dto.ErrorResponse "Invalid input or test has no questions"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 404 {object} dto.ErrorResponse "Test not found"
// @Failure 500 {object} dto.ErrorResponse "Sync failed"
// @Router /moodle/tests/{id}/sync [post]
func (h *MoodleHandler) SyncToMoodle(c *fiber.Ctx) error {
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

	var req dto.SyncMoodleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse(dto.ErrCodeInvalidInput, "invalid request"),
		)
	}

	result, err := h.syncUseCase.Execute(c.Context(), moodleuc.SyncParams{
		TestID:     testID,
		UserID:     userID,
		CourseName: req.CourseName,
	})
	if err != nil {
		if strings.Contains(err.Error(), "test not found") {
			return c.Status(fiber.StatusNotFound).JSON(
				dto.NewErrorResponse(dto.ErrCodeTestNotFound, "test not found"),
			)
		}
		if strings.Contains(err.Error(), "no questions") {
			return c.Status(fiber.StatusBadRequest).JSON(
				dto.NewErrorResponse(dto.ErrCodeTestHasNoQuestions, "test has no questions"),
			)
		}
		if strings.Contains(err.Error(), "failed to sync with Moodle") {
			return c.Status(fiber.StatusInternalServerError).JSON(
				dto.NewErrorResponse(dto.ErrCodeMoodleSyncFailed, err.Error()),
			)
		}
		if strings.Contains(err.Error(), "moodle upload failed") {
			return c.Status(fiber.StatusInternalServerError).JSON(
				dto.NewErrorResponse(dto.ErrCodeMoodleUploadFailed, err.Error()),
			)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeInternalError, "failed to sync test"),
		)
	}

	return c.JSON(result)
}

// GetMoodleCourses godoc
// @Summary Get Moodle courses
// @Description Retrieve list of available courses from Moodle
// @Tags moodle
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.MoodleCoursesResponse
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Failed to retrieve courses"
// @Router /moodle/courses [get]
func (h *MoodleHandler) GetMoodleCourses(c *fiber.Ctx) error {
	result, err := h.getCoursesUseCase.Execute(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeMoodleNotConnected, err.Error()),
		)
	}

	return c.JSON(result)
}

// ValidateMoodleConnection godoc
// @Summary Validate Moodle connection
// @Description Check if the Moodle server connection is working properly
// @Tags moodle
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.MoodleConnectionResponse
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 503 {object} dto.MoodleConnectionResponse "Connection failed"
// @Router /moodle/validate [get]
func (h *MoodleHandler) ValidateMoodleConnection(c *fiber.Ctx) error {
	if err := h.validateConnectionUseCase.Execute(c.Context()); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(dto.MoodleConnectionResponse{
			Connected: false,
			Error:     err.Error(),
		})
	}

	return c.JSON(dto.MoodleConnectionResponse{
		Connected: true,
		Message:   "Moodle connection is valid",
	})
}
