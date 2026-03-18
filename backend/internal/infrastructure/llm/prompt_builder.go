package llm

import "fmt"

// BuildSystemPrompt creates a dynamic system prompt based on generation parameters.
// All instructions are placed here; the user message contains only the source text.
// This function is package-level so it can be reused by all LLM strategies.
func BuildSystemPrompt(params GenerationParams) string {
	difficulty := params.Difficulty
	if difficulty == "" {
		difficulty = "medium"
	}

	singleChoiceCount := params.NumQuestions - params.MultipleChoiceCount

	// Build task description with exact question type counts
	var taskDescription string
	if params.MultipleChoiceCount == 0 {
		taskDescription = fmt.Sprintf(
			"Создай %d тестовых вопросов типа single_choice на основе предоставленного текста.",
			params.NumQuestions,
		)
	} else {
		taskDescription = fmt.Sprintf(
			"Создай %d тестовых вопросов на основе предоставленного текста:\n- %d вопросов типа single_choice\n- %d вопросов типа multiple_choice",
			params.NumQuestions, singleChoiceCount, params.MultipleChoiceCount,
		)
	}

	// Difficulty-specific guidance
	var difficultyGuide string
	switch difficulty {
	case "easy":
		difficultyGuide = `Сложность: easy (лёгкая).
Вопросы проверяют базовые знания и определения. Используй простые и однозначные формулировки. Правильный ответ должен быть очевиден при знании материала, дистракторы — явно неправильными.`
	case "hard":
		difficultyGuide = `Сложность: hard (сложная).
Вопросы проверяют анализ, синтез и критическое мышление. Формулируй вопросы, требующие глубокого понимания. Дистракторы должны быть правдоподобными и близкими к правильному ответу, чтобы отличить их мог только тот, кто хорошо усвоил материал.`
	default:
		difficultyGuide = `Сложность: medium (средняя).
Вопросы проверяют понимание и применение знаний. Часть дистракторов должна быть правдоподобной и требовать обдумывания, но не вводить в заблуждение знающего студента.`
	}

	// Profile-specific guidance block (empty for universal)
	profileBlock := BuildProfileBlock(params.Profile)
	profileSection := ""
	if profileBlock != "" {
		profileSection = "\n\n" + profileBlock
	}

	// Profile-adapted JSON examples
	examplesSection := buildExamplesSection(difficulty, params.MultipleChoiceCount, params.Profile)

	// Profile-adapted formulation examples (ПЛОХО / ХОРОШО)
	formulationExamples := buildFormulationExamples(params.Profile)

	// Include only the type descriptions that are actually used
	questionTypesSection := "ТИПЫ ВОПРОСОВ:\n- single_choice: 4 варианта ответа (ровно 1 правильный, 3 неправильных)"
	if params.MultipleChoiceCount > 0 {
		questionTypesSection += "\n- multiple_choice: 5–6 вариантов ответа (2–3 правильных, 2–3 неправильных)"
	}

	return fmt.Sprintf(`Ты — профессиональный составитель тестовых заданий для образовательных целей. Генерируй вопросы на русском языке строго в формате JSON.

ЗАДАНИЕ:
%s

%s

%s%s

%s

ПРАВИЛА ФОРМУЛИРОВКИ ВОПРОСОВ:
1. Каждый вопрос должен быть САМОДОСТАТОЧНЫМ — понятным без обращения к исходному тексту.
2. НЕ используй фразы: "В примере выше", "Как показано в тексте", "Согласно лекции", "В данном фрагменте".
3. Если нужен конкретный пример — включи его ПОЛНОСТЬЮ прямо в текст вопроса.
4. Проверяй понимание концепций, а не дословное воспроизведение текста.

%s

Верни ТОЛЬКО валидный JSON без дополнительного текста, комментариев и пояснений вне JSON.`,
		taskDescription,
		questionTypesSection,
		difficultyGuide,
		profileSection,
		examplesSection,
		formulationExamples,
	)
}

// BuildProfileBlock returns a profile-specific instruction block describing the discipline domain
// and formulation style. Returns an empty string for ProfileUniversal or unrecognised profiles.
func BuildProfileBlock(profile AcademicProfile) string {
	switch profile {
	case ProfileTechnological:
		return `ПРОФИЛЬ ДИСЦИПЛИНЫ: технологический (программирование, инженерия, математика, информатика).
Акцентируй вопросы на точных определениях, алгоритмах, формулах и разнице между схожими понятиями.
Формулировки должны быть однозначными. Включай фрагменты кода или формулы прямо в текст вопроса, если они необходимы для понимания.`

	case ProfileNaturalScience:
		return `ПРОФИЛЬ ДИСЦИПЛИНЫ: естественно-научный (физика, химия, биология, экология, геология).
Акцентируй вопросы на законах природы, химических процессах, биологических механизмах и причинно-следственных связях.
Используй строгую научную терминологию. Формулируй вопросы на объяснение явлений и закономерностей.`

	case ProfileHumanities:
		return `ПРОФИЛЬ ДИСЦИПЛИНЫ: гуманитарный (история, литература, философия, лингвистика, культурология).
Акцентируй вопросы на авторстве, периодизации, интерпретации текстов и событий, терминологии научных школ.
Дистракторы могут быть правдоподобными альтернативными интерпретациями. Вопросы должны проверять понимание, а не механическое запоминание.`

	case ProfileSocialEconomic:
		return `ПРОФИЛЬ ДИСЦИПЛИНЫ: социально-экономический (экономика, право, социология, менеджмент, политология).
Акцентируй вопросы на применении концепций к ситуациям, разграничении схожих терминов, содержании нормативной базы.
Можно включать короткие практические сценарии. Проверяй понимание принципов и умение их применять.`

	case ProfileCreative:
		return `ПРОФИЛЬ ДИСЦИПЛИНЫ: творческий (изобразительное искусство, дизайн, музыка, архитектура, театр, кино).
Акцентируй вопросы на художественных стилях и направлениях, исторических периодах, творческих техниках, авторах и произведениях.
Формулируй вопросы на идентификацию стилей и эпох. Дистракторы должны быть из близких по времени или стилю направлений.`

	default:
		// ProfileUniversal and unknown profiles — no additional block
		return ""
	}
}

// buildFormulationExamples returns profile-specific ПЛОХО/ХОРОШО examples
// that guide the LLM on how to formulate self-contained questions for the given discipline domain.
func buildFormulationExamples(profile AcademicProfile) string {
	switch profile {
	case ProfileNaturalScience:
		return `ПРИМЕРЫ ПРАВИЛЬНОЙ И НЕПРАВИЛЬНОЙ ФОРМУЛИРОВКИ:

ПЛОХО: "Согласно тексту, что происходит с газом при нагревании?"
ХОРОШО: "Как изменяется давление идеального газа при увеличении его температуры при постоянном объёме?"

ПЛОХО: "В рассмотренном примере, какие вещества вступили в реакцию?"
ХОРОШО: "При смешивании растворов соляной кислоты (HCl) и гидроксида натрия (NaOH) образуются:"

ПЛОХО: "Что означает упомянутый в лекции закон?"
ХОРОШО: "Что утверждает второй закон термодинамики применительно к замкнутым системам?"

ПЛОХО: "Из приведённого фрагмента, чем отличается митоз от мейоза?"
ХОРОШО: "Какое принципиальное отличие мейоза от митоза с точки зрения числа дочерних клеток и их хромосомного набора?"`

	case ProfileHumanities:
		return `ПРИМЕРЫ ПРАВИЛЬНОЙ И НЕПРАВИЛЬНОЙ ФОРМУЛИРОВКИ:

ПЛОХО: "Как говорится в тексте, кто является автором упомянутого произведения?"
ХОРОШО: "Кто является автором романа «Преступление и наказание», впервые опубликованного в 1866 году?"

ПЛОХО: "Согласно лекции, что такое феодальная раздробленность?"
ХОРОШО: "Что в исторической науке понимается под термином «феодальная раздробленность» применительно к Руси XII–XIII веков?"

ПЛОХО: "В данном фрагменте упоминается философская школа — чем она отличается?"
ХОРОШО: "В чём заключается главное отличие эмпиризма от рационализма как философских методов познания?"

ПЛОХО: "Согласно прочитанному, к какому периоду относится данное направление?"
ХОРОШО: "К какому историческому периоду относится возникновение европейского символизма как литературного и художественного направления?"`

	case ProfileSocialEconomic:
		return `ПРИМЕРЫ ПРАВИЛЬНОЙ И НЕПРАВИЛЬНОЙ ФОРМУЛИРОВКИ:

ПЛОХО: "Согласно тексту, что такое ВВП?"
ХОРОШО: "Что измеряет показатель валового внутреннего продукта (ВВП) страны за определённый период?"

ПЛОХО: "Как сказано в лекции, какой орган принимает федеральные законы?"
ХОРОШО: "Какой орган государственной власти Российской Федерации наделён исключительным правом принятия федеральных законов?"

ПЛОХО: "В приведённом примере, что нарушил работодатель?"
ХОРОШО: "Работодатель уволил сотрудника в день его выхода из ежегодного оплачиваемого отпуска без предварительного уведомления. Какую норму трудового законодательства он нарушил?"

ПЛОХО: "Из текста следует, чем монополия отличается от олигополии?"
ХОРОШО: "Чем структура рынка олигополии принципиально отличается от монополии с точки зрения числа участников и механизма ценообразования?"`

	case ProfileCreative:
		return `ПРИМЕРЫ ПРАВИЛЬНОЙ И НЕПРАВИЛЬНОЙ ФОРМУЛИРОВКИ:

ПЛОХО: "Согласно тексту, к какому стилю относится данное произведение?"
ХОРОШО: "Какой художественный стиль характеризуется использованием геометрических форм, монохромной палитры и полным отказом от изобразительности?"

ПЛОХО: "В лекции упомянута техника — как она называется?"
ХОРОШО: "Как называется техника многоголосного музыкального произведения, в котором несколько независимых мелодических линий развиваются по единым контрапунктическим правилам?"

ПЛОХО: "Из приведённого фрагмента, в какой период создавались эти произведения?"
ХОРОШО: "В какой исторический период сложился стиль «модерн» (Art Nouveau) в европейском искусстве и архитектуре?"

ПЛОХО: "Согласно описанию, чем отличается этот режиссёрский метод?"
ХОРОШО: "В чём состоит ключевое отличие системы Станиславского от метода Брехта в подходе актёра к роли?"`

	default:
		// ProfileTechnological and ProfileUniversal
		return `ПРИМЕРЫ ПРАВИЛЬНОЙ И НЕПРАВИЛЬНОЙ ФОРМУЛИРОВКИ:

ПЛОХО: "Согласно лекции, что такое полиморфизм?"
ХОРОШО: "Что такое полиморфизм в объектно-ориентированном программировании?"

ПЛОХО: "В приведённом выше примере, какой метод будет вызван?"
ХОРОШО: "Дан код:
class Animal { void speak() { System.out.println("..."); } }
class Dog extends Animal { void speak() { System.out.println("Woof"); } }
Animal a = new Dog();
Что выведет a.speak()?"

ПЛОХО: "Из текста следует, какова сложность алгоритма?"
ХОРОШО: "Какова временная сложность алгоритма бинарного поиска в отсортированном массиве из N элементов?"

ПЛОХО: "Как сказано в лекции, что делает эта функция?"
ХОРОШО: "Что вернёт функция:
def f(n): return n if n <= 1 else f(n-1) + f(n-2)
при вызове f(6)?"`
	}
}

// buildExamplesSection returns a profile-adapted JSON examples block.
// Examples use vocabulary appropriate to the discipline profile so the LLM
// understands the expected output style for that domain.
func buildExamplesSection(difficulty string, multipleChoiceCount int, profile AcademicProfile) string {
	singleExample := profileSingleChoiceExample(difficulty, profile)
	if multipleChoiceCount == 0 {
		return "ФОРМАТ ОТВЕТА (строго JSON):\n{\n  \"questions\": [\n" + singleExample + "\n  ]\n}"
	}
	multipleExample := profileMultipleChoiceExample(difficulty, profile)
	return "ФОРМАТ ОТВЕТА (строго JSON):\n{\n  \"questions\": [\n" + singleExample + ",\n" + multipleExample + "\n  ]\n}"
}

// profileSingleChoiceExample returns a domain-appropriate single_choice question example.
func profileSingleChoiceExample(difficulty string, profile AcademicProfile) string {
	switch profile {
	case ProfileNaturalScience:
		return `    {
      "question": "Какой закон описывает обратно пропорциональную зависимость между давлением и объёмом газа при постоянной температуре?",
      "type": "single_choice",
      "difficulty": "` + difficulty + `",
      "answers": [
        {"text": "Закон Бойля–Мариотта", "is_correct": true},
        {"text": "Закон Гей-Люссака", "is_correct": false},
        {"text": "Закон Шарля", "is_correct": false},
        {"text": "Закон Авогадро", "is_correct": false}
      ],
      "explanation": "Закон Бойля–Мариотта: при постоянной температуре давление газа обратно пропорционально его объёму (pV = const)."
    }`

	case ProfileHumanities:
		return `    {
      "question": "К какому литературному направлению относится творчество Фёдора Достоевского?",
      "type": "single_choice",
      "difficulty": "` + difficulty + `",
      "answers": [
        {"text": "Реализм", "is_correct": true},
        {"text": "Романтизм", "is_correct": false},
        {"text": "Символизм", "is_correct": false},
        {"text": "Классицизм", "is_correct": false}
      ],
      "explanation": "Достоевский — представитель критического реализма XIX века, его произведения исследуют психологию человека в социальных условиях."
    }`

	case ProfileSocialEconomic:
		return `    {
      "question": "Что является главной характеристикой рынка совершенной конкуренции?",
      "type": "single_choice",
      "difficulty": "` + difficulty + `",
      "answers": [
        {"text": "Большое число продавцов и покупателей, ни один из которых не влияет на цену", "is_correct": true},
        {"text": "Один продавец контролирует весь рынок", "is_correct": false},
        {"text": "Несколько крупных продавцов устанавливают цену совместно", "is_correct": false},
        {"text": "Товары разных продавцов существенно дифференцированы", "is_correct": false}
      ],
      "explanation": "На рынке совершенной конкуренции множество участников торгует однородным товаром, и каждый принимает цену как данную (price taker)."
    }`

	case ProfileCreative:
		return `    {
      "question": "Для какого художественного направления характерна передача мимолётных впечатлений от игры света и цвета, а не точного воспроизведения действительности?",
      "type": "single_choice",
      "difficulty": "` + difficulty + `",
      "answers": [
        {"text": "Импрессионизм", "is_correct": true},
        {"text": "Реализм", "is_correct": false},
        {"text": "Сюрреализм", "is_correct": false},
        {"text": "Экспрессионизм", "is_correct": false}
      ],
      "explanation": "Импрессионизм (от фр. impression — впечатление) возник в 1870-е гг. во Франции. Художники стремились запечатлеть мимолётное состояние природы через игру света и цвета."
    }`

	default:
		// ProfileTechnological and ProfileUniversal
		return `    {
      "question": "Что такое инкапсуляция в объектно-ориентированном программировании?",
      "type": "single_choice",
      "difficulty": "` + difficulty + `",
      "answers": [
        {"text": "Сокрытие внутреннего состояния объекта и предоставление доступа только через публичный интерфейс", "is_correct": true},
        {"text": "Возможность объекта принимать разные формы", "is_correct": false},
        {"text": "Механизм наследования свойств от родительского класса", "is_correct": false},
        {"text": "Способность класса порождать другие классы", "is_correct": false}
      ],
      "explanation": "Инкапсуляция — принцип ООП, при котором данные объекта скрыты от внешнего доступа и доступны только через методы."
    }`
	}
}

// profileMultipleChoiceExample returns a domain-appropriate multiple_choice question example.
func profileMultipleChoiceExample(difficulty string, profile AcademicProfile) string {
	switch profile {
	case ProfileNaturalScience:
		return `    {
      "question": "Какие из перечисленных процессов являются экзотермическими реакциями?",
      "type": "multiple_choice",
      "difficulty": "` + difficulty + `",
      "answers": [
        {"text": "Горение природного газа", "is_correct": true},
        {"text": "Разложение воды электролизом", "is_correct": false},
        {"text": "Нейтрализация кислоты щёлочью", "is_correct": true},
        {"text": "Растворение нитрата аммония в воде", "is_correct": false},
        {"text": "Окисление глюкозы в клетке", "is_correct": true}
      ],
      "explanation": "Экзотермические реакции выделяют тепло. Горение, нейтрализация и клеточное дыхание сопровождаются выделением энергии. Электролиз и растворение NH4NO3 — эндотермические процессы."
    }`

	case ProfileHumanities:
		return `    {
      "question": "Какие из перечисленных черт характерны для эпохи Просвещения (XVIII век)?",
      "type": "multiple_choice",
      "difficulty": "` + difficulty + `",
      "answers": [
        {"text": "Вера в силу разума и науки как инструментов прогресса", "is_correct": true},
        {"text": "Критика религиозного догматизма и церковной власти", "is_correct": true},
        {"text": "Идеализация средневекового уклада и рыцарства", "is_correct": false},
        {"text": "Концепция естественных прав человека", "is_correct": true},
        {"text": "Отрицание роли государства в жизни общества", "is_correct": false}
      ],
      "explanation": "Просвещение провозглашало приоритет разума, критиковало церковь и феодальный строй, разрабатывало теории естественного права. Идеализация средневековья — черта романтизма."
    }`

	case ProfileSocialEconomic:
		return `    {
      "question": "Какие из перечисленных источников относятся к нормативно-правовым актам в российской правовой системе?",
      "type": "multiple_choice",
      "difficulty": "` + difficulty + `",
      "answers": [
        {"text": "Федеральный конституционный закон", "is_correct": true},
        {"text": "Постановление Правительства РФ", "is_correct": true},
        {"text": "Решение суда по конкретному делу", "is_correct": false},
        {"text": "Указ Президента РФ", "is_correct": true},
        {"text": "Правовой обычай", "is_correct": false}
      ],
      "explanation": "Нормативно-правовые акты — официальные документы с общеобязательными правилами. Судебные решения и правовые обычаи — иные источники права, не относящиеся к НПА."
    }`

	case ProfileCreative:
		return `    {
      "question": "Какие из перечисленных особенностей характерны для архитектуры барокко?",
      "type": "multiple_choice",
      "difficulty": "` + difficulty + `",
      "answers": [
        {"text": "Пышный декор и обилие орнаментов", "is_correct": true},
        {"text": "Динамичные криволинейные формы", "is_correct": true},
        {"text": "Строгая симметрия и минимализм деталей", "is_correct": false},
        {"text": "Использование тромплёя и иллюзионистической живописи", "is_correct": true},
        {"text": "Прямые линии и функциональность как главный принцип", "is_correct": false}
      ],
      "explanation": "Барокко (XVII–XVIII вв.) отличается пышностью, динамизмом и иллюзионизмом. Строгость и минимализм характерны для классицизма, функциональность — для модернизма."
    }`

	default:
		// ProfileTechnological and ProfileUniversal
		return `    {
      "question": "Какие из следующих утверждений верны для интерфейсов в языке Go?",
      "type": "multiple_choice",
      "difficulty": "` + difficulty + `",
      "answers": [
        {"text": "Интерфейс реализуется неявно — достаточно реализовать все его методы", "is_correct": true},
        {"text": "Интерфейс может содержать поля данных", "is_correct": false},
        {"text": "Пустой интерфейс interface{} удовлетворяет любой тип", "is_correct": true},
        {"text": "Для реализации интерфейса нужно явно указать ключевое слово implements", "is_correct": false},
        {"text": "Интерфейс может встраивать другие интерфейсы", "is_correct": true}
      ],
      "explanation": "В Go интерфейсы реализуются неявно. Пустой интерфейс принимает любой тип. Интерфейсы могут встраивать другие интерфейсы. Поля данных и implements отсутствуют."
    }`
	}
}
