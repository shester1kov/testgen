export const MAX_FILE_SIZE = parseInt(import.meta.env.VITE_MAX_FILE_SIZE) || 52428800 // 50MB
export const SUPPORTED_FORMATS = import.meta.env.VITE_SUPPORTED_FORMATS?.split(',') || [
  'pdf',
  'docx',
  'pptx',
  'txt',
]

export const FILE_TYPE_LABELS: Record<string, string> = {
  pdf: 'PDF Document',
  docx: 'Word Document',
  pptx: 'PowerPoint Presentation',
  txt: 'Text File',
}

export const DOCUMENT_STATUS_LABELS: Record<string, string> = {
  uploaded: 'Uploaded',
  parsing: 'Parsing...',
  parsed: 'Parsed',
  error: 'Error',
}

export const TEST_STATUS_LABELS: Record<string, string> = {
  draft: 'Draft',
  published: 'Published',
  archived: 'Archived',
}

export const QUESTION_TYPE_LABELS: Record<string, string> = {
  single_choice: 'Single Choice',
  multiple_choice: 'Multiple Choice',
  true_false: 'True/False',
  short_answer: 'Short Answer',
}

export const DIFFICULTY_LABELS: Record<string, string> = {
  easy: 'Easy',
  medium: 'Medium',
  hard: 'Hard',
}

export const DIFFICULTY_COLORS: Record<string, string> = {
  easy: 'green',
  medium: 'yellow',
  hard: 'red',
}
