-- Add error_msg column to documents table
ALTER TABLE documents ADD COLUMN error_msg TEXT;

-- Update file_type constraint to include 'md' format
ALTER TABLE documents DROP CONSTRAINT IF EXISTS documents_file_type_check;
ALTER TABLE documents ADD CONSTRAINT documents_file_type_check
    CHECK (file_type IN ('pdf', 'docx', 'pptx', 'txt', 'md'));
