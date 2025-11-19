import * as yup from 'yup'
import { MAX_FILE_SIZE, SUPPORTED_FORMATS } from './constants'

// Auth validation schemas
export const loginSchema = yup.object({
  email: yup.string().email('Invalid email address').required('Email is required'),
  password: yup.string().min(6, 'Password must be at least 6 characters').required('Password is required'),
})

export const registerSchema = yup.object({
  email: yup.string().email('Invalid email address').required('Email is required'),
  password: yup
    .string()
    .min(6, 'Password must be at least 6 characters')
    .matches(/[a-zA-Z]/, 'Password must contain at least one letter')
    .matches(/[0-9]/, 'Password must contain at least one number')
    .required('Password is required'),
  full_name: yup.string().min(2, 'Name must be at least 2 characters').required('Full name is required'),
  role: yup.string().oneOf(['admin', 'teacher', 'student'], 'Invalid role').required('Role is required'),
})

// Document validation
export const documentUploadSchema = yup.object({
  title: yup.string().min(3, 'Title must be at least 3 characters').required('Title is required'),
  file: yup
    .mixed()
    .required('File is required')
    .test('fileSize', 'File size is too large', value => {
      if (!value) return false
      return (value as File).size <= MAX_FILE_SIZE
    })
    .test('fileType', 'Unsupported file format', value => {
      if (!value) return false
      const extension = (value as File).name.split('.').pop()?.toLowerCase()
      return SUPPORTED_FORMATS.includes(extension || '')
    }),
})

// Test validation
export const testCreationSchema = yup.object({
  title: yup.string().min(3, 'Title must be at least 3 characters').required('Title is required'),
  description: yup.string().optional(),
})

export const testGenerationSchema = yup.object({
  document_id: yup.string().required('Please select a document'),
  title: yup.string().min(3, 'Title must be at least 3 characters').required('Title is required'),
  description: yup.string().optional(),
  num_questions: yup
    .number()
    .min(1, 'At least 1 question is required')
    .max(50, 'Maximum 50 questions allowed')
    .required('Number of questions is required'),
  question_types: yup
    .array()
    .of(yup.string().oneOf(['single_choice', 'multiple_choice', 'true_false', 'short_answer']))
    .min(1, 'Select at least one question type')
    .required('Question types are required'),
  difficulty: yup.string().oneOf(['easy', 'medium', 'hard'], 'Invalid difficulty').required('Difficulty is required'),
})

// Helper functions
export function validateFileSize(file: File): boolean {
  return file.size <= MAX_FILE_SIZE
}

export function validateFileType(file: File): boolean {
  const extension = file.name.split('.').pop()?.toLowerCase()
  return SUPPORTED_FORMATS.includes(extension || '')
}

export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}
