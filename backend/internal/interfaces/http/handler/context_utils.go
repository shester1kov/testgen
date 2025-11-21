package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// getUserIDFromContext safely extracts the authenticated user ID from Fiber context.
// It supports both uuid.UUID and string representations and returns false when the
// value is missing or invalid to avoid panics on type assertions.
func getUserIDFromContext(c *fiber.Ctx) (uuid.UUID, bool) {
	rawUserID := c.Locals("userID")
	if rawUserID == nil {
		return uuid.Nil, false
	}

	switch v := rawUserID.(type) {
	case uuid.UUID:
		if v == uuid.Nil {
			return uuid.Nil, false
		}
		return v, true
	case string:
		parsedID, err := uuid.Parse(v)
		if err != nil || parsedID == uuid.Nil {
			return uuid.Nil, false
		}
		return parsedID, true
	default:
		return uuid.Nil, false
	}
}
