export interface Document {
  id: string
  user_id: string
  user_name?: string  // Only for admin
  user_email?: string // Only for admin
  title: string
  file_name: string
  file_path: string
  file_type: FileType
  file_size: number
  parsed_text?: string
  status: DocumentStatus
  error_msg?: string
  created_at: string
  updated_at: string
}

export const FileType = {
  PDF: 'pdf',
  DOCX: 'docx',
  PPTX: 'pptx',
  TXT: 'txt',
  MD: 'md',
} as const

export type FileType = (typeof FileType)[keyof typeof FileType]

export const DocumentStatus = {
  UPLOADED: 'uploaded',
  PARSING: 'parsing',
  PARSED: 'parsed',
  ERROR: 'error',
} as const

export type DocumentStatus = (typeof DocumentStatus)[keyof typeof DocumentStatus]

export interface DocumentUploadRequest {
  title: string
  file: File
}

export interface DocumentParseRequest {
  document_id: string
}
