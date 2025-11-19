export interface Document {
  id: string
  user_id: string
  title: string
  file_name: string
  file_path: string
  file_type: FileType
  file_size: number
  parsed_text?: string
  status: DocumentStatus
  created_at: string
  updated_at: string
}

export enum FileType {
  PDF = 'pdf',
  DOCX = 'docx',
  PPTX = 'pptx',
  TXT = 'txt',
}

export enum DocumentStatus {
  UPLOADED = 'uploaded',
  PARSING = 'parsing',
  PARSED = 'parsed',
  ERROR = 'error',
}

export interface DocumentUploadRequest {
  title: string
  file: File
}

export interface DocumentParseRequest {
  document_id: string
}
