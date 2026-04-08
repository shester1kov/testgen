-- Rollback migration 000007
DROP INDEX IF EXISTS idx_users_preferred_profile_id;
ALTER TABLE users DROP COLUMN IF EXISTS preferred_academic_profile_id;

DROP INDEX IF EXISTS idx_test_gen_configs_profile_id;
DROP TABLE IF EXISTS test_generation_configs;
