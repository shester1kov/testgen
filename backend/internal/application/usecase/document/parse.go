package document

import (
	"context"
	"fmt"
	"github.com/shester1kov/testgen-backend/internal/domain"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/parser"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/storage"
	"go.uber.org/zap"
)

// ParseUseCase handles document parsing
type ParseUseCase struct {
	documentRepo  repository.DocumentRepository
	storage       storage.Storage
	parserFactory *parser.DocumentParserFactory
	userRepo      repository.UserRepository
	logger        *zap.Logger
}

// NewParseUseCase creates a new parse use case
func NewParseUseCase(
	documentRepo repository.DocumentRepository,
	storage storage.Storage,
	parserFactory *parser.DocumentParserFactory,
	userRepo repository.UserRepository,
	logger *zap.Logger,
) *ParseUseCase {
	return &ParseUseCase{
		documentRepo:  documentRepo,
		storage:       storage,
		parserFactory: parserFactory,
		userRepo:      userRepo,
		logger:        logger,
	}
}

// Execute executes the parse use case
func (uc *ParseUseCase) Execute(
	ctx context.Context,
	documentID uuid.UUID,
	userID uuid.UUID,
) error {
	uc.logger.Info(
		"Starting document parsing",
		zap.String("document_id", documentID.String()),
		zap.String("user_id", userID.String()),
	)

	// Get document from database
	document, err := uc.documentRepo.FindByID(ctx, documentID)
	if err != nil {
		uc.logger.Error("Document not found",
			zap.Error(err),
			zap.String("document_id", documentID.String()),
		)
		return fmt.Errorf("%w: %v", domain.ErrDocumentNotFound, err)
	}

	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to fetch user: %w", err)
	}

	// Verify ownership
	if document.UserID != userID && !user.IsAdmin() {
		uc.logger.Warn("Unauthorized access attempt",
			zap.String("document_id", documentID.String()),
			zap.String("requesting_user_id", userID.String()),
			zap.String("owner_user_id", document.UserID.String()),
		)
		return domain.ErrUnauthorized
	}

	// Check if already parsed
	if document.IsParsed() {
		uc.logger.Warn("Document already parsed",
			zap.String("document_id", documentID.String()),
		)
		return domain.ErrAlreadyParsed
	}

	// Mark as parsing
	document.MarkAsParsing()
	if err := uc.documentRepo.Update(ctx, document); err != nil {
		uc.logger.Error("Failed to update document status",
			zap.Error(err),
			zap.String("document_id", documentID.String()),
		)
		return fmt.Errorf("failed to update document status: %w", err)
	}

	// Get appropriate parser
	docParser, err := uc.parserFactory.CreateParser(string(document.FileType))
	if err != nil {
		uc.logger.Error("Failed to create parser",
			zap.Error(err),
			zap.String("file_type", string(document.FileType)),
		)
		document.MarkAsError(err.Error())
		uc.documentRepo.Update(ctx, document)
		return fmt.Errorf("failed to create parser: %w", err)
	}

	// Download file from storage
	uc.logger.Info("Downloading file from storage",
		zap.String("object_key", document.ObjectKey),
	)

	file, err := uc.storage.Download(ctx, document.ObjectKey)
	if err != nil {
		uc.logger.Error("Failed to download file from storage",
			zap.Error(err),
			zap.String("object_key", document.ObjectKey),
		)
		document.MarkAsError(err.Error())
		uc.documentRepo.Update(ctx, document)
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer file.Close()

	// Parse document
	uc.logger.Info("Parsing document content",
		zap.String("document_id", documentID.String()),
		zap.String("file_type", string(document.FileType)),
	)

	parsedText, err := docParser.Parse(file)
	if err != nil {
		uc.logger.Error("Failed to parse document",
			zap.Error(err),
			zap.String("document_id", documentID.String()),
		)
		document.MarkAsError(err.Error())
		uc.documentRepo.Update(ctx, document)
		return fmt.Errorf("failed to parse document: %w", err)
	}

	// Update document with parsed text
	document.MarkAsParsed(parsedText)
	if err := uc.documentRepo.Update(ctx, document); err != nil {
		uc.logger.Error("Failed to save parsed text",
			zap.Error(err),
			zap.String("document_id", documentID.String()),
		)
		return fmt.Errorf("failed to save parsed text: %w", err)
	}

	uc.logger.Info("Document parsed successfully",
		zap.String("document_id", documentID.String()),
		zap.Int("parsed_text_length", len(parsedText)),
	)

	return nil
}
