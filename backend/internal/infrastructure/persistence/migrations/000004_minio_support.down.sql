-- Rollback: Remove MinIO support from documents table
-- This migration reverts changes made in 000004_minio_support.up.sql
-- Date: 2026-02-12

-- ============================================
-- 1. Drop content_type column
-- ============================================
ALTER TABLE documents
DROP COLUMN IF EXISTS content_type;


-- ============================================
-- 2. Rename column back: object_key → file_path
-- ============================================
ALTER TABLE documents
RENAME COLUMN object_key TO file_path;

-- Восстанавливаем старый комментарий
COMMENT ON COLUMN documents.file_path IS 'File path to uploaded document';

-- ============================================
-- Rollback complete
-- ============================================