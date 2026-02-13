-- Migration: Add MinIO support to documents table
-- Changes:
--   1. Rename file_path → object_key (semantic clarity)
--   2. Add storage_type column (minio/local)
--   3. Add content_type column (MIME type)
--   4. Add index on storage_type
-- Date: 2026-02-12

-- ============================================
-- 1. Rename column: file_path → object_key
-- ============================================
-- Причина: object_key более точно отражает назначение
-- Примеры:
--   - "documents/550e8400-e29b-41d4-a716-446655440000/file.pdf"
--   - "documents/550e8400-e29b-41d4-a716-446655440000/presentation.pptx"
ALTER TABLE documents
RENAME COLUMN file_path TO object_key;

-- Обновляем комментарий к колонке
COMMENT ON COLUMN documents.object_key IS 'MinIO object key (e.g., documents/{user-id}/{uuid}.{ext})';

-- ============================================
-- 2. Add content_type column
-- ============================================
-- MIME-тип файла для правильной отдачи через HTTP
-- Будет заполняться автоматически при загрузке
ALTER TABLE documents
ADD COLUMN content_type VARCHAR(100);

COMMENT ON COLUMN documents.content_type IS 'MIME content type for HTTP Content-Type header (e.g., application/pdf)';


-- ============================================
-- Migration complete
-- ============================================
-- Note: Existing file_path values will be preserved as object_key
-- You need to manually migrate files from filesystem to MinIO
-- and update object_key values to MinIO format