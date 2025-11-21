package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/application/dto"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/moodle"
)

// MoodleHandler handles Moodle integration operations
type MoodleHandler struct {
	testRepo     repository.TestRepository
	questionRepo repository.QuestionRepository
	answerRepo   repository.AnswerRepository
	xmlExporter  *moodle.MoodleXMLExporter
	moodleClient *moodle.Client
}

// NewMoodleHandler creates a new Moodle handler
func NewMoodleHandler(
	testRepo repository.TestRepository,
	questionRepo repository.QuestionRepository,
	answerRepo repository.AnswerRepository,
	xmlExporter *moodle.MoodleXMLExporter,
	moodleClient *moodle.Client,
) *MoodleHandler {
	return &MoodleHandler{
		testRepo:     testRepo,
		questionRepo: questionRepo,
		answerRepo:   answerRepo,
		xmlExporter:  xmlExporter,
		moodleClient: moodleClient,
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
// @Router /tests/{id}/export-xml [get]
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

	// Get test
	test, err := h.testRepo.FindByID(c.Context(), testID)
	if err != nil || test.UserID != userID {
		return c.Status(fiber.StatusNotFound).JSON(
			dto.NewErrorResponse(dto.ErrCodeTestNotFound, "test not found"),
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
	c.Set("Content-Type", "application/xml")
	c.Set("Content-Disposition", "attachment; filename="+test.Title+".xml")

	return c.SendString(xmlContent)
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
// @Router /tests/{id}/sync-moodle [post]
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

	// Get test
	test, err := h.testRepo.FindByID(c.Context(), testID)
	if err != nil || test.UserID != userID {
		return c.Status(fiber.StatusNotFound).JSON(
			dto.NewErrorResponse(dto.ErrCodeTestNotFound, "test not found"),
		)
	}

	// Get questions
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

	// Get answers
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

	// Upload to Moodle
	uploadResp, err := h.moodleClient.UploadQuiz(c.Context(), moodle.UploadQuizRequest{
		CourseName: req.CourseName,
		QuizName:   test.Title,
		XMLContent: xmlContent,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeMoodleSyncFailed, "failed to sync with Moodle: "+err.Error()),
		)
	}

	if !uploadResp.Success {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeMoodleUploadFailed, "moodle upload failed: "+uploadResp.Message),
		)
	}

	// Update test with Moodle sync info
	test.MarkMoodleSynced(uploadResp.QuizID)
	if err := h.testRepo.Update(c.Context(), test); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeDatabaseError, "failed to update test"),
		)
	}

	return c.JSON(dto.MoodleSyncResponse{
		Message:  "test synchronized with Moodle",
		MoodleID: uploadResp.QuizID,
		CourseID: uploadResp.CourseID,
	})
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
	courses, err := h.moodleClient.GetCourses(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(dto.ErrCodeMoodleNotConnected, "failed to retrieve Moodle courses: "+err.Error()),
		)
	}

	// Convert to DTOs
	courseDTOs := make([]dto.MoodleCourse, len(courses))
	for i, course := range courses {
		courseDTOs[i] = dto.MoodleCourse{
			ID:        fmt.Sprintf("%d", course.ID),
			Name:      course.FullName,
			ShortName: course.ShortName,
		}
	}

	return c.JSON(dto.MoodleCoursesResponse{
		Courses: courseDTOs,
	})
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
	if err := h.moodleClient.ValidateConnection(c.Context()); err != nil {
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
