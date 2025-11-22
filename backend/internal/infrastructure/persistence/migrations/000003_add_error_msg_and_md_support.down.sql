-- Remove error_msg column from documents table
ALTER TABLE documents DROP COLUMN IF EXISTS error_msg;

-- Revert file_type constraint to original values (without 'md')
ALTER TABLE documents DROP CONSTRAINT IF EXISTS documents_file_type_check;
ALTER TABLE documents ADD CONSTRAINT documents_file_type_check
    CHECK (file_type IN ('pdf', 'docx', 'pptx', 'txt'));
