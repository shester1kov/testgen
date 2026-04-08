package llm

import "fmt"

// BuildSystemPrompt creates a dynamic system prompt based on generation parameters.
// All instructions are placed here; the user message contains only the source text.
// Profile-specific content (instruction, formulation examples, question examples) is
// sourced from ProfileData, which is loaded from the database by the generate use case.
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

	// Profile-specific guidance block from DB; empty string when absent or universal
	profileSection := ""
	if params.ProfileData != nil && params.ProfileData.Instruction != "" {
		profileSection = "\n\n" + params.ProfileData.Instruction
	}

	// Profile-adapted JSON examples from DB
	var examplesSection string
	if params.ProfileData != nil {
		examplesSection = buildExamplesSectionFromData(difficulty, params.MultipleChoiceCount, params.ProfileData.QuestionExamples)
	} else {
		examplesSection = "ФОРМАТ ОТВЕТА (строго JSON):\n{\n  \"questions\": []\n}"
	}

	// Profile-adapted formulation examples (ПЛОХО / ХОРОШО) from DB
	var formulationExamples string
	if params.ProfileData != nil {
		formulationExamples = buildFormulationExamplesFromData(params.ProfileData.Formulations)
	}

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

// buildFormulationExamplesFromData assembles the ПЛОХО/ХОРОШО section from DB-loaded data.
func buildFormulationExamplesFromData(formulations []FormulationExample) string {
	if len(formulations) == 0 {
		return ""
	}
	result := "ПРИМЕРЫ ПРАВИЛЬНОЙ И НЕПРАВИЛЬНОЙ ФОРМУЛИРОВКИ:"
	for _, f := range formulations {
		result += "\n\nПЛОХО: \"" + f.Bad + "\"\nХОРОШО: \"" + f.Good + "\""
	}
	return result
}

// buildExamplesSectionFromData assembles the JSON examples section from DB-loaded question examples.
// The difficulty field is injected at runtime since it is not stored in the DB.
func buildExamplesSectionFromData(difficulty string, multipleChoiceCount int, examples []QuestionExample) string {
	var singleExample, multipleExample string
	for _, ex := range examples {
		if ex.QuestionType == SingleChoice && singleExample == "" {
			singleExample = formatQuestionExample(ex, difficulty)
		}
		if ex.QuestionType == MultipleChoice && multipleExample == "" {
			multipleExample = formatQuestionExample(ex, difficulty)
		}
	}

	if singleExample == "" {
		return "ФОРМАТ ОТВЕТА (строго JSON):\n{\n  \"questions\": []\n}"
	}

	if multipleChoiceCount == 0 || multipleExample == "" {
		return "ФОРМАТ ОТВЕТА (строго JSON):\n{\n  \"questions\": [\n" + singleExample + "\n  ]\n}"
	}
	return "ФОРМАТ ОТВЕТА (строго JSON):\n{\n  \"questions\": [\n" + singleExample + ",\n" + multipleExample + "\n  ]\n}"
}

// formatQuestionExample renders a QuestionExample as an indented JSON string for the prompt.
func formatQuestionExample(ex QuestionExample, difficulty string) string {
	answersJSON := "["
	for i, a := range ex.Answers {
		correctStr := "false"
		if a.IsCorrect {
			correctStr = "true"
		}
		if i > 0 {
			answersJSON += ","
		}
		answersJSON += "\n        {\"text\": \"" + a.Text + "\", \"is_correct\": " + correctStr + "}"
	}
	answersJSON += "\n      ]"

	return `    {
      "question": "` + ex.QuestionText + `",
      "type": "` + string(ex.QuestionType) + `",
      "difficulty": "` + difficulty + `",
      "answers": ` + answersJSON + `,
      "explanation": "` + ex.Explanation + `"
    }`
}
