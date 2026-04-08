-- Migration 000007: Test generation configs and user profile preference
-- test_generation_configs: immutable record of LLM params used per generated test (1:0..1 with tests)
-- users.preferred_academic_profile_id: optional teacher preference for default profile

-- ============================================================
-- 1. test_generation_configs
-- ============================================================
CREATE TABLE test_generation_configs (
    -- test_id is both PK and FK: enforces 1:0..1 relationship with tests
    test_id             UUID        PRIMARY KEY
                          REFERENCES tests(id) ON DELETE CASCADE,
    academic_profile_id UUID        NOT NULL
                          REFERENCES academic_profiles(id) ON DELETE RESTRICT,
    difficulty          VARCHAR(50) NOT NULL
                          CHECK (difficulty IN ('easy', 'medium', 'hard')),
    num_questions       SMALLINT    NOT NULL
                          CHECK (num_questions BETWEEN 1 AND 50),
    llm_provider        VARCHAR(50) NOT NULL
                          CHECK (llm_provider IN ('perplexity', 'openai', 'yandexgpt')),
    language            VARCHAR(10) NOT NULL DEFAULT 'ru',
    -- No updated_at: this record is immutable (historical generation fact)
    created_at          TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE  test_generation_configs                     IS 'Immutable record of LLM generation parameters for AI-generated tests; absent for manually created tests';
COMMENT ON COLUMN test_generation_configs.test_id             IS 'PK = FK: one-to-zero-or-one with tests';
COMMENT ON COLUMN test_generation_configs.academic_profile_id IS 'Profile whose instruction block was sent to the LLM';
COMMENT ON COLUMN test_generation_configs.language            IS 'BCP-47 language tag used for generation (e.g. ru, en)';

CREATE INDEX idx_test_gen_configs_profile_id
    ON test_generation_configs(academic_profile_id);

-- ============================================================
-- 2. users.preferred_academic_profile_id
-- ============================================================
ALTER TABLE users
    ADD COLUMN preferred_academic_profile_id UUID
        REFERENCES academic_profiles(id) ON DELETE SET NULL;

COMMENT ON COLUMN users.preferred_academic_profile_id IS 'Teacher default academic profile; NULL means no preference set (application defaults to universal)';

-- Partial index: only index rows where preference is set
CREATE INDEX idx_users_preferred_profile_id
    ON users(preferred_academic_profile_id)
    WHERE preferred_academic_profile_id IS NOT NULL;
