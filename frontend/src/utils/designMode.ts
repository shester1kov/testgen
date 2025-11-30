/**
 * Design Mode Utilities
 *
 * Provides mock authentication and data for UI design tools (Anima, Figma)
 * When VITE_DESIGN_MODE=true, the app bypasses real authentication
 */

export const isDesignMode = (): boolean => {
  return import.meta.env.VITE_DESIGN_MODE === 'true'
}

/**
 * Mock user for design mode
 * Returns a fake authenticated user without backend communication
 */
export const getMockUser = () => {
  return {
    id: 'design-mode-user-123',
    email: 'admin@example.com',
    full_name: 'Администратор (Design Mode)',
    role: 'admin' as const
  }
}

/**
 * Mock authentication response
 */
export const getMockAuthResponse = () => {
  return {
    token: 'design-mode-mock-token-xyz',
    user: getMockUser()
  }
}

/**
 * Mock documents for design mode
 */
export const getMockDocuments = () => {
  return [
    {
      id: '1',
      user_id: 'design-mode-user-123',
      title: 'Лекция 1: Введение в программирование',
      file_name: 'lecture_01.pdf',
      file_type: 'pdf',
      file_size: 1024000,
      status: 'parsed',
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString()
    },
    {
      id: '2',
      user_id: 'design-mode-user-123',
      title: 'Лекция 2: Основы Python',
      file_name: 'lecture_02.docx',
      file_type: 'docx',
      file_size: 512000,
      status: 'parsed',
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString()
    },
    {
      id: '3',
      user_id: 'design-mode-user-123',
      title: 'Презентация: Структуры данных',
      file_name: 'presentation.pptx',
      file_type: 'pptx',
      file_size: 2048000,
      status: 'uploading',
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString()
    }
  ]
}

/**
 * Mock tests for design mode
 */
export const getMockTests = () => {
  return [
    {
      id: '1',
      user_id: 'design-mode-user-123',
      title: 'Тест по основам программирования',
      description: 'Проверка знаний базовых концепций',
      total_questions: 15,
      status: 'published',
      moodle_synced: true,
      created_at: new Date().toISOString()
    },
    {
      id: '2',
      user_id: 'design-mode-user-123',
      title: 'Тест по Python',
      description: 'Синтаксис и основные библиотеки',
      total_questions: 20,
      status: 'draft',
      moodle_synced: false,
      created_at: new Date().toISOString()
    }
  ]
}

/**
 * Mock users for design mode
 */
export const getMockUsers = () => {
  return [
    {
      id: '1',
      email: 'admin@example.com',
      full_name: 'Иван Администратор',
      role: 'admin'
    },
    {
      id: '2',
      email: 'teacher@example.com',
      full_name: 'Мария Преподаватель',
      role: 'teacher'
    },
    {
      id: '3',
      email: 'student1@example.com',
      full_name: 'Петр Студентов',
      role: 'student'
    },
    {
      id: '4',
      email: 'student2@example.com',
      full_name: 'Анна Учащаяся',
      role: 'student'
    }
  ]
}

/**
 * Log design mode status on app initialization
 */
export const logDesignModeStatus = () => {
  if (isDesignMode()) {
    console.log(
      '%c🎨 DESIGN MODE ENABLED',
      'background: #ff6b00; color: white; padding: 8px 16px; font-size: 14px; font-weight: bold;'
    )
    console.log(
      '%cAuthentication bypassed. Using mock data for UI design.',
      'color: #ff6b00; font-size: 12px;'
    )
  }
}
