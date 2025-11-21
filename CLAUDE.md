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
- **Validation**: go-playground/validator
- **Document parsing**:
  - unidoc/unioffice (DOCX, PPTX)
  - ledongthuc/pdf (PDF)
  - Standard library (TXT)
- **Auth**: JWT tokens
- **Testing**: testify, go-mock
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
- **HTTP Client**: Axios
- **Form validation**: VeeValidate + Yup
- **Testing**: Vitest, Vue Test Utils
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
  - Типы вопросов (single choice, multiple choice, true/false)
  - Уровень сложности
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
**Где**: `internal/interfaces/http/middleware/`

**Назначение**: Обработка сквозной функциональности (auth, logging, CORS)

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
- Component tests: Vue components (Vitest + Vue Test Utils)
- Store tests: Pinia stores
- E2E tests: User flows (опционально Playwright)

### Coverage Target
- Backend: >70%
- Frontend: >60%

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

## Contact & Resources

- Instructor: Макиевский С.Е.
- University: МИРЭА - Институт перспективных технологий и индустриального программирования
- Course: Создание программного обеспечения
- Year: 2024-2025
