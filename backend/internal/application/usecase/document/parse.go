package document

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/shester1kov/testgen-backend/internal/domain/repository"
	"github.com/shester1kov/testgen-backend/internal/infrastructure/parser"
)

// ParseUseCase handles document parsing
type ParseUseCase struct {
	documentRepo  repository.DocumentRepository
	parserFactory *parser.DocumentParserFactory
}

// NewParseUseCase creates a new parse use case
func NewParseUseCase(documentRepo repository.DocumentRepository, parserFactory *parser.DocumentParserFactory) *ParseUseCase {
	return &ParseUseCase{
		documentRepo:  documentRepo,
		parserFactory: parserFactory,
	}
}

// Execute executes the parse use case
func (uc *ParseUseCase) Execute(ctx context.Context, documentID uuid.UUID, userID uuid.UUID) error {
	// Get document from database
	document, err := uc.documentRepo.FindByID(ctx, documentID)
	if err != nil {
		return fmt.Errorf("document not found: %w", err)
	}

	// Verify ownership
	if document.UserID != userID {
		return fmt.Errorf("unauthorized access to document")
	}

	// Check if already parsed
	if document.IsParsed() {
		return fmt.Errorf("document already parsed")
	}

	// Mark as parsing
	document.MarkAsParsing()
	if err := uc.documentRepo.Update(ctx, document); err != nil {
		return fmt.Errorf("failed to update document status: %w", err)
	}

	// Get appropriate parser
	docParser, err := uc.parserFactory.CreateParser(string(document.FileType))
	if err != nil {
		document.MarkAsError()
		uc.documentRepo.Update(ctx, document)
		return fmt.Errorf("failed to create parser: %w", err)
	}

	// Open file
	file, err := os.Open(document.FilePath)
	if err != nil {
		document.MarkAsError()
		uc.documentRepo.Update(ctx, document)
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Parse document
	parsedText, err := docParser.Parse(file)
	if err != nil {
		document.MarkAsError()
		uc.documentRepo.Update(ctx, document)
		return fmt.Errorf("failed to parse document: %w", err)
	}

	// Update document with parsed text
	document.MarkAsParsed(parsedText)
	if err := uc.documentRepo.Update(ctx, document); err != nil {
		return fmt.Errorf("failed to save parsed text: %w", err)
	}

	return nil
}
