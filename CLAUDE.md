# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Система генерации тестовых заданий на основе документов** - информационная система для автоматического создания тестовых вопросов из учебных материалов с интеграцией в Moodle.

**Предметная область**: Образовательные технологии, автоматизация создания учебно-методических материалов

**Применение**: Университеты, образовательные учреждения, онлайн-платформы обучения

**Методология разработки**: Waterfall (каскадная модель)

**Архитектура**: Распределенный монолит (Frontend + Backend REST API)

## Technology Stack

### Backend (Go)

- **Framework**: Fiber v2 (Express-like web framework)
- **ORM**: GORM (PostgreSQL)
- **DI**: Wire (Dependency Injection)
- **Migrations**: golang-migrate
- **Monitoring**: Prometheus metrics export
- **Logging**: Uber Zap (structured logging with request ID tracking)
- **Validation**: go-playground/validator
- **Document parsing**:
  - unidoc/unioffice (DOCX, PPTX)
  - ledongthuc/pdf (PDF)
  - Standard library (TXT)
- **Auth**: JWT tokens (HTTP-only cookies + Authorization header)
- **Testing**: testify, go-mock
- **API Documentation**: Swagger/OpenAPI 2.0 with swag
- **Config**: godotenv

### Frontend (Vue 3)

- **Build tool**: Vite 7
- **Framework**: Vue 3 (Composition API)
- **Language**: TypeScript
- **Routing**: Vue Router 4
- **State management**: Pinia
- **UI Framework**: Tailwind CSS v4 (Utility-first CSS)
- **UI Components**: Headless UI (Accessible components)
- **Icons**: Heroicons (Official Tailwind icons)
- **HTTP Client**: Axios (with request/response logging interceptors)
- **Form validation**: VeeValidate + Yup
- **Testing**: Vitest, Vue Test Utils (112 tests with comprehensive coverage)
- **Logging**: Custom logger utility with DEBUG/INFO/WARN/ERROR levels
- **Code style**: ESLint + Prettier

### ML/AI

- **LLM API**:
  - Perplexity API (для генерации вопросов)
  - OpenAI API (fallback)
  - YandexGPT (российская альтернатива)
- **Local models** (опционально):
  - Ollama (llama2, mistral)
  - Hugging Face Transformers (для экспериментов)

### Database

- **Primary DB**: PostgreSQL 15
- **Caching**: Redis (опционально для будущего масштабирования)

### Infrastructure

- **Containerization**: Docker + Docker Compose
- **Load Balancer**: Nginx
- **Monitoring**: Prometheus + Grafana (опционально)
- **Logging**: Structured logging (zerolog/zap)

### Integration

- **Moodle**: REST API integration (XML-RPC or Web Services)
- **Export formats**: JSON, CSV, Moodle XML

## Architecture

### High-Level Architecture

```
┌─────────────┐
│   Browser   │
└──────┬──────┘
       │ HTTPS
┌──────▼──────────────────┐
│   Nginx (Load Balancer) │
│   - SSL/TLS             │
│   - Static files        │
│   - Reverse proxy       │
└──────┬──────────────────┘
       │
   ┌───┴────┐
   │        │
┌──▼───┐ ┌─▼────────────────┐
│ Vue  │ │ Go Fiber Backend │
│ SPA  │ │ REST API         │
└──────┘ │ - Auth           │
         │ - Document parse │
         │ - Test generation│
         └────┬─────────────┘
              │
    ┌─────────┼─────────┐
    │         │         │
┌───▼─────┐ ┌──▼──┐  ┌───▼────────┐
│ Postgres│ │Redis│  │ LLM API    │
│         │ │cache│  │ Perplexity │
└─────────┘ └─────┘  └────────────┘
```

### Backend Structure (Clean Architecture + DDD)

```
backend/
├── cmd/
│   └── api/
│       └── main.go                 # Entry point
├── internal/
│   ├── domain/                     # Domain layer (entities, value objects)
│   │   ├── entity/
│   │   │   ├── user.go
│   │   │   ├── document.go
│   │   │   ├── test.go
│   │   │   └── question.go
│   │   ├── repository/             # Repository interfaces
│   │   │   ├── user_repository.go
│   │   │   ├── document_repository.go
│   │   │   └── test_repository.go
│   │   └── service/                # Domain services interfaces
│   │       ├── auth_service.go
│   │       ├── generator_service.go
│   │       └── moodle_service.go
│   ├── application/                # Application layer (use cases)
│   │   ├── usecase/
│   │   │   ├── auth/
│   │   │   │   ├── login.go
│   │   │   │   └── register.go
│   │   │   ├── document/
│   │   │   │   ├── upload.go
│   │   │   │   └── parse.go
│   │   │   └── test/
│   │   │       ├── generate.go
│   │   │       ├── export.go
│   │   │       └── moodle_sync.go
│   │   └── dto/                    # Data Transfer Objects
│   │       ├── auth_dto.go
│   │       ├── document_dto.go
│   │       └── test_dto.go
│   ├── infrastructure/             # Infrastructure layer
│   │   ├── persistence/            # Database implementations
│   │   │   ├── postgres/
│   │   │   │   ├── user_repo.go
│   │   │   │   ├── document_repo.go
│   │   │   │   └── test_repo.go
│   │   │   └── migrations/
│   │   │       └── *.sql
│   │   ├── parser/                 # Document parsers
│   │   │   ├── pdf_parser.go
│   │   │   ├── docx_parser.go
│   │   │   ├── pptx_parser.go
│   │   │   └── txt_parser.go
│   │   ├── llm/                    # LLM integrations
│   │   │   ├── perplexity_client.go
│   │   │   └── prompt_builder.go
│   │   ├── moodle/                 # Moodle integration
│   │   │   └── client.go
│   │   └── monitoring/
│   │       └── prometheus.go
│   └── interfaces/                 # Interface layer (controllers, middleware)
│       ├── http/
│       │   ├── handler/
│       │   │   ├── auth_handler.go
│       │   │   ├── document_handler.go
│       │   │   └── test_handler.go
│       │   ├── middleware/
│       │   │   ├── auth.go
│       │   │   ├── logger.go
│       │   │   ├── cors.go
│       │   │   └── error.go
│       │   └── router/
│       │       └── router.go
│       └── validator/
│           └── validator.go
├── pkg/                            # Shared packages
│   ├── config/
│   │   └── config.go
│   ├── logger/
│   │   └── logger.go
│   ├── errors/
│   │   └── errors.go
│   └── utils/
│       └── jwt.go
├── wire.go                         # Wire DI setup
├── wire_gen.go                     # Generated by Wire
├── go.mod
└── go.sum
```

### Frontend Structure (Feature-Based)

```
frontend/
├── public/
│   └── favicon.ico
├── src/
│   ├── assets/                     # Static assets
│   │   ├── styles/
│   │   │   ├── main.scss
│   │   │   └── variables.scss
│   │   └── images/
│   ├── components/                 # Shared components
│   │   ├── common/
│   │   │   ├── AppHeader.vue
│   │   │   ├── AppSidebar.vue
│   │   │   ├── AppButton.vue
│   │   │   └── AppModal.vue
│   │   └── ui/
│   │       ├── FileUploader.vue
│   │       ├── QuestionCard.vue
│   │       └── LoadingSpinner.vue
│   ├── features/                   # Feature modules
│   │   ├── auth/
│   │   │   ├── components/
│   │   │   │   ├── LoginForm.vue
│   │   │   │   └── RegisterForm.vue
│   │   │   ├── composables/
│   │   │   │   └── useAuth.ts
│   │   │   ├── stores/
│   │   │   │   └── authStore.ts
│   │   │   └── types/
│   │   │       └── auth.types.ts
│   │   ├── documents/
│   │   │   ├── components/
│   │   │   │   ├── DocumentUpload.vue
│   │   │   │   ├── DocumentList.vue
│   │   │   │   └── DocumentPreview.vue
│   │   │   ├── composables/
│   │   │   │   └── useDocuments.ts
│   │   │   ├── stores/
│   │   │   │   └── documentsStore.ts
│   │   │   └── types/
│   │   │       └── document.types.ts
│   │   └── tests/
│   │       ├── components/
│   │       │   ├── TestGenerator.vue
│   │       │   ├── QuestionEditor.vue
│   │       │   ├── TestPreview.vue
│   │       │   └── MoodleExport.vue
│   │       ├── composables/
│   │       │   └── useTests.ts
│   │       ├── stores/
│   │       │   └── testsStore.ts
│   │       └── types/
│   │           └── test.types.ts
│   ├── layouts/
│   │   ├── DefaultLayout.vue
│   │   └── AuthLayout.vue
│   ├── router/
│   │   ├── index.ts
│   │   └── guards.ts
│   ├── services/                   # API services
│   │   ├── api.ts                  # Axios instance
│   │   ├── authService.ts
│   │   ├── documentService.ts
│   │   └── testService.ts
│   ├── stores/
│   │   └── index.ts                # Pinia root store
│   ├── types/
│   │   ├── global.d.ts
│   │   └── api.types.ts
│   ├── utils/
│   │   ├── validators.ts
│   │   ├── formatters.ts
│   │   └── constants.ts
│   ├── App.vue
│   └── main.ts
├── .env.example
├── .eslintrc.js
├── .prettierrc
├── index.html
├── package.json
├── tsconfig.json
├── vite.config.ts
└── vitest.config.ts
```

## Structured Logging

### Backend Logging (Go)

#### Implementation

**Library**: Uber Zap - высокопроизводительная библиотека для структурированного логирования

**Location**: `pkg/logger/`

#### Features

1. **Structured Logs**: JSON и console форматы
2. **Log Levels**: debug, info, warn, error
3. **Request Tracking**: Автоматическое добавление request ID к каждому запросу
4. **Context Fields**: Поддержка добавления полей для контекста (user_id, action, etc.)
5. **HTTP Middleware**: Автоматическое логирование всех HTTP запросов с метриками:
   - Method, Path, Status Code
   - Duration (время выполнения)
   - IP address, User-Agent
   - Request ID
   - User ID (если аутентифицирован)

### Usage Examples

```go
// Initialize logger
appLogger, err := logger.New(logger.Config{
    Level:      "info",
    OutputPath: "stdout",
    Format:     "json",
})

// Simple logging
appLogger.Info("Server started", zap.String("port", "8080"))
appLogger.Error("Database error", zap.Error(err))

// With fields
appLogger.InfoWithFields("User logged in", map[string]interface{}{
    "user_id": "123",
    "email": "user@example.com",
    "ip": "192.168.1.1",
})

// Context logger
contextLogger := appLogger.WithField("request_id", "abc123")
contextLogger.Info("Processing request")
```

### HTTP Request Logging Output

```json
{
  "level": "info",
  "timestamp": "2024-01-20T15:04:05.123Z",
  "caller": "logger/middleware.go:45",
  "message": "HTTP request completed",
  "method": "POST",
  "path": "/api/v1/auth/login",
  "status": 200,
  "duration": "45ms",
  "ip": "127.0.0.1",
  "user_agent": "Mozilla/5.0...",
  "request_id": "20240120150405.123456"
}
```

### Configuration

Environment variables:

- `LOG_LEVEL`: debug, info, warn, error (default: info)
- `LOG_FORMAT`: console, json (default: console for dev, json for prod)

#### Backend Testing

Comprehensive test suite in `pkg/logger/logger_test.go`:

- 13 test cases covering all logging functionality
- Tests for log levels, structured fields, file output
- 100% code coverage for core logger functionality

### Frontend Logging (TypeScript)

#### Implementation

**Location**: `frontend/src/utils/logger.ts`

**Purpose**: Unified logging system for debugging frontend operations in development and production

#### Features

1. **Log Levels**: DEBUG, INFO, WARN, ERROR (enum-based)
2. **Environment-aware**: Automatic detection of development/production mode
3. **Browser Console Integration**: All logs output to browser console with proper formatting
4. **HTTP Request/Response Logging**: Automatic logging via Axios interceptors
5. **Store Action Logging**: Pinia store actions logged with payload and errors
6. **Structured Output**: Consistent format with timestamp, level, category, message, and data

#### Architecture

```typescript
enum LogLevel {
  DEBUG = 'DEBUG',
  INFO = 'INFO',
  WARN = 'WARN',
  ERROR = 'ERROR',
}

class Logger {
  // Core logging methods
  debug(message: string, category?: string, data?: unknown): void
  info(message: string, category?: string, data?: unknown): void
  warn(message: string, category?: string, data?: unknown): void
  error(message: string, category?: string, data?: unknown): void

  // HTTP-specific logging
  logRequest(method: string, url: string, data?: unknown): void
  logResponse(method: string, url: string, status: number, data?: unknown): void
  logError(method: string, url: string, error: Error): void

  // Store-specific logging
  logStoreAction(store: string, action: string, payload?: unknown): void
  logStoreError(store: string, action: string, error: Error): void
}
```

#### Integration Points

**1. Axios Interceptors** (`src/services/api.ts`):

```typescript
// Request logging
api.interceptors.request.use(config => {
  logger.logRequest(config.method?.toUpperCase() || 'GET', config.url || '', config.data)
  return config
})

// Response logging
api.interceptors.response.use(
  response => {
    logger.logResponse(method, url, response.status, response.data)
    return response.data
  },
  error => {
    logger.logError(method, url, error)
    throw error
  }
)
```

**2. Pinia Store Actions** (`src/features/auth/stores/authStore.ts`):

```typescript
async function login(credentials: LoginRequest) {
  logger.logStoreAction('authStore', 'login', { email: credentials.email })
  try {
    const response = await authService.login(credentials)
    logger.info('User logged in successfully', 'authStore', { userId: response.user.id })
    return response
  } catch (err: any) {
    logger.logStoreError('authStore', 'login', err)
    throw err
  }
}
```

#### Console Output Example

```
[2024-01-20 15:04:05] [DEBUG] [HTTP] GET /api/v1/auth/me
[2024-01-20 15:04:05] [DEBUG] [HTTP] GET /api/v1/auth/me - 200
[2024-01-20 15:04:06] [DEBUG] [STORE] authStore.login { email: "user@example.com" }
[2024-01-20 15:04:07] [INFO] [authStore] User logged in successfully { userId: "123", role: "student" }
```

#### Configuration

- **Development**: All log levels enabled (DEBUG, INFO, WARN, ERROR)
- **Production**: INFO, WARN, ERROR only (DEBUG disabled)
- **Console Output**: Always enabled for browser DevTools debugging

#### Frontend Testing

Comprehensive test suite in `src/utils/__tests__/logger.spec.ts`:

- 41 test cases covering all logging functionality
- Tests for all log levels (DEBUG, INFO, WARN, ERROR)
- HTTP request/response/error logging tests
- Store action/error logging tests
- Environment detection and level filtering
- Console spy validation for output verification

## Database Schema (ERD - Crow's Foot Notation)

```sql
-- Users (Пользователи системы)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('admin', 'teacher', 'student')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Documents (Загруженные документы)
CREATE TABLE documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(500) NOT NULL,
    file_name VARCHAR(500) NOT NULL,
    file_path VARCHAR(1000) NOT NULL,
    file_type VARCHAR(50) NOT NULL CHECK (file_type IN ('pdf', 'docx', 'pptx', 'txt')),
    file_size BIGINT NOT NULL,
    parsed_text TEXT,
    status VARCHAR(50) DEFAULT 'uploaded' CHECK (status IN ('uploaded', 'parsing', 'parsed', 'error')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Tests (Сгенерированные тесты)
CREATE TABLE tests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    document_id UUID REFERENCES documents(id) ON DELETE SET NULL,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    total_questions INTEGER DEFAULT 0,
    status VARCHAR(50) DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'archived')),
    moodle_synced BOOLEAN DEFAULT FALSE,
    moodle_test_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Questions (Вопросы в тестах)
CREATE TABLE questions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    test_id UUID NOT NULL REFERENCES tests(id) ON DELETE CASCADE,
    question_text TEXT NOT NULL,
    question_type VARCHAR(50) DEFAULT 'single_choice' CHECK (
        question_type IN ('single_choice', 'multiple_choice', 'true_false', 'short_answer')
    ),
    difficulty VARCHAR(50) DEFAULT 'medium' CHECK (difficulty IN ('easy', 'medium', 'hard')),
    points DECIMAL(5,2) DEFAULT 1.0,
    order_num INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Answers (Варианты ответов)
CREATE TABLE answers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    answer_text TEXT NOT NULL,
    is_correct BOOLEAN DEFAULT FALSE,
    order_num INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- User Activity Log (Журнал действий пользователей)
CREATE TABLE activity_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(255) NOT NULL,
    entity_type VARCHAR(100),
    entity_id UUID,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_documents_user_id ON documents(user_id);
CREATE INDEX idx_tests_user_id ON tests(user_id);
CREATE INDEX idx_tests_document_id ON tests(document_id);
CREATE INDEX idx_questions_test_id ON questions(test_id);
CREATE INDEX idx_answers_question_id ON answers(question_id);
CREATE INDEX idx_activity_logs_user_id ON activity_logs(user_id);
CREATE INDEX idx_activity_logs_created_at ON activity_logs(created_at);
```

## Core Features (MVP)

### 1. Аутентификация и авторизация

- Регистрация пользователей (teacher, student, admin)
- Вход/выход (JWT tokens)
- Разграничение прав доступа по ролям
- Журнал действий пользователей

### 2. Управление документами

- Загрузка файлов (PDF, DOCX, PPTX, TXT)
- Валидация формата и размера (max 50MB)
- Парсинг текста из документов
- Просмотр списка загруженных документов
- Удаление документов

### 3. Генерация тестов (LLM)

- Автоматическая генерация вопросов из текста документа
- Настройка параметров:
  - Количество вопросов
  - Типы вопросов - single choice, в будущем (multiple choice, true/false)
- Редактирование сгенерированных вопросов
- Предварительный просмотр теста

### 4. Управление тестами

- CRUD операции с тестами
- Добавление/удаление/редактирование вопросов вручную
- Изменение порядка вопросов
- Статусы тестов (draft, published, archived)

### 5. Интеграция с Moodle

- Экспорт тестов в Moodle XML формат
- Синхронизация с Moodle через REST API
- Отслеживание статуса синхронизации

### 6. Мониторинг и логирование

- Prometheus метрики (requests, latency, errors)
- Структурированные логи
- Аудит действий пользователей

## Design Patterns (минимум 3)

### 1. Repository Pattern

**Где**: `internal/domain/repository/` + `internal/infrastructure/persistence/`

**Назначение**: Абстракция доступа к данным, разделение бизнес-логики и слоя данных

```go
// Domain layer - interface
type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
    FindByEmail(ctx context.Context, email string) (*entity.User, error)
    FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
}

// Infrastructure layer - implementation
type postgresUserRepository struct {
    db *gorm.DB
}
```

### 2. Factory Pattern

**Где**: `internal/infrastructure/parser/`

**Назначение**: Создание парсеров документов в зависимости от типа файла

```go
type DocumentParserFactory interface {
    CreateParser(fileType string) (DocumentParser, error)
}

// Returns PDFParser, DOCXParser, PPTXParser, or TXTParser
```

### 3. Strategy Pattern

**Где**: `internal/infrastructure/llm/`

**Назначение**: Выбор LLM провайдера (Perplexity, OpenAI, YandexGPT)

```go
type LLMStrategy interface {
    GenerateQuestions(ctx context.Context, text string, params GenerationParams) ([]Question, error)
}

type PerplexityStrategy struct {}
type OpenAIStrategy struct {}
type YandexGPTStrategy struct {}
```

### 4. Dependency Injection (Wire)

**Где**: `wire.go`

**Назначение**: Автоматическое внедрение зависимостей

### 5. Middleware Chain

**Где**: `internal/interfaces/http/middleware/`, `pkg/logger/middleware.go`

**Назначение**: Обработка сквозной функциональности (auth, logging, CORS, request tracking)

```go
// Request flow through middleware chain:
app.Use(recover.New())              // 1. Panic recovery
app.Use(logger.RequestIDMiddleware()) // 2. Request ID generation
app.Use(logger.HTTPMiddleware(log))  // 3. Structured logging
app.Use(cors.New(...))               // 4. CORS headers
app.Use(auth.JWTMiddleware())        // 5. Authentication
```

**IMPORTANT: CORS Configuration for HTTP-only Cookies**

When using HTTP-only cookies with `AllowCredentials: true`, you CANNOT use wildcard "*" for `AllowOrigins`. You must explicitly list all allowed origins:

```go
app.Use(cors.New(cors.Config{
    AllowOrigins:     "http://localhost:3000, http://localhost:5173",
    AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
    AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
    AllowCredentials: true, // Required for cookies
}))
```

**Why:**

- HTTP-only cookies require `AllowCredentials: true`
- Browsers block `Access-Control-Allow-Origin: *` with credentials for security
- Must explicitly list each origin (development: localhost ports, production: domain URLs)

### 6. Adapter Pattern (Wrapper Pattern)

**Где**: `pkg/logger/logger.go`

**Назначение**: Обертка над zap.Logger с дополнительным API для удобства использования

```go
type Logger struct {
    *zap.Logger
    config Config
}

// Wraps zap functionality with convenience methods
func (l *Logger) WithFields(fields map[string]interface{}) *Logger
func (l *Logger) InfoWithFields(msg string, fields map[string]interface{})
```

## API Documentation

### Swagger/OpenAPI

**Полная интерактивная документация API доступна через Swagger UI:**

- **URL**: `http://localhost:8080/swagger/index.html`
- **Формат**: OpenAPI 2.0 (Swagger)
- **Генерация**: swag CLI tool
- **Файлы**: `backend/docs/swagger.json`, `backend/docs/swagger.yaml`

**Особенности:**

- Автоматическая генерация из godoc комментариев
- Полное описание всех эндпоинтов с примерами
- Интерактивное тестирование API прямо в браузере
- Документация всех DTO структур
- Описание ошибок и статус кодов
- JWT Bearer Auth поддержка в UI

**Генерация документации:**

```bash
cd backend
swag init -g cmd/api/main.go -o docs
```

**Все эндпоинты задокументированы с:**

- Summary и Description
- Request/Response DTOs
- Security требования (BearerAuth)
- Все возможные HTTP статус коды
- Группировка по тегам (auth, documents, tests, moodle, users)

## API Endpoints (REST)

### Authentication

```
POST   /api/v1/auth/register       # Регистрация
POST   /api/v1/auth/login          # Вход
POST   /api/v1/auth/logout         # Выход
GET    /api/v1/auth/me             # Текущий пользователь
```

### Documents

```
POST   /api/v1/documents           # Загрузка документа
GET    /api/v1/documents           # Список документов
GET    /api/v1/documents/:id       # Детали документа
DELETE /api/v1/documents/:id       # Удаление документа
POST   /api/v1/documents/:id/parse # Парсинг документа
```

### Tests

```
POST   /api/v1/tests               # Создание теста
GET    /api/v1/tests               # Список тестов
GET    /api/v1/tests/:id           # Детали теста
PUT    /api/v1/tests/:id           # Обновление теста
DELETE /api/v1/tests/:id           # Удаление теста
POST   /api/v1/tests/:id/generate  # Генерация вопросов (LLM)
POST   /api/v1/tests/:id/export    # Экспорт в Moodle XML
POST   /api/v1/tests/:id/sync      # Синхронизация с Moodle
```

### Questions

```
POST   /api/v1/tests/:testId/questions           # Добавить вопрос
PUT    /api/v1/tests/:testId/questions/:id       # Обновить вопрос
DELETE /api/v1/tests/:testId/questions/:id       # Удалить вопрос
PUT    /api/v1/tests/:testId/questions/reorder   # Изменить порядок
```

### Monitoring

```
GET    /metrics                    # Prometheus metrics
GET    /health                     # Health check
```

## Commands

### Backend Development

```bash
# Установка зависимостей
cd backend
go mod download

# Генерация Wire DI
go install github.com/google/wire/cmd/wire@latest
wire

# Запуск миграций
migrate -path internal/infrastructure/persistence/migrations \
        -database "postgres://user:pass@localhost:5432/testgen?sslmode=disable" up

# Запуск сервера (development)
go run cmd/api/main.go

# Запуск тестов
go test ./... -v -cover

# Запуск тестов с покрытием
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Линтинг
golangci-lint run

# Сборка
go build -o bin/api cmd/api/main.go
```

### Frontend Development

```bash
# Установка зависимостей
cd frontend
npm install

# Запуск dev сервера
npm run dev

# Сборка для продакшена
npm run build

# Предпросмотр продакшн сборки
npm run preview

# Запуск тестов
npm run test

# Запуск тестов с UI
npm run test:ui

# Проверка типов TypeScript
npm run type-check

# Линтинг
npm run lint

# Форматирование
npm run format
```

### Docker

```bash
# Сборка и запуск всех сервисов
docker-compose up -d

# Пересборка сервисов
docker-compose up -d --build

# Просмотр логов
docker-compose logs -f

# Остановка сервисов
docker-compose down

# Остановка с удалением volumes
docker-compose down -v

# Запуск только БД
docker-compose up -d postgres

# Запуск миграций в контейнере
docker-compose exec backend migrate -path /migrations \
    -database "postgres://user:pass@postgres:5432/testgen?sslmode=disable" up
```

### Database

```bash
# Создание новой миграции
migrate create -ext sql -dir backend/internal/infrastructure/persistence/migrations -seq init_schema

# Применение миграций
migrate -path backend/internal/infrastructure/persistence/migrations \
        -database "postgres://localhost:5432/testgen?sslmode=disable" up

# Откат миграции
migrate -path backend/internal/infrastructure/persistence/migrations \
        -database "postgres://localhost:5432/testgen?sslmode=disable" down 1

# Подключение к PostgreSQL
docker-compose exec postgres psql -U testgen_user -d testgen_db
```

## Environment Variables

### Backend (.env)

```env
# Server
PORT=8080
ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=testgen_user
DB_PASSWORD=testgen_password
DB_NAME=testgen_db
DB_SSLMODE=disable

# JWT
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRATION=24h

# Cookie Configuration
COOKIE_NAME=testgen_token
COOKIE_DOMAIN=
COOKIE_PATH=/
COOKIE_SECURE=false  # Set to true in production with HTTPS
COOKIE_HTTP_ONLY=true
COOKIE_SAME_SITE=Lax

# CORS Configuration (IMPORTANT!)
# ⚠️ When using HTTP-only cookies (credentials), CANNOT use wildcard "*" for AllowOrigins
# Must explicitly list allowed origins and set AllowCredentials=true
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173  # Vite dev server ports
# In production: ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com

# File Upload
MAX_FILE_SIZE=52428800  # 50MB
UPLOAD_DIR=./uploads

# LLM
PERPLEXITY_API_KEY=your-perplexity-key
OPENAI_API_KEY=your-openai-key
YANDEX_GPT_API_KEY=your-yandex-key
LLM_PROVIDER=perplexity  # perplexity, openai, yandex

# Moodle
MOODLE_URL=https://moodle.example.com
MOODLE_TOKEN=your-moodle-webservice-token

# Logging
LOG_LEVEL=info  # debug, info, warn, error
LOG_FORMAT=console  # console (dev) or json (production)

# Monitoring
ENABLE_METRICS=true
```

### Frontend (.env)

```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_MAX_FILE_SIZE=52428800
VITE_SUPPORTED_FORMATS=pdf,docx,pptx,txt
```

## Development Workflow (Waterfall)

### Фаза 1: Анализ (Analysis)

1. Анализ требований к системе генерации тестов
2. Исследование конкурентов (Kahoot, Quizlet, Google Forms)
3. Анализ технологий (LLM модели, парсеры документов, Moodle API)
4. Определение функциональных и нефункциональных требований

### Фаза 2: Проектирование (Design)

1. UML Use Case диаграмма (актёры и их действия)
2. UML Class диаграмма (паттерны проектирования)
3. UML Activity диаграмма (процесс генерации теста)
4. UML Sequence диаграмма (взаимодействие компонентов)
5. ERD диаграмма базы данных (Crow's Foot Notation)
6. Мокап-дизайн интерфейса (Figma)

### Фаза 3: Разработка (Implementation)

1. Настройка инфраструктуры (Docker, PostgreSQL, Nginx)
2. Backend разработка (Go + Fiber)
3. Frontend разработка (Vue + TypeScript)
4. Интеграция с LLM API
5. Интеграция с Moodle
6. Применение паттернов проектирования

### Фаза 4: Тестирование (Testing)

1. Unit тесты (backend + frontend)
2. Integration тесты (API endpoints)
3. E2E тесты (пользовательские сценарии)
4. Тестирование бизнес-процессов
5. Нагрузочное тестирование

### Фаза 5: Внедрение (Deployment)

1. Написание технической документации
2. Инструкция по развертыванию
3. Презентация проекта
4. Защита курсовой работы

## User Roles

### Admin

- Управление пользователями
- Просмотр всех документов и тестов
- Доступ к логам и метрикам

### Teacher

- Загрузка документов
- Генерация и редактирование тестов
- Экспорт в Moodle
- Просмотр своих тестов

### Student (для будущего расширения)

- Прохождение тестов
- Просмотр результатов

## Security Requirements (ГОСТ Р ИСО/МЭК 27001-2012)

1. **Аутентификация**: JWT tokens с истечением срока действия
2. **Авторизация**: Role-based access control (RBAC)
3. **Шифрование**: HTTPS/TLS для всех соединений
4. **Хранение паролей**: bcrypt hashing
5. **Валидация входных данных**: защита от SQL injection, XSS
6. **Журналирование**: аудит всех действий пользователей
7. **Ограничение файлов**: проверка типов и размеров файлов

## UI/UX Requirements (ГОСТ Р ИСО 9241)

1. **Единообразие**: Material Design (Vuetify)
2. **Обратная связь**: Loading states, success/error messages
3. **Доступность**: Keyboard navigation, ARIA attributes
4. **Интуитивность**: Понятная навигация, подсказки
5. **Адаптивность**: Responsive дизайн (desktop, tablet)
6. **Цветовая палитра**: Контрастные цвета, читаемость текста

## Testing Strategy

### Backend Tests

- Unit tests: Domain entities, services
- Repository tests: Database operations (with testcontainers)
- Handler tests: HTTP endpoints (mocked dependencies)
- Integration tests: Full API workflow

### Frontend Tests

**Test Structure:** All tests follow co-location pattern in `__tests__/` directories with `.spec.ts` naming convention.

**Current Coverage: 112 tests passing**

#### Test Files

1. **Component Tests** (`src/features/auth/components/__tests__/`):
   - `LoginForm.spec.ts` - 6 tests (form rendering, submission, error handling, loading states, navigation)
   - `RegisterForm.spec.ts` - 5 tests (form without role selection, submission, error display, loading states, links)

2. **Store Tests** (`src/features/auth/stores/__tests__/`):
   - `authStore.spec.ts` - 24 tests covering:
     - Positive: login, register, logout, user fetch, initialization
     - Negative: invalid credentials, duplicate email, weak password, network errors, token errors, concurrent operations

3. **Service Tests** (`src/services/__tests__/`):
   - `authService.spec.ts` - 25 tests covering:
     - Positive: login, register, logout, getMe
     - Negative: validation errors, network errors, unauthorized access, special characters, SQL injection attempts, XSS attempts, long inputs

4. **Utility Tests** (`src/utils/__tests__/`):
   - `logger.spec.ts` - 41 tests (all log levels, HTTP logging, store logging, fields, formatters)
   - `validators.spec.ts` - 5 tests (file size, file type, formatters)
   - `formatters.spec.ts` - 6 tests (date formatting, relative time, text truncation)

#### Key Testing Patterns

- **Mock-based testing** with Vitest `vi.mock()`
- **User registration without role** - role assigned by backend (default: student)
- **AuthResponse includes token field** - required for JWT authentication
- **Comprehensive negative testing** - 60+ negative test cases for security and edge cases
- **localStorage handling** - tests account for undefined/null differences in test environment

#### Coverage Target

- Backend: >70%
- Frontend: >60% (currently achieved with 112 tests)

## Performance Requirements

- API response time: <500ms (95th percentile)
- Document parsing: <10s for files up to 50MB
- Test generation: <30s for 20 questions
- Concurrent users: 50+ (MVP)

## Future Enhancements (Магистерская диссертация)

1. **ML модели**: Fine-tuned Transformers для генерации вопросов
2. **Микросервисы**: Разделение на независимые сервисы
3. **Масштабирование**: Kubernetes, horizontal scaling
4. **Расширенная аналитика**: Статистика по тестам, сложности
5. **Адаптивные тесты**: Подбор вопросов по уровню знаний
6. **Мобильное приложение**: React Native
7. **Расширенная интеграция**: Google Classroom, Microsoft Teams

## Important Notes

- Весь код должен быть на английском (комментарии, переменные, функции)
- Коммиты в Git на английском
- Документация может быть на русском (для курсовой)
- Следовать SOLID принципам
- Избегать OWASP Top 10 уязвимостей
- Код должен проходить линтеры (golangci-lint, ESLint)
- Обязательное покрытие тестами критичных функций

## Git Workflow

```bash
# Feature branch workflow
git checkout -b feature/document-upload
git add .
git commit -m "feat: implement document upload functionality"
git push origin feature/document-upload

# Conventional commits
feat: новая функциональность
fix: исправление бага
docs: документация
refactor: рефакторинг
test: тесты
chore: инфраструктура
```
