package middleware

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
)

// ActivityLogger is a middleware factory that records user actions to the activity_logs table.
type ActivityLogger struct {
	repo repository.ActivityLogRepository
}

func NewActivityLogger(repo repository.ActivityLogRepository) *ActivityLogger {
	return &ActivityLogger{repo: repo}
}

// Log returns a Fiber middleware that logs the given action on successful responses (HTTP 2xx).
// entityType is the type of entity being acted on (e.g. "document", "test", "user").
// entityID is extracted from route params: ":id", ":testId", or ":questionId".
func (al *ActivityLogger) Log(action, entityType string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		// Only log successful operations
		if c.Response().StatusCode() >= 400 {
			return err
		}

		log := &entity.ActivityLog{
			Action:     action,
			EntityType: entityType,
			IPAddress:  c.IP(),
			UserAgent:  c.Get("User-Agent"),
		}

		// Extract authenticated user ID (nil for public routes like login/register)
		if uid, ok := extractUserID(c); ok {
			log.UserID = &uid
		}

		// Extract entity ID from common route params
		if rawID := c.Params("id"); rawID != "" {
			if entityID, parseErr := uuid.Parse(rawID); parseErr == nil {
				log.EntityID = &entityID
			}
		} else if rawID := c.Params("questionId"); rawID != "" {
			if entityID, parseErr := uuid.Parse(rawID); parseErr == nil {
				log.EntityID = &entityID
			}
		} else if rawID := c.Params("testId"); rawID != "" {
			if entityID, parseErr := uuid.Parse(rawID); parseErr == nil {
				log.EntityID = &entityID
			}
		}

		go func() {
			_ = al.repo.Create(context.Background(), log)
		}()

		return err
	}
}

func extractUserID(c *fiber.Ctx) (uuid.UUID, bool) {
	raw := c.Locals("userID")
	if raw == nil {
		return uuid.Nil, false
	}
	switch v := raw.(type) {
	case uuid.UUID:
		if v == uuid.Nil {
			return uuid.Nil, false
		}
		return v, true
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil || parsed == uuid.Nil {
			return uuid.Nil, false
		}
		return parsed, true
	default:
		return uuid.Nil, false
	}
}
