-- Migration 000005: Create academic_profiles reference table
-- Promotes AcademicProfile from a Go constant to a DB-persisted reference table
-- with per-profile temperature, display labels, and LLM instruction block.

CREATE TABLE academic_profiles (
    id                  UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    code                VARCHAR(50)  UNIQUE NOT NULL
                          CHECK (code IN (
                              'technological', 'natural_science', 'humanities',
                              'social_economic', 'creative', 'universal'
                          )),
    display_name        VARCHAR(100) NOT NULL,
    description         TEXT,
    -- Profile instruction block sent to the LLM in the system prompt.
    -- NULL for 'universal' (no additional block needed).
    profile_instruction TEXT,
    -- LLM temperature: lower for exact/technical domains, higher for interpretive ones.
    temperature         NUMERIC(3,2) NOT NULL DEFAULT 0.50
                          CHECK (temperature BETWEEN 0.00 AND 1.00),
    sort_order          SMALLINT     NOT NULL DEFAULT 0,
    created_at          TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE  academic_profiles                 IS 'Reference table for educational discipline profiles used to customise LLM question generation';
COMMENT ON COLUMN academic_profiles.code            IS 'Machine-readable key matching llm.AcademicProfile Go constants';
COMMENT ON COLUMN academic_profiles.profile_instruction IS 'Full profile block injected into the LLM system prompt; NULL for universal profile';
COMMENT ON COLUMN academic_profiles.temperature     IS 'LLM sampling temperature: 0.20 (technical) → 0.70 (creative)';

-- Seed all 6 profiles
INSERT INTO academic_profiles (code, display_name, description, profile_instruction, temperature, sort_order)
VALUES
(
    'technological',
    'Технологический',
    'Программирование, инженерия, математика, информатика',
    'ПРОФИЛЬ ДИСЦИПЛИНЫ: технологический (программирование, инженерия, математика, информатика).
Акцентируй вопросы на точных определениях, алгоритмах, формулах и разнице между схожими понятиями.
Формулировки должны быть однозначными. Включай фрагменты кода или формулы прямо в текст вопроса, если они необходимы для понимания.',
    0.20,
    1
),
(
    'natural_science',
    'Естественно-научный',
    'Физика, химия, биология, экология, геология',
    'ПРОФИЛЬ ДИСЦИПЛИНЫ: естественно-научный (физика, химия, биология, экология, геология).
Акцентируй вопросы на законах природы, химических процессах, биологических механизмах и причинно-следственных связях.
Используй строгую научную терминологию. Формулируй вопросы на объяснение явлений и закономерностей.',
    0.30,
    2
),
(
    'humanities',
    'Гуманитарный',
    'История, литература, философия, лингвистика, культурология',
    'ПРОФИЛЬ ДИСЦИПЛИНЫ: гуманитарный (история, литература, философия, лингвистика, культурология).
Акцентируй вопросы на авторстве, периодизации, интерпретации текстов и событий, терминологии научных школ.
Дистракторы могут быть правдоподобными альтернативными интерпретациями. Вопросы должны проверять понимание, а не механическое запоминание.',
    0.65,
    3
),
(
    'social_economic',
    'Социально-экономический',
    'Экономика, право, социология, менеджмент, политология',
    'ПРОФИЛЬ ДИСЦИПЛИНЫ: социально-экономический (экономика, право, социология, менеджмент, политология).
Акцентируй вопросы на применении концепций к ситуациям, разграничении схожих терминов, содержании нормативной базы.
Можно включать короткие практические сценарии. Проверяй понимание принципов и умение их применять.',
    0.45,
    4
),
(
    'creative',
    'Творческий',
    'Изобразительное искусство, дизайн, музыка, архитектура, театр, кино',
    'ПРОФИЛЬ ДИСЦИПЛИНЫ: творческий (изобразительное искусство, дизайн, музыка, архитектура, театр, кино).
Акцентируй вопросы на художественных стилях и направлениях, исторических периодах, творческих техниках, авторах и произведениях.
Формулируй вопросы на идентификацию стилей и эпох. Дистракторы должны быть из близких по времени или стилю направлений.',
    0.70,
    5
),
(
    'universal',
    'Универсальный',
    'Общий профиль без специализации',
    NULL,
    0.50,
    6
);
