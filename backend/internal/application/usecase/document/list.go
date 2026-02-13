package document

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/application/dto"
	"github.com/shester1kov/testgen-backend/internal/domain/entity"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
	"go.uber.org/zap"
)

// ListUseCase handles document listing
type ListUseCase struct {
	documentRepo repository.DocumentRepository
	userRepo     repository.UserRepository
	logger       *zap.Logger
}

// NewListUseCase creates a new list use case
func NewListUseCase(
	documentRepo repository.DocumentRepository,
	userRepo repository.UserRepository,
	logger *zap.Logger,
) *ListUseCase {
	return &ListUseCase{
		documentRepo: documentRepo,
		userRepo:     userRepo,
		logger:       logger,
	}
}

// Execute executes the list use case
func (uc *ListUseCase) Execute(ctx context.Context, userID uuid.UUID, page, pageSize int) (*dto.DocumentListResponse, error) {
	uc.logger.Info(
		"Listing documents",
		zap.String("user_id", userID.String()),
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
	)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	var documents []*entity.Document
	var total int64

	// Get user to check role
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user")
	}

	// Admin sees all documents, others see only their own
	if user.IsAdmin() {
		documents, err = uc.documentRepo.FindAll(ctx, pageSize, offset)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch documents")
		}

		total, err = uc.documentRepo.CountAll(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to count documents")
		}
	} else {
		documents, err = uc.documentRepo.FindByUserID(ctx, userID, pageSize, offset)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch documents")
		}

		total, err = uc.documentRepo.CountByUserID(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to count documents")
		}
	}

	result := make([]dto.DocumentUploadResponse, len(documents))
	for i, doc := range documents {
		var parsedText *string
		if doc.ParsedText != "" {
			parsedText = &doc.ParsedText
		}
		var errorMsg *string
		if doc.ErrorMsg != "" {
			errorMsg = &doc.ErrorMsg
		}

		// Include user info for admin
		var userName *string
		var userEmail *string
		if user.IsAdmin() && doc.User.ID != uuid.Nil {
			userName = &doc.User.FullName
			userEmail = &doc.User.Email
		}

		result[i] = dto.DocumentUploadResponse{
			ID:         doc.ID.String(),
			UserID:     doc.UserID.String(),
			UserName:   userName,
			UserEmail:  userEmail,
			Title:      doc.Title,
			FileName:   doc.FileName,
			FileType:   string(doc.FileType),
			FileSize:   doc.FileSize,
			ParsedText: parsedText,
			Status:     string(doc.Status),
			ErrorMsg:   errorMsg,
			CreatedAt:  doc.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	uc.logger.Info("Documents listed successfully",
		zap.String("user_id", userID.String()),
		zap.Int("count", len(documents)),
		zap.Int64("total", total),
	)

	return &dto.DocumentListResponse{
		Documents: result,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
	}, nil
}
