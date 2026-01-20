// Export format types
export enum ExportFormat {
  JSON = 'json',
  MOODLE_XML = 'moodle_xml',
  STEPIK_CSV = 'stepik_csv',
  GIFT = 'gift',
  AIKEN = 'aiken',
  BLACKBOARD = 'blackboard',
  QTI = 'qti',
}

export interface FormatInfo {
  format: ExportFormat
  name: string
  description: string
}

export interface ExportRequest {
  testId: string
  format: ExportFormat
}
