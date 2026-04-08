-- Migration 000006: Create profile example tables and seed all content
-- formulation_examples: ПЛОХО/ХОРОШО pairs guiding the LLM on self-contained question phrasing
-- question_examples:    Domain-appropriate JSON example questions for the LLM format section

-- ============================================================
-- Table: academic_profile_formulation_examples
-- ============================================================
CREATE TABLE academic_profile_formulation_examples (
    id                  UUID     PRIMARY KEY DEFAULT uuid_generate_v4(),
    academic_profile_id UUID     NOT NULL
                          REFERENCES academic_profiles(id) ON DELETE CASCADE,
    bad_example         TEXT     NOT NULL,
    good_example        TEXT     NOT NULL,
    order_num           SMALLINT NOT NULL DEFAULT 0
);

COMMENT ON TABLE  academic_profile_formulation_examples             IS 'ПЛОХО/ХОРОШО phrasing examples per profile injected into the LLM system prompt';
COMMENT ON COLUMN academic_profile_formulation_examples.bad_example IS 'Example of a context-dependent (bad) question formulation';
COMMENT ON COLUMN academic_profile_formulation_examples.good_example IS 'Self-contained (good) reformulation of the same question';

CREATE UNIQUE INDEX uq_profile_formulation_order
    ON academic_profile_formulation_examples(academic_profile_id, order_num);
CREATE INDEX idx_profile_formulation_profile_id
    ON academic_profile_formulation_examples(academic_profile_id);

-- ============================================================
-- Table: academic_profile_question_examples
-- ============================================================
CREATE TABLE academic_profile_question_examples (
    id                  UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    academic_profile_id UUID        NOT NULL
                          REFERENCES academic_profiles(id) ON DELETE CASCADE,
    question_type       VARCHAR(50) NOT NULL
                          CHECK (question_type IN ('single_choice', 'multiple_choice')),
    question_text       TEXT        NOT NULL,
    -- JSON array: [{text: "...", is_correct: true/false}, ...]
    -- 4 items for single_choice (1 correct), 5-6 for multiple_choice (2-3 correct)
    answers             JSONB       NOT NULL,
    explanation         TEXT        NOT NULL,
    order_num           SMALLINT    NOT NULL DEFAULT 0
);

COMMENT ON TABLE  academic_profile_question_examples              IS 'Domain-specific JSON question examples shown to the LLM to demonstrate the expected output format and vocabulary';
COMMENT ON COLUMN academic_profile_question_examples.answers      IS 'JSON array of {text, is_correct} objects; difficulty field is injected at runtime by the prompt builder';
COMMENT ON COLUMN academic_profile_question_examples.question_type IS 'single_choice or multiple_choice — determines answer count and correct-answer cardinality';

CREATE UNIQUE INDEX uq_profile_qexample_type_order
    ON academic_profile_question_examples(academic_profile_id, question_type, order_num);
CREATE INDEX idx_profile_qexample_profile_id
    ON academic_profile_question_examples(academic_profile_id);

-- ============================================================
-- Seed: formulation examples
-- ============================================================

-- technological
INSERT INTO academic_profile_formulation_examples (academic_profile_id, bad_example, good_example, order_num)
SELECT id,
       'Согласно лекции, что такое полиморфизм?',
       'Что такое полиморфизм в объектно-ориентированном программировании?',
       1
FROM academic_profiles WHERE code = 'technological';

INSERT INTO academic_profile_formulation_examples (academic_profile_id, bad_example, good_example, order_num)
SELECT id,
       'В приведённом выше примере, какой метод будет вызван?',
       'Дан код:
class Animal { void speak() { System.out.println("..."); } }
class Dog extends Animal { void speak() { System.out.println("Woof"); } }
Animal a = new Dog();
Что выведет a.speak()?',
       2
FROM academic_profiles WHERE code = 'technological';

INSERT INTO academic_profile_formulation_examples (academic_profile_id, bad_example, good_example, order_num)
SELECT id,
       'Из текста следует, какова сложность алгоритма?',
       'Какова временная сложность алгоритма бинарного поиска в отсортированном массиве из N элементов?',
       3
FROM academic_profiles WHERE code = 'technological';

INSERT INTO academic_profile_formulation_examples (academic_profile_id, bad_example, good_example, order_num)
SELECT id,
       'Как сказано в лекции, что делает эта функция?',
       'Что вернёт функция:
def f(n): return n if n <= 1 else f(n-1) + f(n-2)
при вызове f(6)?',
       4
FROM academic_profiles WHERE code = 'technological';

-- natural_science
INSERT INTO academic_profile_formulation_examples (academic_profile_id, bad_example, good_example, order_num)
SELECT id,
       'Согласно тексту, что происходит с газом при нагревании?',
       'Как изменяется давление идеального газа при увеличении его температуры при постоянном объёме?',
       1
FROM academic_profiles WHERE code = 'natural_science';

INSERT INTO academic_profile_formulation_examples (academic_profile_id, bad_example, good_example, order_num)
SELECT id,
       'В рассмотренном примере, какие вещества вступили в реакцию?',
       'При смешивании растворов соляной кислоты (HCl) и гидроксида натрия (NaOH) образуются:',
       2
FROM academic_profiles WHERE code = 'natural_science';

INSERT INTO academic_profile_formulation_examples (academic_profile_id, bad_example, good_example, order_num)
SELECT id,
       'Что означает упомянутый в лекции закон?',
       'Что утверждает второй закон термодинамики применительно к замкнутым системам?',
       3
FROM academic_profiles WHERE code = 'natural_science';

INSERT INTO academic_profile_formulation_examples (academic_profile_id, bad_example, good_example, order_num)
SELECT id,
       'Из приведённого фрагмента, чем отличается митоз от мейоза?',
       'Какое принципиальное отличие мейоза от митоза с точки зрения числа дочерних клеток и их хромосомного набора?',
       4
FROM academic_profiles WHERE code = 'natural_science';

-- humanities
INSERT INTO academic_profile_formulation_examples (academic_profile_id, bad_example, good_example, order_num)
SELECT id,
       'Как говорится в тексте, кто является автором упомянутого произведения?',
       'Кто является автором романа «Преступление и наказание», впервые опубликованного в 1866 году?',
       1
FROM academic_profiles WHERE code = 'humanities';

INSERT INTO academic_profile_formulation_examples (academic_profile_id, bad_example, good_example, order_num)
SELECT id,
       'Согласно лекции, что такое феодальная раздробленность?',
       'Что в исторической науке понимается под термином «феодальная раздробленность» применительно к Руси XII–XIII веков?',
       2
FROM academic_profiles WHERE code = 'humanities';

INSERT INTO academic_profile_formulation_examples (academic_profile_id, bad_example, good_example, order_num)
SELECT id,
       'В данном фрагменте упоминается философская школа — чем она отличается?',
       'В чём заключается главное отличие эмпиризма от рационализма как философских методов познания?',
       3
FROM academic_profiles WHERE code = 'humanities';

INSERT INTO academic_profile_formulation_examples (academic_profile_id, bad_example, good_example, order_num)
SELECT id,
       'Согласно прочитанному, к какому периоду относится данное направление?',
       'К какому историческому периоду относится возникновение европейского символизма как литературного и художественного направления?',
       4
FROM academic_profiles WHERE code = 'humanities';

-- social_economic
INSERT INTO academic_profile_formulation_examples (academic_profile_id, bad_example, good_example, order_num)
SELECT id,
       'Согласно тексту, что такое ВВП?',
       'Что измеряет показатель валового внутреннего продукта (ВВП) страны за определённый период?',
       1
FROM academic_profiles WHERE code = 'social_economic';

INSERT INTO academic_profile_formulation_examples (academic_profile_id, bad_example, good_example, order_num)
SELECT id,
       'Как сказано в лекции, какой орган принимает федеральные законы?',
       'Какой орган государственной власти Российской Федерации наделён исключительным правом принятия федеральных законов?',
       2
FROM academic_profiles WHERE code = 'social_economic';

INSERT INTO academic_profile_formulation_examples (academic_profile_id, bad_example, good_example, order_num)
SELECT id,
       'В приведённом примере, что нарушил работодатель?',
       'Работодатель уволил сотрудника в день его выхода из ежегодного оплачиваемого отпуска без предварительного уведомления. Какую норму трудового законодательства он нарушил?',
       3
FROM academic_profiles WHERE code = 'social_economic';

INSERT INTO academic_profile_formulation_examples (academic_profile_id, bad_example, good_example, order_num)
SELECT id,
       'Из текста следует, чем монополия отличается от олигополии?',
       'Чем структура рынка олигополии принципиально отличается от монополии с точки зрения числа участников и механизма ценообразования?',
       4
FROM academic_profiles WHERE code = 'social_economic';

-- creative
INSERT INTO academic_profile_formulation_examples (academic_profile_id, bad_example, good_example, order_num)
SELECT id,
       'Согласно тексту, к какому стилю относится данное произведение?',
       'Какой художественный стиль характеризуется использованием геометрических форм, монохромной палитры и полным отказом от изобразительности?',
       1
FROM academic_profiles WHERE code = 'creative';

INSERT INTO academic_profile_formulation_examples (academic_profile_id, bad_example, good_example, order_num)
SELECT id,
       'В лекции упомянута техника — как она называется?',
       'Как называется техника многоголосного музыкального произведения, в котором несколько независимых мелодических линий развиваются по единым контрапунктическим правилам?',
       2
FROM academic_profiles WHERE code = 'creative';

INSERT INTO academic_profile_formulation_examples (academic_profile_id, bad_example, good_example, order_num)
SELECT id,
       'Из приведённого фрагмента, в какой период создавались эти произведения?',
       'В какой исторический период сложился стиль «модерн» (Art Nouveau) в европейском искусстве и архитектуре?',
       3
FROM academic_profiles WHERE code = 'creative';

INSERT INTO academic_profile_formulation_examples (academic_profile_id, bad_example, good_example, order_num)
SELECT id,
       'Согласно описанию, чем отличается этот режиссёрский метод?',
       'В чём состоит ключевое отличие системы Станиславского от метода Брехта в подходе актёра к роли?',
       4
FROM academic_profiles WHERE code = 'creative';

-- universal (inherits technological examples)
INSERT INTO academic_profile_formulation_examples (academic_profile_id, bad_example, good_example, order_num)
SELECT id,
       'Согласно лекции, что такое полиморфизм?',
       'Что такое полиморфизм в объектно-ориентированном программировании?',
       1
FROM academic_profiles WHERE code = 'universal';

INSERT INTO academic_profile_formulation_examples (academic_profile_id, bad_example, good_example, order_num)
SELECT id,
       'В приведённом выше примере, какой метод будет вызван?',
       'Дан код:
class Animal { void speak() { System.out.println("..."); } }
class Dog extends Animal { void speak() { System.out.println("Woof"); } }
Animal a = new Dog();
Что выведет a.speak()?',
       2
FROM academic_profiles WHERE code = 'universal';

INSERT INTO academic_profile_formulation_examples (academic_profile_id, bad_example, good_example, order_num)
SELECT id,
       'Из текста следует, какова сложность алгоритма?',
       'Какова временная сложность алгоритма бинарного поиска в отсортированном массиве из N элементов?',
       3
FROM academic_profiles WHERE code = 'universal';

INSERT INTO academic_profile_formulation_examples (academic_profile_id, bad_example, good_example, order_num)
SELECT id,
       'Как сказано в лекции, что делает эта функция?',
       'Что вернёт функция:
def f(n): return n if n <= 1 else f(n-1) + f(n-2)
при вызове f(6)?',
       4
FROM academic_profiles WHERE code = 'universal';

-- ============================================================
-- Seed: question examples
-- ============================================================

-- technological — single_choice
INSERT INTO academic_profile_question_examples (academic_profile_id, question_type, question_text, answers, explanation, order_num)
SELECT id,
       'single_choice',
       'Что такое инкапсуляция в объектно-ориентированном программировании?',
       '[
           {"text": "Сокрытие внутреннего состояния объекта и предоставление доступа только через публичный интерфейс", "is_correct": true},
           {"text": "Возможность объекта принимать разные формы", "is_correct": false},
           {"text": "Механизм наследования свойств от родительского класса", "is_correct": false},
           {"text": "Способность класса порождать другие классы", "is_correct": false}
       ]'::jsonb,
       'Инкапсуляция — принцип ООП, при котором данные объекта скрыты от внешнего доступа и доступны только через методы.',
       1
FROM academic_profiles WHERE code = 'technological';

-- technological — multiple_choice
INSERT INTO academic_profile_question_examples (academic_profile_id, question_type, question_text, answers, explanation, order_num)
SELECT id,
       'multiple_choice',
       'Какие из следующих утверждений верны для интерфейсов в языке Go?',
       '[
           {"text": "Интерфейс реализуется неявно — достаточно реализовать все его методы", "is_correct": true},
           {"text": "Интерфейс может содержать поля данных", "is_correct": false},
           {"text": "Пустой интерфейс interface{} удовлетворяет любой тип", "is_correct": true},
           {"text": "Для реализации интерфейса нужно явно указать ключевое слово implements", "is_correct": false},
           {"text": "Интерфейс может встраивать другие интерфейсы", "is_correct": true}
       ]'::jsonb,
       'В Go интерфейсы реализуются неявно. Пустой интерфейс принимает любой тип. Интерфейсы могут встраивать другие интерфейсы. Поля данных и implements отсутствуют.',
       1
FROM academic_profiles WHERE code = 'technological';

-- natural_science — single_choice
INSERT INTO academic_profile_question_examples (academic_profile_id, question_type, question_text, answers, explanation, order_num)
SELECT id,
       'single_choice',
       'Какой закон описывает обратно пропорциональную зависимость между давлением и объёмом газа при постоянной температуре?',
       '[
           {"text": "Закон Бойля–Мариотта", "is_correct": true},
           {"text": "Закон Гей-Люссака", "is_correct": false},
           {"text": "Закон Шарля", "is_correct": false},
           {"text": "Закон Авогадро", "is_correct": false}
       ]'::jsonb,
       'Закон Бойля–Мариотта: при постоянной температуре давление газа обратно пропорционально его объёму (pV = const).',
       1
FROM academic_profiles WHERE code = 'natural_science';

-- natural_science — multiple_choice
INSERT INTO academic_profile_question_examples (academic_profile_id, question_type, question_text, answers, explanation, order_num)
SELECT id,
       'multiple_choice',
       'Какие из перечисленных процессов являются экзотермическими реакциями?',
       '[
           {"text": "Горение природного газа", "is_correct": true},
           {"text": "Разложение воды электролизом", "is_correct": false},
           {"text": "Нейтрализация кислоты щёлочью", "is_correct": true},
           {"text": "Растворение нитрата аммония в воде", "is_correct": false},
           {"text": "Окисление глюкозы в клетке", "is_correct": true}
       ]'::jsonb,
       'Экзотермические реакции выделяют тепло. Горение, нейтрализация и клеточное дыхание сопровождаются выделением энергии. Электролиз и растворение NH4NO3 — эндотермические процессы.',
       1
FROM academic_profiles WHERE code = 'natural_science';

-- humanities — single_choice
INSERT INTO academic_profile_question_examples (academic_profile_id, question_type, question_text, answers, explanation, order_num)
SELECT id,
       'single_choice',
       'К какому литературному направлению относится творчество Фёдора Достоевского?',
       '[
           {"text": "Реализм", "is_correct": true},
           {"text": "Романтизм", "is_correct": false},
           {"text": "Символизм", "is_correct": false},
           {"text": "Классицизм", "is_correct": false}
       ]'::jsonb,
       'Достоевский — представитель критического реализма XIX века, его произведения исследуют психологию человека в социальных условиях.',
       1
FROM academic_profiles WHERE code = 'humanities';

-- humanities — multiple_choice
INSERT INTO academic_profile_question_examples (academic_profile_id, question_type, question_text, answers, explanation, order_num)
SELECT id,
       'multiple_choice',
       'Какие из перечисленных черт характерны для эпохи Просвещения (XVIII век)?',
       '[
           {"text": "Вера в силу разума и науки как инструментов прогресса", "is_correct": true},
           {"text": "Критика религиозного догматизма и церковной власти", "is_correct": true},
           {"text": "Идеализация средневекового уклада и рыцарства", "is_correct": false},
           {"text": "Концепция естественных прав человека", "is_correct": true},
           {"text": "Отрицание роли государства в жизни общества", "is_correct": false}
       ]'::jsonb,
       'Просвещение провозглашало приоритет разума, критиковало церковь и феодальный строй, разрабатывало теории естественного права. Идеализация средневековья — черта романтизма.',
       1
FROM academic_profiles WHERE code = 'humanities';

-- social_economic — single_choice
INSERT INTO academic_profile_question_examples (academic_profile_id, question_type, question_text, answers, explanation, order_num)
SELECT id,
       'single_choice',
       'Что является главной характеристикой рынка совершенной конкуренции?',
       '[
           {"text": "Большое число продавцов и покупателей, ни один из которых не влияет на цену", "is_correct": true},
           {"text": "Один продавец контролирует весь рынок", "is_correct": false},
           {"text": "Несколько крупных продавцов устанавливают цену совместно", "is_correct": false},
           {"text": "Товары разных продавцов существенно дифференцированы", "is_correct": false}
       ]'::jsonb,
       'На рынке совершенной конкуренции множество участников торгует однородным товаром, и каждый принимает цену как данную (price taker).',
       1
FROM academic_profiles WHERE code = 'social_economic';

-- social_economic — multiple_choice
INSERT INTO academic_profile_question_examples (academic_profile_id, question_type, question_text, answers, explanation, order_num)
SELECT id,
       'multiple_choice',
       'Какие из перечисленных источников относятся к нормативно-правовым актам в российской правовой системе?',
       '[
           {"text": "Федеральный конституционный закон", "is_correct": true},
           {"text": "Постановление Правительства РФ", "is_correct": true},
           {"text": "Решение суда по конкретному делу", "is_correct": false},
           {"text": "Указ Президента РФ", "is_correct": true},
           {"text": "Правовой обычай", "is_correct": false}
       ]'::jsonb,
       'Нормативно-правовые акты — официальные документы с общеобязательными правилами. Судебные решения и правовые обычаи — иные источники права, не относящиеся к НПА.',
       1
FROM academic_profiles WHERE code = 'social_economic';

-- creative — single_choice
INSERT INTO academic_profile_question_examples (academic_profile_id, question_type, question_text, answers, explanation, order_num)
SELECT id,
       'single_choice',
       'Для какого художественного направления характерна передача мимолётных впечатлений от игры света и цвета, а не точного воспроизведения действительности?',
       '[
           {"text": "Импрессионизм", "is_correct": true},
           {"text": "Реализм", "is_correct": false},
           {"text": "Сюрреализм", "is_correct": false},
           {"text": "Экспрессионизм", "is_correct": false}
       ]'::jsonb,
       'Импрессионизм (от фр. impression — впечатление) возник в 1870-е гг. во Франции. Художники стремились запечатлеть мимолётное состояние природы через игру света и цвета.',
       1
FROM academic_profiles WHERE code = 'creative';

-- creative — multiple_choice
INSERT INTO academic_profile_question_examples (academic_profile_id, question_type, question_text, answers, explanation, order_num)
SELECT id,
       'multiple_choice',
       'Какие из перечисленных особенностей характерны для архитектуры барокко?',
       '[
           {"text": "Пышный декор и обилие орнаментов", "is_correct": true},
           {"text": "Динамичные криволинейные формы", "is_correct": true},
           {"text": "Строгая симметрия и минимализм деталей", "is_correct": false},
           {"text": "Использование тромплёя и иллюзионистической живописи", "is_correct": true},
           {"text": "Прямые линии и функциональность как главный принцип", "is_correct": false}
       ]'::jsonb,
       'Барокко (XVII–XVIII вв.) отличается пышностью, динамизмом и иллюзионизмом. Строгость и минимализм характерны для классицизма, функциональность — для модернизма.',
       1
FROM academic_profiles WHERE code = 'creative';

-- universal — single_choice (same as technological)
INSERT INTO academic_profile_question_examples (academic_profile_id, question_type, question_text, answers, explanation, order_num)
SELECT id,
       'single_choice',
       'Что такое инкапсуляция в объектно-ориентированном программировании?',
       '[
           {"text": "Сокрытие внутреннего состояния объекта и предоставление доступа только через публичный интерфейс", "is_correct": true},
           {"text": "Возможность объекта принимать разные формы", "is_correct": false},
           {"text": "Механизм наследования свойств от родительского класса", "is_correct": false},
           {"text": "Способность класса порождать другие классы", "is_correct": false}
       ]'::jsonb,
       'Инкапсуляция — принцип ООП, при котором данные объекта скрыты от внешнего доступа и доступны только через методы.',
       1
FROM academic_profiles WHERE code = 'universal';

-- universal — multiple_choice (same as technological)
INSERT INTO academic_profile_question_examples (academic_profile_id, question_type, question_text, answers, explanation, order_num)
SELECT id,
       'multiple_choice',
       'Какие из следующих утверждений верны для интерфейсов в языке Go?',
       '[
           {"text": "Интерфейс реализуется неявно — достаточно реализовать все его методы", "is_correct": true},
           {"text": "Интерфейс может содержать поля данных", "is_correct": false},
           {"text": "Пустой интерфейс interface{} удовлетворяет любой тип", "is_correct": true},
           {"text": "Для реализации интерфейса нужно явно указать ключевое слово implements", "is_correct": false},
           {"text": "Интерфейс может встраивать другие интерфейсы", "is_correct": true}
       ]'::jsonb,
       'В Go интерфейсы реализуются неявно. Пустой интерфейс принимает любой тип. Интерфейсы могут встраивать другие интерфейсы. Поля данных и implements отсутствуют.',
       1
FROM academic_profiles WHERE code = 'universal';
