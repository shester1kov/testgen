package document

import (
	"context"
	"fmt"
	"github.com/shester1kov/testgen-backend/internal/domain"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/application/dto"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
	"go.uber.org/zap"
)

// GetUseCase handles retrieving a single document
type GetUseCase struct {
	documentRepo repository.DocumentRepository
	userRepo     repository.UserRepository
	logger       *zap.Logger
}

// NewGetUseCase creates a new get use case
func NewGetUseCase(documentRepo repository.DocumentRepository, userRepo repository.UserRepository, logger *zap.Logger) *GetUseCase {
	return &GetUseCase{
		documentRepo: documentRepo,
		userRepo:     userRepo,
		logger:       logger,
	}
}

// Execute executes the get use case
// includeParsedText: if true, includes full parsed text in response (can be large!)
func (uc *GetUseCase) Execute(ctx context.Context, documentID uuid.UUID, userID uuid.UUID, includeParsedText bool) (*dto.DocumentUploadResponse, error) {
	uc.logger.Info(
		"Retrieving document",
		zap.String("document_id", documentID.String()),
		zap.String("user_id", userID.String()),
		zap.Bool("include_parsed_text", includeParsedText),
	)

	// Fetch document
	document, err := uc.documentRepo.FindByID(ctx, documentID)
	if err != nil {
		uc.logger.Error("Failed to fetch document",
			zap.Error(err),
			zap.String("document_id", documentID.String()),
		)
		return nil, fmt.Errorf("%w: %v", domain.ErrDocumentNotFound, err)
	}

	// Check if user is admin
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	// Verify ownership
	if document.UserID != userID && !user.IsAdmin() {
		uc.logger.Warn("Access denied: document ownership verification failed",
			zap.String("document_id", documentID.String()),
			zap.String("requesting_user_id", userID.String()),
			zap.String("owner_user_id", document.UserID.String()),
			zap.Bool("is_admin", user.IsAdmin()),
		)
		return nil, domain.ErrAccessDenied
	}

	uc.logger.Info("Document retrieved successfully",
		zap.String("document_id", document.ID.String()),
		zap.String("title", document.Title),
	)

	// Map entity to DTO
	// Only include ParsedText if explicitly requested (to avoid large responses)
	var parsedText *string
	if includeParsedText && document.ParsedText != "" {
		parsedText = &document.ParsedText
	}

	var errorMsg *string
	if document.ErrorMsg != "" {
		errorMsg = &document.ErrorMsg
	}

	return &dto.DocumentUploadResponse{
		ID:         document.ID.String(),
		UserID:     document.UserID.String(),
		Title:      document.Title,
		FileName:   document.FileName,
		FileType:   string(document.FileType),
		FileSize:   document.FileSize,
		ParsedText: parsedText,
		Status:     string(document.Status),
		ErrorMsg:   errorMsg,
		CreatedAt:  document.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}
