export interface Test {
  id: string
  user_id: string
  user_name?: string  // Only for admin
  user_email?: string // Only for admin
  document_id?: string
  title: string
  description?: string
  total_questions: number
  status: TestStatus
  moodle_synced: boolean
  moodle_test_id?: string
  created_at: string
  updated_at: string
  questions?: Question[]
}

export enum TestStatus {
  DRAFT = 'draft',
  PUBLISHED = 'published',
  ARCHIVED = 'archived',
}

export interface Question {
  id: string
  test_id: string
  question_text: string
  question_type: QuestionType
  difficulty: Difficulty
  points: number
  order_num: number
  created_at: string
  updated_at: string
  answers: Answer[]
}

export enum QuestionType {
  SINGLE_CHOICE = 'single_choice',
  MULTIPLE_CHOICE = 'multiple_choice',
  TRUE_FALSE = 'true_false',
  SHORT_ANSWER = 'short_answer',
}

export enum Difficulty {
  EASY = 'easy',
  MEDIUM = 'medium',
  HARD = 'hard',
}

export interface Answer {
  id: string
  question_id: string
  answer_text: string
  is_correct: boolean
  order_num: number
  created_at: string
}

export interface TestGenerationRequest {
  document_id: string
  title: string
  description?: string
  num_questions: number
  question_types: QuestionType[]
  difficulty: Difficulty
}

export interface TestExportRequest {
  test_id: string
  format: 'json' | 'csv' | 'moodle_xml'
}

export interface MoodleSyncRequest {
  test_id: string
  moodle_course_id: string
}
