-- Drop triggers
DROP TRIGGER IF EXISTS update_questions_updated_at ON questions;
DROP TRIGGER IF EXISTS update_tests_updated_at ON tests;
DROP TRIGGER IF EXISTS update_documents_updated_at ON documents;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_activity_logs_created_at;
DROP INDEX IF EXISTS idx_activity_logs_user_id;
DROP INDEX IF EXISTS idx_answers_question_id;
DROP INDEX IF EXISTS idx_questions_order;
DROP INDEX IF EXISTS idx_questions_test_id;
DROP INDEX IF EXISTS idx_tests_status;
DROP INDEX IF EXISTS idx_tests_document_id;
DROP INDEX IF EXISTS idx_tests_user_id;
DROP INDEX IF EXISTS idx_documents_status;
DROP INDEX IF EXISTS idx_documents_user_id;

-- Drop tables
DROP TABLE IF EXISTS activity_logs;
DROP TABLE IF EXISTS answers;
DROP TABLE IF EXISTS questions;
DROP TABLE IF EXISTS tests;
DROP TABLE IF EXISTS documents;
DROP TABLE IF EXISTS users;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";
