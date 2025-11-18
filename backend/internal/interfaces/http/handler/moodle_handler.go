package handler

import (
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

// ExportToXML exports a test to Moodle XML format
func (h *MoodleHandler) ExportToXML(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	testID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid test ID"})
	}

	// Get test
	test, err := h.testRepo.FindByID(c.Context(), testID)
	if err != nil || test.UserID != userID {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "test not found"})
	}

	// Get questions for the test
	questions, err := h.questionRepo.FindByTestID(c.Context(), testID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve questions"})
	}

	if len(questions) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "test has no questions"})
	}

	// Get answers for each question
	answersMap := make(map[string][]*entity.Answer)
	for _, q := range questions {
		answers, err := h.answerRepo.FindByQuestionID(c.Context(), q.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve answers"})
		}
		answersMap[q.ID.String()] = answers
	}

	// Export to XML
	xmlContent, err := h.xmlExporter.Export(test, questions, answersMap)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to export XML"})
	}

	// Set headers for file download
	c.Set("Content-Type", "application/xml")
	c.Set("Content-Disposition", "attachment; filename="+test.Title+".xml")

	return c.SendString(xmlContent)
}

// SyncToMoodle synchronizes a test with Moodle
func (h *MoodleHandler) SyncToMoodle(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	testID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid test ID"})
	}

	var req dto.SyncMoodleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	// Get test
	test, err := h.testRepo.FindByID(c.Context(), testID)
	if err != nil || test.UserID != userID {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "test not found"})
	}

	// Get questions
	questions, err := h.questionRepo.FindByTestID(c.Context(), testID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve questions"})
	}

	if len(questions) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "test has no questions"})
	}

	// Get answers
	answersMap := make(map[string][]*entity.Answer)
	for _, q := range questions {
		answers, err := h.answerRepo.FindByQuestionID(c.Context(), q.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve answers"})
		}
		answersMap[q.ID.String()] = answers
	}

	// Export to XML
	xmlContent, err := h.xmlExporter.Export(test, questions, answersMap)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to export XML"})
	}

	// Upload to Moodle
	uploadResp, err := h.moodleClient.UploadQuiz(c.Context(), moodle.UploadQuizRequest{
		CourseName: req.CourseName,
		QuizName:   test.Title,
		XMLContent: xmlContent,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to sync with Moodle: " + err.Error()})
	}

	if !uploadResp.Success {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "moodle upload failed: " + uploadResp.Message})
	}

	// Update test with Moodle sync info
	test.MarkMoodleSynced(uploadResp.QuizID)
	if err := h.testRepo.Update(c.Context(), test); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update test"})
	}

	return c.JSON(fiber.Map{
		"message":   "test synchronized with Moodle",
		"moodle_id": uploadResp.QuizID,
		"course_id": uploadResp.CourseID,
	})
}

// GetMoodleCourses retrieves available Moodle courses
func (h *MoodleHandler) GetMoodleCourses(c *fiber.Ctx) error {
	courses, err := h.moodleClient.GetCourses(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve Moodle courses: " + err.Error()})
	}

	return c.JSON(fiber.Map{
		"courses": courses,
	})
}

// ValidateMoodleConnection checks if Moodle connection is working
func (h *MoodleHandler) ValidateMoodleConnection(c *fiber.Ctx) error {
	if err := h.moodleClient.ValidateConnection(c.Context()); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"connected": false,
			"error":     err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"connected": true,
		"message":   "Moodle connection is valid",
	})
}
