# План исследовательской реализации Test Generation System

## 1. Основная идея проекта

### 1.1. Что разрабатываем

**Test Generation System** — веб‑платформа для автоматической генерации тестовых заданий из учебных материалов с использованием технологий искусственного интеллекта (LLM) и интеграцией с Learning Management System Moodle.

**Полное название**: Система автоматической генерации тестовых заданий на основе документов с интеграцией в Moodle

**Тип системы**: Информационная система (распределенный монолит)

### 1.2. Проблема, которую решает проект

**Основная проблема**: Преподаватели тратят значительное время на ручное создание тестовых заданий из учебных материалов. Этот процесс:

- Отнимает 3-5 часов на один тест из 20 вопросов
- Подвержен человеческим ошибкам (опечатки, некорректные формулировки)
- Не масштабируется при большом количестве дисциплин
- Требует экспертизы в предметной области И педагогическом дизайне

**Дополнительные проблемы**:

- Отсутствие инструментов для быстрой генерации разнообразных вопросов
- Сложность интеграции тестов в существующие LMS
- Необходимость валидации качества вопросов перед использованием

### 1.3. Целевая аудитория

**Первичная аудитория**:

- **Преподаватели вузов и колледжей** — основные пользователи, создающие тесты для своих дисциплин
- **Методисты** — специалисты, разрабатывающие учебно-методические материалы
- **Администраторы LMS** — управление системой и пользователями

**Вторичная аудитория**:

- **Ассистенты преподавателей** — помощь в подготовке тестовых материалов
- **Студенты** (перспектива) — прохождение тестов для валидации качества

**География применения**: Университеты, образовательные учреждения, онлайн-платформы обучения

### 1.4. Роли пользователей и их задачи

#### Admin (Администратор)

**Зона ответственности**: Управление системой и пользователями

**Основные задачи**:

- Управление пользователями (просмотр списка, изменение ролей)
- Назначение ролей (admin, teacher, student)
- Контроль аудит-логов (журнал действий пользователей)
- Конфигурация интеграций (Moodle, LLM провайдеры)
- Мониторинг системы (Grafana, Prometheus)
- Просмотр всех документов и тестов в системе

**Доступ**: Все эндпоинты API, включая `/api/v1/users/:id/role`

#### Teacher (Преподаватель)

**Зона ответственности**: Создание и управление учебными материалами

**Основные задачи**:

- Загрузка учебных материалов (PDF, DOCX, PPTX, TXT, MD)
- Парсинг текста из документов
- Генерация тестовых вопросов с помощью LLM
- Редактирование сгенерированных вопросов
- Настройка параметров генерации (количество вопросов, типы, сложность)
- Экспорт тестов в Moodle XML формат
- Синхронизация с Moodle
- Просмотр метрик качества вопросов
- Управление своими документами и тестами

**Доступ**: Эндпоинты документов, тестов, Moodle интеграции

#### Student (Студент)

**Зона ответственности**: Взаимодействие с тестами (будущая функциональность)

**Планируемые задачи**:

- Прохождение сгенерированных тестов
- Валидация качества вопросов через обратную связь
- Просмотр результатов тестирования

**Доступ**: Только назначенные тесты (read-only)

**Примечание**: При регистрации пользователь автоматически получает роль `student`. Только администратор может изменять роли.

## 2. Анализ требований к проекту

### 2.1. Ожидания пользователей

**Преподаватели ожидают**:

- Быструю генерацию корректных вопросов из разнородных форматов файлов (PDF, DOCX, PPTX, TXT, MD)
- Возможность ручной доработки вопросов перед экспортом (UI для правок и валидаций)
- Надёжную интеграцию с Moodle (XML импорт, настройка банков вопросов)
- Прозрачную оценку качества: предпросмотр, метрики сложности, истории генераций
- Интуитивный интерфейс без необходимости технических знаний

**Администраторы ожидают**:

- Полный контроль над системой и пользователями
- Мониторинг работоспособности и производительности
- Аудит действий пользователей для соблюдения политик безопасности
- Гибкую настройку интеграций (LLM провайдеры, Moodle)

**Студенты ожидают** (будущая функциональность):

- Удобный интерфейс для прохождения тестов
- Немедленную обратную связь по результатам

### 2.2. Функциональные требования

#### FR-1: Аутентификация и авторизация

- **FR-1.1**: Регистрация пользователей с автоматическим назначением роли `student`
- **FR-1.2**: Вход в систему с JWT токенами в HTTP-only cookies (защита от XSS)
- **FR-1.3**: Поддержка двух способов авторизации: Cookie (приоритет) и Authorization header
- **FR-1.4**: Выход из системы с очисткой cookie
- **FR-1.5**: RBAC (Role-Based Access Control) на основе таблицы roles в БД

**Реализовано**:

- Middleware для проверки JWT токенов
- Middleware для проверки ролей пользователей
- Эндпоинты: `/auth/register`, `/auth/login`, `/auth/logout`, `/auth/me`

#### FR-2: Управление пользователями (Admin only)

- **FR-2.1**: Просмотр списка всех пользователей с пагинацией
- **FR-2.2**: Изменение ролей пользователей (admin/teacher/student)
- **FR-2.3**: Защита endpoint'ов middleware для проверки роли администратора

**Реализовано**:

- `/api/v1/users` (GET) - список пользователей
- `/api/v1/users/:id/role` (PUT) - изменение роли

#### FR-3: Управление документами

- **FR-3.1**: Загрузка документов форматов PDF, DOCX, PPTX, TXT, MD (до 50MB)
- **FR-3.2**: Автоматический парсинг текста из документов
- **FR-3.3**: Просмотр списка загруженных документов с пагинацией
- **FR-3.4**: Просмотр деталей документа (метаданные, parsed text)
- **FR-3.5**: Удаление документов
- **FR-3.6**: Валидация типа и размера файла

**Реализовано**:

- Parser factory с поддержкой 5 форматов
- Хранение файлов в `./uploads`
- Статусы документов: `uploaded`, `parsing`, `parsed`, `error`

#### FR-4: Генерация тестов с LLM

- **FR-4.1**: Генерация вопросов из текста документа с помощью LLM
- **FR-4.2**: Настройка параметров генерации:
  - Количество вопросов (1-50)
  - Типы вопросов (single_choice, multiple_choice, true_false, short_answer)
  - Уровень сложности (easy, medium, hard)
- **FR-4.3**: Поддержка нескольких LLM провайдеров (Perplexity, OpenAI, YandexGPT)
- **FR-4.4**: Создание теста с автоматической генерацией вопросов

**Реализовано**:

- Strategy Pattern для выбора LLM провайдера
- Factory Pattern для создания LLM стратегий
- Промпты для генерации качественных вопросов

#### FR-5: Управление тестами

- **FR-5.1**: CRUD операции с тестами
- **FR-5.2**: Редактирование вопросов и ответов
- **FR-5.3**: Изменение порядка вопросов
- **FR-5.4**: Статусы тестов (draft, published, archived)
- **FR-5.5**: Просмотр списка тестов с пагинацией

**Реализовано**:

- Полный REST API для управления тестами
- Связь тестов с документами (опционально)
- Автоматический подсчет total_questions

#### FR-6: Интеграция с Moodle

- **FR-6.1**: Экспорт тестов в Moodle XML формат
- **FR-6.2**: Синхронизация тестов с Moodle через REST API
- **FR-6.3**: Получение списка курсов Moodle
- **FR-6.4**: Проверка подключения к Moodle
- **FR-6.5**: Отслеживание статуса синхронизации (moodle_synced, moodle_test_id)

**Реализовано**:

- XML экспортер для Moodle формата
- Moodle client для работы с Web Services API
- Эндпоинты: `/tests/:id/export-xml`, `/tests/:id/sync-moodle`, `/moodle/courses`, `/moodle/validate`

#### FR-7: Мониторинг и логирование

- **FR-7.1**: Экспорт метрик в формате Prometheus
- **FR-7.2**: Структурированное логирование в JSON формате
- **FR-7.3**: Трекинг request ID для всех запросов
- **FR-7.4**: Логирование всех HTTP запросов с метриками (duration, status, path)
- **FR-7.5**: Аудит действий пользователей

**Реализовано**:

- Prometheus metrics endpoint (`/metrics`)
- Health check endpoint (`/health`)
- Uber Zap logger с JSON форматом
- Middleware для request ID и HTTP логирования
- Grafana dashboards для визуализации

### 2.3. Нефункциональные требования

#### NFR-1: Производительность

- API response time: <500ms (95th percentile)
- Парсинг документов: <10s для файлов до 50MB
- Генерация теста: <30s для 20 вопросов
- Поддержка 50+ одновременных пользователей

#### NFR-2: Безопасность

- JWT токены с истечением срока действия (24h)
- Bcrypt hashing для паролей (cost 10)
- HTTPS/TLS для всех соединений
- HTTP-only cookies для защиты от XSS
- CORS с явным указанием разрешенных origins
- Валидация входных данных на всех слоях
- SQL injection защита через GORM prepared statements
- XSS защита через санитизацию HTML

#### NFR-3: Масштабируемость

- Горизонтальное масштабирование backend (stateless design)
- Поддержка Redis для кэширования (опционально)
- Пагинация для всех списковых эндпоинтов
- Оптимизация запросов к БД с индексами

#### NFR-4: Надежность

- Покрытие тестами ≥80% для критических модулей
- Graceful shutdown для корректного завершения запросов
- Retry механизмы для интеграций (LLM, Moodle)
- Логирование ошибок для анализа

#### NFR-5: Usability (Эргономика)

- Адаптация под ГОСТ Р ИСО 9241 (UI/UX)
- Responsive дизайн (desktop 1280px+)
- Контрастные цвета для читаемости
- Клавиатурная навигация
- Понятные сообщения об ошибках

## 3. Архитектура проекта

### 3.1. Общая архитектура системы

**Архитектурный стиль**: Распределенный монолит (Distributed Monolith)

**Архитектурный подход**: Clean Architecture + Domain-Driven Design (DDD)

```
┌─────────────────────┐
│      Browser        │
└──────────┬──────────┘
           │ HTTPS
┌──────────▼──────────────────────┐
│   Nginx (Load Balancer)         │
│   - SSL/TLS терминация          │
│   - Reverse proxy               │
│   - Статические файлы frontend  │
└──────────┬──────────────────────┘
           │
     ┌─────┴──────┐
     │            │
┌────▼────┐  ┌───▼──────────────────────┐
│  Vue 3  │  │ Go Fiber Backend         │
│   SPA   │  │ REST API                 │
│ (Vite)  │  │ - Auth (JWT)             │
└─────────┘  │ - Document parsing       │
             │ - Test generation (LLM)  │
             │ - Moodle integration     │
             │ - Monitoring (Prometheus)│
             └────┬─────────────────────┘
                  │
     ┌────────────┼────────────┬──────────────┐
     │            │            │              │
┌────▼────┐  ┌───▼────┐  ┌───▼──────┐  ┌───▼─────────┐
│Postgres │  │ Redis  │  │ LLM API  │  │File Storage │
│   15    │  │   7    │  │Perplexity│  │  ./uploads  │
│         │  │(cache) │  │OpenAI    │  │             │
│         │  │        │  │YandexGPT │  │             │
└─────────┘  └────────┘  └──────────┘  └─────────────┘
     │
     │ metrics & logs
     │
┌────▼──────────────────────────────┐
│  Observability Stack              │
│  - Prometheus (metrics)           │
│  - Grafana (dashboards)           │
│  - Loki (log aggregation)         │
│  - Promtail (log collector)       │
└───────────────────────────────────┘
```

### 3.2. Backend архитектура (Clean Architecture)

#### Слоистая архитектура

```
┌──────────────────────────────────────────────────────────┐
│                  Interfaces Layer                        │
│  (HTTP Handlers, Middleware, Router, Validators)         │
│  - Fiber HTTP handlers                                   │
│  - JWT & Role middleware                                 │
│  - Request validation                                    │
└──────────────────┬───────────────────────────────────────┘
                   │
┌──────────────────▼───────────────────────────────────────┐
│              Application Layer                           │
│  (Use Cases, DTOs, Business Logic)                       │
│  - Login/Register use cases                              │
│  - Document upload/parse use cases                       │
│  - Test generation/export use cases                      │
└──────────────────┬───────────────────────────────────────┘
                   │
┌──────────────────▼───────────────────────────────────────┐
│                 Domain Layer                             │
│  (Entities, Repository Interfaces, Business Rules)       │
│  - User, Document, Test, Question, Answer entities       │
│  - Repository interfaces (contracts)                     │
│  - Domain validation logic                               │
└──────────────────┬───────────────────────────────────────┘
                   │
┌──────────────────▼───────────────────────────────────────┐
│            Infrastructure Layer                          │
│  (DB, Parsers, LLM, Moodle, Logging, Monitoring)         │
│  - PostgreSQL repositories (GORM)                        │
│  - Document parsers (PDF, DOCX, PPTX, TXT, MD)           │
│  - LLM strategies (Perplexity, OpenAI, YandexGPT)        │
│  - Moodle client & XML exporter                          │
│  - Uber Zap logger                                       │
│  - Prometheus metrics                                    │
└──────────────────────────────────────────────────────────┘
```

#### Структура директорий Backend

```
backend/
├── cmd/api/                    # Entry point
│   └── main.go                 # App initialization, Wire DI
├── internal/
│   ├── domain/                 # Domain Layer (чистая бизнес-логика)
│   │   ├── entity/             # Domain entities
│   │   │   ├── user.go         # User entity с методами (IsAdmin, SetPassword)
│   │   │   ├── role.go         # Role entity (admin, teacher, student)
│   │   │   ├── document.go     # Document entity
│   │   │   ├── test.go         # Test entity
│   │   │   ├── question.go     # Question entity
│   │   │   └── answer.go       # Answer entity
│   │   └── repository/         # Repository interfaces (контракты)
│   │       ├── user_repository.go
│   │       ├── document_repository.go
│   │       ├── test_repository.go
│   │       ├── question_repository.go
│   │       └── answer_repository.go
│   ├── application/            # Application Layer (use cases)
│   │   ├── dto/                # Data Transfer Objects
│   │   │   ├── auth_dto.go
│   │   │   ├── document_dto.go
│   │   │   ├── test_dto.go
│   │   │   ├── user_dto.go
│   │   │   └── response_dto.go
│   │   └── usecase/            # Use cases (бизнес-логика)
│   │       ├── auth/
│   │       │   ├── login.go
│   │       │   └── register.go
│   │       ├── document/
│   │       │   ├── upload.go
│   │       │   ├── parse.go
│   │       │   ├── get.go
│   │       │   ├── list.go
│   │       │   └── delete.go
│   │       └── test/
│   │           ├── generate.go      # LLM генерация
│   │           └── export_moodle.go # Экспорт в Moodle XML
│   ├── infrastructure/         # Infrastructure Layer
│   │   ├── persistence/
│   │   │   ├── postgres/       # GORM репозитории
│   │   │   │   ├── database.go
│   │   │   │   ├── user_repo.go
│   │   │   │   ├── role_repo.go
│   │   │   │   ├── document_repo.go
│   │   │   │   ├── test_repo.go
│   │   │   │   ├── question_repo.go
│   │   │   │   └── answer_repo.go
│   │   │   └── migrations/     # SQL миграции
│   │   │       ├── 000001_init_schema.up.sql
│   │   │       ├── 000002_create_roles_table.up.sql
│   │   │       └── 000003_add_error_msg_and_md_support.up.sql
│   │   ├── parser/             # Document parsers
│   │   │   ├── parser.go       # Parser interface
│   │   │   ├── pdf_parser.go
│   │   │   ├── docx_parser.go
│   │   │   ├── pptx_parser.go
│   │   │   ├── txt_parser.go
│   │   │   └── md_parser.go
│   │   ├── llm/                # LLM integrations
│   │   │   ├── llm_strategy.go      # Strategy interface
│   │   │   ├── factory.go           # Factory для создания стратегий
│   │   │   ├── perplexity_strategy.go
│   │   │   ├── openai_strategy.go
│   │   │   └── yandex_strategy.go
│   │   └── moodle/             # Moodle integration
│   │       ├── client.go            # REST API клиент
│   │       └── xml_exporter.go      # XML формат экспорта
│   └── interfaces/             # Interface Layer
│       ├── http/
│       │   ├── handler/        # HTTP handlers (controllers)
│       │   │   ├── auth_handler.go
│       │   │   ├── document_handler.go
│       │   │   ├── test_handler.go
│       │   │   ├── user_handler.go
│       │   │   ├── moodle_handler.go
│       │   │   └── stats_handler.go
│       │   ├── middleware/     # Middleware цепочка
│       │   │   ├── auth.go     # JWT authentication
│       │   │   └── role.go     # RBAC authorization
│       │   └── router/
│       │       └── router.go   # Fiber routes configuration
│       └── validator/
│           └── validator.go    # Input validation
├── pkg/                        # Shared packages
│   ├── config/                 # Configuration management
│   │   └── config.go           # Env variables loader
│   ├── logger/                 # Structured logging
│   │   ├── logger.go           # Zap logger wrapper
│   │   └── middleware.go       # HTTP logging middleware
│   ├── monitoring/             # Prometheus metrics
│   │   └── prometheus.go
│   ├── security/               # Security utilities
│   │   └── xss.go              # XSS sanitization
│   └── errors/                 # Error handling
│       └── errors.go
├── docs/                       # Swagger documentation
│   ├── swagger.json
│   └── swagger.yaml
├── wire.go                     # Wire DI configuration
└── wire_gen.go                 # Generated by Wire
```

### 3.3. Frontend архитектура (Feature-based)

#### Структура директорий Frontend

```
frontend/src/
├── features/                   # Feature modules (domain-oriented)
│   ├── auth/                   # Authentication feature
│   │   ├── components/
│   │   │   ├── LoginForm.vue
│   │   │   └── RegisterForm.vue
│   │   ├── stores/
│   │   │   └── authStore.ts    # Pinia store
│   │   └── types/
│   │       └── auth.types.ts
│   ├── documents/              # Documents management
│   │   ├── components/
│   │   │   ├── DocumentUpload.vue
│   │   │   ├── DocumentList.vue
│   │   │   └── DocumentCard.vue
│   │   ├── services/
│   │   │   └── documentService.ts
│   │   ├── stores/
│   │   │   └── documentsStore.ts
│   │   └── types/
│   │       └── document.types.ts
│   ├── tests/                  # Tests management
│   │   ├── components/
│   │   │   ├── TestCard.vue
│   │   │   ├── QuestionList.vue
│   │   │   └── QuestionEditModal.vue
│   │   ├── stores/
│   │   │   └── testsStore.ts
│   │   └── types/
│   │       └── test.types.ts
│   └── users/                  # User management (Admin)
│       ├── components/
│       │   └── UserTable.vue
│       └── stores/
│           └── usersStore.ts
├── layouts/                    # Page layouts
│   ├── DefaultLayout.vue       # Main app layout (с header/sidebar)
│   └── AuthLayout.vue          # Auth pages layout
├── router/                     # Vue Router
│   ├── index.ts                # Routes configuration
│   └── guards.ts               # Navigation guards (auth check)
├── services/                   # API services
│   ├── api.ts                  # Axios instance с interceptors
│   ├── authService.ts
│   ├── documentService.ts
│   ├── testService.ts
│   ├── userService.ts
│   └── statsService.ts
├── stores/                     # Global Pinia stores
│   └── index.ts
├── utils/                      # Utility functions
│   ├── logger.ts               # Frontend logger (DEBUG/INFO/WARN/ERROR)
│   ├── formatters.ts           # Date/number formatters
│   ├── validators.ts           # File validators
│   └── constants.ts
├── views/                      # Page components
│   ├── DashboardView.vue
│   ├── DocumentsView.vue
│   ├── DocumentDetailsView.vue
│   ├── TestsView.vue
│   ├── TestDetailsView.vue
│   ├── CreateTestView.vue
│   ├── EditTestView.vue
│   ├── UsersView.vue
│   └── NotFoundView.vue
└── App.vue                     # Root component
```

### 3.4. Database Schema (PostgreSQL)

**Основные таблицы**:

```sql
-- Роли пользователей
users (id, email, password_hash, full_name, role_id, created_at, updated_at, deleted_at)
roles (id, name, description, created_at)

-- Документы
documents (id, user_id, title, file_name, file_path, file_type, file_size,
           parsed_text, status, error_message, created_at, updated_at, deleted_at)

-- Тесты и вопросы
tests (id, user_id, document_id, title, description, total_questions,
       status, moodle_synced, moodle_test_id, created_at, updated_at, deleted_at)

questions (id, test_id, question_text, question_type, difficulty, points,
           order_num, created_at, updated_at)

answers (id, question_id, answer_text, is_correct, order_num, created_at)

-- Аудит
activity_logs (id, user_id, action, entity_type, entity_id, ip_address,
               user_agent, created_at)
```

**Индексы для производительности**:

- `users(email)` - UNIQUE для быстрой аутентификации
- `documents(user_id)` - для списка документов пользователя
- `tests(user_id, document_id)` - для фильтрации тестов
- `questions(test_id, order_num)` - для упорядоченного получения вопросов
- `answers(question_id)` - для получения вариантов ответов

### 3.5. Компоненты системы и их взаимодействие

#### Компоненты

1. **Frontend (Vue 3 SPA)**
   - UI для загрузки материалов
   - Редактор вопросов
   - Просмотр метрик и статистики
   - Управление пользователями (Admin)

2. **Backend (Go + Fiber)**
   - REST API (Swagger документация)
   - Оркестрация LLM-пайплайна
   - Валидация и безопасность (JWT, RBAC)
   - Экспорт в Moodle XML

3. **PostgreSQL 15**
   - Хранение пользователей и ролей
   - Хранение документов, тестов, вопросов
   - История генераций и экспортов
   - Аудит-логи

4. **Redis 7** (опционально)
   - Кэш промежуточных шагов
   - Сессии экспорта

5. **LLM API**
   - Perplexity API - основной провайдер
   - OpenAI API - fallback
   - YandexGPT - российская альтернатива

6. **Nginx**
   - Балансировка запросов
   - SSL/TLS терминация
   - Раздача статики frontend
   - Reverse proxy для backend

7. **Observability Stack**
   - Prometheus - сбор метрик
   - Grafana - визуализация (дашборды)
   - Loki - агрегация логов
   - Promtail - коллектор логов из Docker

#### Взаимодействие компонентов

**Сценарий 1: Загрузка и парсинг документа**

```
User → Vue SPA → POST /api/v1/documents (multipart/form-data)
                      ↓
                 Auth Middleware (JWT проверка)
                      ↓
                 Document Handler
                      ↓
                 Upload Use Case
                      ↓
                 Document Repository → PostgreSQL (INSERT)
                      ↓
                 Parser Factory → выбор парсера (PDF/DOCX/PPTX/TXT/MD)
                      ↓
                 Document Repository → PostgreSQL (UPDATE parsed_text, status)
                      ↓
                 Response ← Vue SPA (обновление UI)
```

**Сценарий 2: Генерация теста с LLM**

```
User → Vue SPA → POST /api/v1/tests/generate
                      ↓
                 Auth Middleware + Role Middleware (teacher/admin)
                      ↓
                 Test Handler
                      ↓
                 Generate Use Case
                      ↓
                 Document Repository → PostgreSQL (SELECT parsed_text)
                      ↓
                 LLM Factory → выбор стратегии (Perplexity/OpenAI/YandexGPT)
                      ↓
                 LLM API (HTTP request с промптом)
                      ↓
                 LLM Response (JSON с вопросами)
                      ↓
                 Test Repository → PostgreSQL (INSERT test, questions, answers)
                      ↓
                 Response ← Vue SPA (отображение сгенерированного теста)
```

**Сценарий 3: Экспорт в Moodle**

```
User → Vue SPA → GET /api/v1/tests/:id/export-xml
                      ↓
                 Auth Middleware
                      ↓
                 Moodle Handler
                      ↓
                 Export Use Case
                      ↓
                 Test Repository → PostgreSQL (SELECT с вопросами и ответами)
                      ↓
                 Moodle XML Exporter (формирование XML)
                      ↓
                 Response ← Vue SPA (скачивание XML файла)
```

**Мониторинг и логирование**

```
Backend HTTP Requests → Request ID Middleware
                              ↓
                         HTTP Logger Middleware
                              ↓
                         Uber Zap → JSON logs → Docker stdout
                              ↓
                         Promtail (scrape logs)
                              ↓
                         Loki (store & index)
                              ↓
                         Grafana (visualization)

Backend Metrics → Prometheus Client (HTTP middleware)
                       ↓
                  /metrics endpoint
                       ↓
                  Prometheus (scrape)
                       ↓
                  Grafana (dashboards)
```

## 4. Технологический стек

### 4.1 Backend (Go 1.23.1)

**Основной фреймворк - Fiber v2**:

- Высокая производительность (построен на fasthttp)
- Express-like API (знакомый синтаксис для Node.js разработчиков)
- Поддержка цепочек middleware
- Поддержка WebSocket
- Cookie и Session management
- Поддержка шаблонизаторов
- Раздача статических файлов

**ORM и работа с БД**:

- GORM - полнофункциональный ORM
  - Автоматические миграции
  - Ассоциации (BelongsTo, HasOne, HasMany, ManyToMany)
  - Хуки (BeforeCreate, AfterUpdate, BeforeDelete и др.)
  - Мягкое удаление (поле deleted_at)
  - Prepared Statements (защита от SQL injection)
  - Поддержка транзакций
  - Scopes для переиспользования запросов
- golang-migrate - версионирование миграций
  - Up/Down миграции
  - Поддержка отката (rollback)
  - Инструмент командной строки

**Dependency Injection - Wire**:

- Внедрение зависимостей на этапе компиляции (нет runtime overhead)
- Типобезопасность
- Генерация кода
- Автоматическая проверка зависимостей
- Наборы провайдеров

**Мониторинг и логирование**:

- Uber Zap logger
  - Структурированное логирование (JSON + Console)
  - Множество целей вывода
  - Уровни логирования (Debug, Info, Warn, Error, Fatal)
  - Контекстные поля (request_id, user_id, action)
  - Высокая производительность
- Prometheus client
  - HTTP метрики (количество запросов, длительность, активные)
  - Пользовательские метрики
  - Histogram для задержек
  - Counter для событий

**Парсинг документов**:

- unidoc/unioffice (DOCX, PPTX)
  - Парсинг полной структуры документа
  - Извлечение текста
  - Чтение метаданных
- ledongthuc/pdf (PDF)
  - Извлечение текста
  - Постраничный парсинг
- goldmark (Markdown)
  - Совместимость с CommonMark
  - Поддержка расширений
- Standard library (TXT)

**Безопасность**:

- golang-jwt/jwt v5
  - Подпись HS256
  - Валидация claims
  - Истечение токенов
- golang.org/x/crypto/bcrypt
  - Адаптивное хеширование
  - Генерация соли
  - Фактор сложности 10

**Тестирование**:

- testify/suite - наборы тестов
- testify/assert - проверки утверждений
- testify/mock - создание моков
- testify/require - критические проверки
- Built-in testing - бенчмарки

### 4.2 Frontend (Vue 3.5)

**Build Tool - Vite 7**:

- Молниеносный HMR (Hot Module Replacement)
- Продакшн сборки на основе Rollup
- Встроенная поддержка TypeScript
- Экосистема плагинов
- Препроцессинг CSS
- Оптимизация ассетов

**Язык - TypeScript 5.6**:

- Типобезопасность
- Интерфейсы и типы
- Перечисления (enums)
- Обобщения (generics)
- Декораторы

**Роутинг - Vue Router 4**:

- Режим истории (History mode)
- Навигационные guards (beforeEach, beforeResolve, afterEach)
- Ленивая загрузка
- Вложенные маршруты
- Именованные маршруты
- Мета-поля маршрутов

**State Management - Pinia**:

- Стиль Composition API
- TypeScript в приоритете
- Интеграция с DevTools
- Горячая замена модулей
- Система плагинов
- Actions, Getters, State

**UI Framework**:

- Tailwind CSS v4
  - Подход utility-first
  - JIT компиляция
  - Темный режим
  - Адаптивный дизайн
  - Пользовательская конфигурация
- Headless UI
  - Полная доступность
  - Компоненты без стилей
  - Dialog, Menu, Listbox и др.
- Heroicons
  - 292+ иконки
  - Версии Solid и Outline
  - Формат SVG

**HTTP - Axios**:

- Перехватчики запросов/ответов
- Автоматическое преобразование JSON
- Отмена запросов
- Отслеживание прогресса
- Поддержка cookies

**Validation**:

- VeeValidate
  - Валидация на уровне полей
  - Валидация на уровне форм
  - Пользовательские правила
- Yup
  - Валидация схем
  - Вывод типов
  - Пользовательские валидаторы

**Testing - Vitest**:

- API совместимое с Jest
- Быстрое выполнение (на основе Vite)
- ESM в приоритете
- Поддержка TypeScript
- Отчеты о покрытии (c8)
- Режим UI
- 112+ тестов

**Логирование**:

- Custom Logger utility
  - Перечисление LogLevel (DEBUG, INFO, WARN, ERROR)
  - Интеграция с консолью браузера
  - Структурированный вывод
  - HTTP логирование через interceptors
  - Логирование действий в store

### 4.3 AI/ML

**LLM Providers**:

1. Perplexity API (основной)
   - Модель: llama-3.1-sonar-large-128k-online
   - Контекстное окно: 128k токенов
   - Интеграция онлайн-поиска
   - Режим JSON
   - Поддержка потоковой передачи

2. OpenAI API (fallback)
   - Модели: GPT-4, GPT-3.5-turbo
   - Вызов функций
   - Возможности компьютерного зрения
   - Поддержка дообучения

3. YandexGPT (российский)
   - Модели: yandexgpt-lite, yandexgpt
   - Аутентификация на основе папок
   - IAM токен
   - Потоковая передача

**Implementation**:

- Strategy Pattern
- Factory Pattern
- Конфигурируемый выбор провайдера
- Обработка ошибок с fallback
- Ограничение частоты запросов
- Управление таймаутами

**Промпт Engineering**:

```
Системный промпт:
"Ты - эксперт по созданию образовательных тестов.
Создай {num_questions} вопросов из текста ниже.
Типы: {question_types}
Сложность: {difficulty}
Формат ответа: JSON"
```

### 4.4 Database

**PostgreSQL 15**:

- Первичные ключи UUID
- JSONB для метаданных
- Полнотекстовый поиск (tsvector)
- Частичные индексы
- Конкурентные индексы
- Ограничения внешних ключей
- Ограничения проверки
- Триггеры для аудита
- Безопасность на уровне строк (опционально)

**Миграции**:

```sql
-- 000001_init_schema.up.sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    role_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
```

**Redis 7** (опционально):

- Кэширование часто запрашиваемых данных
- Хранилище сессий
- Счетчики ограничения частоты запросов
- Pub/Sub для обновлений в реальном времени

### 4.5 Infrastructure

**Docker Compose**:

```yaml
services:
  backend:
    build: ./backend
    ports: ["8080:8080"]
    depends_on: [postgres, redis]

  frontend:
    build: ./frontend
    ports: ["5173:5173"]

  postgres:
    image: postgres:15
    volumes: [postgres_data:/var/lib/postgresql/data]

  redis:
    image: redis:7-alpine

  nginx:
    image: nginx:alpine
    ports: ["80:80", "443:443"]

  prometheus:
    image: prom/prometheus

  grafana:
    image: grafana/grafana

  loki:
    image: grafana/loki

  promtail:
    image: grafana/promtail
```

**Nginx Configuration**:

```nginx
upstream backend {
    server backend:8080;
}

upstream frontend {
    server frontend:5173;
}

server {
    listen 80;
    server_name localhost;

    # Frontend
    location / {
        proxy_pass http://frontend;
    }

    # Backend API
    location /api/ {
        proxy_pass http://backend;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    # Monitoring
    location /metrics {
        proxy_pass http://backend/metrics;
    }
}
```

**Grafana Dashboards**:

1. TestGen API Dashboard
   - Request Rate (req/s)
   - Response Time (p50, p95, p99)
   - HTTP Status Codes
   - Requests by Endpoint
   - Active Connections

2. TestGen Logs Dashboard
   - All Backend Logs
   - Error Logs
   - HTTP 5xx Errors
   - Log Levels Over Time
   - Top 10 Endpoints

## 5. Безопасность

### 5.1 Аутентификация

**JWT Tokens**:

- Algorithm: HS256
- Secret: 256-bit key (из env)
- Expiration: 24 hours
- Claims: user_id, email, role_id, exp, iat
- Token в HTTP-only cookie + Authorization header

**Процесс аутентификации**:

```
1. POST /auth/login {email, password}
2. Backend проверяет bcrypt hash
3. Генерация JWT token
4. Set-Cookie: testgen_token=<jwt>; HttpOnly; Secure; SameSite=Lax
5. Response: {user, token}
```

**Middleware chain**:

```
Request → Recovery → RequestID → CORS → JWT Auth → Role Check → Handler
```

### 5.2 Авторизация (RBAC)

**Роли в БД**:

```sql
CREATE TABLE roles (
    id UUID PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL, -- admin, teacher, student
    description TEXT
);
```

**Permission matrix**:

| Endpoint | Admin | Teacher | Student |
|----------|-------|---------|---------|
| POST /documents | ✅ | ✅ | ❌ |
| GET /documents | ✅ (all) | ✅ (own) | ❌ |
| POST /tests/generate | ✅ | ✅ | ❌ |
| GET /tests | ✅ (all) | ✅ (own) | ✅ (assigned) |
| PUT /users/:id/role | ✅ | ❌ | ❌ |
| GET /users | ✅ | ✅ | ❌ |

**Middleware implementation**:

```go
func AdminOnly() fiber.Handler {
    return func(c *fiber.Ctx) error {
        user := c.Locals("user").(*entity.User)
        if !user.IsAdmin() {
            return fiber.ErrForbidden
        }
        return c.Next()
    }
}
```

### 5.3 Защита от уязвимостей

**SQL Injection**:

- Подготовленные выражения GORM
- Параметризованные запросы
- Без конкатенации сырого SQL

**XSS (Cross-Site Scripting)**:

- HTTP-only cookies (нет доступа из JavaScript)
- Заголовки CSP
- Санитизация HTML на backend
- Санитизация на frontend при отображении

**CSRF (Cross-Site Request Forgery)**:

- Cookies с SameSite=Lax
- Паттерн двойной отправки cookie (опционально)
- Валидация токенов для изменяющих запросов

**CORS**:

```go
app.Use(cors.New(cors.Config{
    AllowOrigins: "http://localhost:3000,http://localhost:5173",
    AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
    AllowHeaders: "Origin,Content-Type,Accept,Authorization",
    AllowCredentials: true, // для cookies
}))
```

**Rate Limiting**:

- Nginx rate limiting (10 req/s per IP)
- Backend rate limiting для критичных эндпоинтов
- Redis счетчики

**Input Validation**:

- go-playground/validator на backend
- VeeValidate + Yup на frontend
- Валидация типов файлов
- Валидация размера файлов (макс 50MB)
- Валидация формата email
- Валидация надежности пароля

### 5.4 Audit Logging

**Activity Logs Table**:

```sql
CREATE TABLE activity_logs (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    action VARCHAR(255) NOT NULL, -- 'login', 'create_test', 'export_xml'
    entity_type VARCHAR(100), -- 'test', 'document', 'user'
    entity_id UUID,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_activity_logs_user_id ON activity_logs(user_id);
CREATE INDEX idx_activity_logs_created_at ON activity_logs(created_at);
```

**Logged actions**:

- Вход/выход пользователя
- Изменение ролей (действие Admin)
- Генерация тестов
- Экспорт тестов
- Загрузка/удаление документов
- Синхронизация с Moodle

## 6. План разработки и управление

### 6.1 Фаза 1: Сбор требований (1-2 недели)

**Задачи**:

1. Интервью с преподавателями (целевая аудитория)
2. Анализ существующих решений (Kahoot, Quizlet, Google Forms)
3. Изучение Moodle API и XML формата
4. Изучение возможностей LLM для генерации вопросов
5. Определение MVP scope
6. Создание User Stories
7. Приоритизация требований (MoSCoW)

**Артефакты**:

- Документ спецификации требований
- Пользовательские истории (Admin, Teacher, Student)
- Use Case диаграммы (UML)
- Функциональные требования (FR-1...FR-7)
- Нефункциональные требования (NFR-1...NFR-5)

**Критерии завершения**:

- Утверждены требования stakeholders
- Определен MVP scope
- Создан backlog

### 6.2 Фаза 2: Проектирование (1 неделя)

**Задачи**:

1. Архитектурное проектирование (Clean Architecture)
2. Проектирование БД (ERD диаграмма)
3. API контракты (Swagger/OpenAPI)
4. Выбор технологий (Go, Vue, PostgreSQL, LLM)
5. UML диаграммы (Sequence, Class, Activity)
6. UI/UX wireframes
7. Дизайн паттернов (Repository, Factory, Strategy)

**Артефакты**:

- Документ проектирования ПО (SDD)
- ERD диаграмма (Crow's Foot Notation)
- Спецификация API (Swagger)
- UML диаграммы
- Каркасы интерфейсов (Figma/Sketch)
- Документ технологического стека
- Записи архитектурных решений (ADR)

**Критерии завершения**:

- Утвержден дизайн архитектуры
- Определена схема БД
- Согласованы API контракты
- Выбраны технологии

### 6.3 Фаза 3: Реализация Backend (3-4 недели)

**Sprint 1 (неделя 1): Infrastructure + Auth**

- Настройка проекта (Go modules, структура директорий)
- Конфигурация Docker Compose
- Настройка PostgreSQL + миграции
- Настройка Wire DI
- Аутентификация JWT
- User CRUD + управление ролями
- Unit тесты для auth

**Sprint 2 (неделя 2): Document Management**

- Обработчик загрузки документов
- Parser Factory (PDF, DOCX, PPTX, TXT, MD)
- Document CRUD операции
- Асинхронный парсинг
- Обработка ошибок
- Unit тесты для парсеров

**Sprint 3 (неделя 3): LLM Integration + Test Generation**

- LLM Strategy Pattern
- LLM Factory
- Интеграция с Perplexity/OpenAI/YandexGPT
- Use case генерации тестов
- Question/Answer CRUD
- Unit тесты для LLM

**Sprint 4 (неделя 4): Moodle Integration + Monitoring**

- Экспортер Moodle XML
- Moodle REST API клиент
- Метрики Prometheus
- Логирование Zap
- Проверки работоспособности
- Интеграционные тесты

**Артефакты**:

- Рабочий backend API
- Документация Swagger
- Unit тесты (покрытие ≥80%)
- Интеграционные тесты
- Docker образы

### 6.4 Фаза 4: Реализация Frontend (3-4 недели)

**Sprint 1 (неделя 1): Setup + Auth**

- Настройка проекта Vite
- Конфигурация TypeScript
- Настройка Tailwind CSS
- Настройка Vue Router
- Хранилища Pinia
- Формы входа/регистрации
- Защитники маршрутов auth
- Unit тесты для auth

**Sprint 2 (неделя 2): Document Management UI**

- Компонент загрузки документов
- Представление списка документов
- Детальное представление документа
- Валидация файлов
- Индикаторы прогресса
- Unit тесты

**Sprint 3 (неделя 3): Test Management UI**

- Форма генерации тестов
- Представление списка тестов
- Детальное представление теста
- Редактор вопросов
- Редактор ответов
- Unit тесты

**Sprint 4 (неделя 4): Admin + Polish**

- Управление пользователями (Admin)
- Панель со статистикой
- UI экспорта Moodle
- Обработка ошибок
- Состояния загрузки
- Адаптивный дизайн
- Unit тесты

**Артефакты**:

- Рабочее frontend приложение
- Адаптивный UI
- Unit тесты (112+ тестов)
- Библиотека компонентов

### 6.5 Фаза 5: Тестирование + Hardening (2 недели)

**Week 1: Testing**

- Интеграционное тестирование (E2E)
- Тестирование производительности (нагрузочные тесты с k6)
- Тестирование безопасности (SAST, сканирование зависимостей)
- Кроссбраузерное тестирование
- Тестирование доступности (WCAG 2.1)
- Юзабилити-тестирование
- Исправление багов

**Week 2: Hardening**

- Усиление безопасности (тестирование на проникновение)
- Оптимизация производительности
- Улучшения обработки ошибок
- Улучшения логирования
- Настройка мониторинга (Prometheus + Grafana)
- Проверка документации
- Ревью кода

**Артефакты**:

- Отчет по тестированию
- Отчет по производительности
- Отчет по аудиту безопасности
- Исправления багов
- Оптимизации

### 6.6 Фаза 6: Документация + Релиз (1 неделя)

**Документация**:

- Руководство пользователя
- Документация API (Swagger)
- Руководство по развертыванию
- Документация архитектуры
- Руководства по эксплуатации
- Обучающие материалы

**Релиз**:

- Развертывание в продакшн
- Настройка мониторинга
- Настройка резервного копирования
- Обучающие сессии для пользователей
- Примечания к релизу
- План поддержки после релиза

**Артефакты**:

- Полная документация
- Продакшн развертывание
- Обучающие материалы
- Примечания к релизу

## 7. Макеты и прототипы

### 7.1 User Flow

**Преподаватель - Создание теста**:

```
1. Логин → Dashboard
2. Переход к документам
3. Загрузка документа (PDF/DOCX/PPTX/TXT/MD)
4. Ожидание парсинга (прогресс-бар)
5. Просмотр распарсенного текста
6. Клик "Генерировать тест"
7. Настройка параметров:
   - Количество вопросов: 20
   - Типы вопросов: Single Choice
   - Сложность: Medium
8. Отправка на генерацию
9. Ожидание LLM (прогресс-бар)
10. Просмотр сгенерированного теста
11. Редактирование вопросов при необходимости
12. Экспорт в Moodle XML
13. Скачивание XML файла
14. Импорт в Moodle
```

**Администратор - Управление пользователями**:

```
1. Логин → Dashboard
2. Переход к пользователям
3. Просмотр списка пользователей (с пагинацией)
4. Выбор пользователя
5. Изменение роли (teacher/student/admin)
6. Подтверждение изменения
7. Просмотр журнала активности
```

### 7.2 UI Components

**Dashboard View**:

- Заголовок (логотип, меню пользователя, выход)
- Боковая панель (навигация)
- Карточки статистики (количество документов, тестов, вопросов)
- Недавняя активность
- Быстрые действия (Загрузить документ, Создать тест)

**Document Upload**:

- Область перетаскивания
- Кнопка выбора файла
- Бейдж поддерживаемых форматов
- Индикатор максимального размера (50MB)
- Прогресс-бар
- Сообщения об успехе/ошибке

**Test Generation Form**:

- Селектор документа (выпадающий список)
- Количество вопросов (слайдер 1-50)
- Типы вопросов (чекбоксы)
- Сложность (радио-кнопки: Easy/Medium/Hard)
- Кнопка генерации
- Кнопка отмены

**Question Editor**:

- Текст вопроса (текстовая область)
- Селектор типа вопроса
- Селектор сложности
- Баллы (числовой ввод)
- Варианты ответов (динамический список)
- Индикатор правильного ответа (чекбокс/радио)
- Кнопка добавления ответа
- Кнопка удаления ответа
- Кнопка сохранения
- Кнопка отмены

### 7.3 Responsive Design

**Breakpoints** (Tailwind):

- sm: 640px (mobile)
- md: 768px (tablet)
- lg: 1024px (desktop)
- xl: 1280px (large desktop)

**Mobile adaptations**:

- Сворачиваемая боковая панель → гамбургер-меню
- Карточки в стопку вместо сетки
- Кнопки для сенсорного управления (мин 44x44px)
- Упрощенные формы
- Нижняя панель навигации

### 7.4 Accessibility (ГОСТ Р ИСО 9241)

**Клавиатурная навигация**:

- Порядок Tab
- Индикаторы фокуса
- Клавиатурные сокращения (Ctrl+N - новый документ, Ctrl+T - новый тест)
- Ссылка пропуска к контенту
- Захват фокуса в модальных окнах

**Контрастность**:

- Соответствие WCAG AA
- 4.5:1 для обычного текста
- 3:1 для крупного текста
- Режим высокой контрастности

**Скринридеры**:

- Семантический HTML (header, nav, main, footer)
- ARIA метки
- ARIA live регионы для динамического контента
- Alt текст для изображений
- Role атрибуты

**Масштабирование**:

- Поддержка масштабирования до 200%
- Rem единицы для шрифтов
- Адаптивная разметка
- Без горизонтальной прокрутки

## 8. Стратегия тестирования

### 8.1 Backend Testing (Go)

**Unit Tests**:

```go
// Пример: тест сущности User
func TestUser_SetPassword(t *testing.T) {
    user := &entity.User{}
    password := "TestPass123!"

    err := user.SetPassword(password)
    assert.NoError(t, err)
    assert.NotEmpty(t, user.PasswordHash)
    assert.True(t, user.CheckPassword(password))
    assert.False(t, user.CheckPassword("WrongPassword"))
}

// Пример: тест репозитория
func TestUserRepository_Create(t *testing.T) {
    db := setupTestDB(t)
    repo := postgres.NewUserRepository(db)

    user := &entity.User{
        Email: "test@example.com",
        FullName: "Test User",
    }
    user.SetPassword("password123")

    err := repo.Create(context.Background(), user)
    assert.NoError(t, err)
    assert.NotEqual(t, uuid.Nil, user.ID)
}
```

**Integration Tests**:

```go
// Пример: тест HTTP обработчика
func TestAuthHandler_Login(t *testing.T) {
    app, cleanup := setupTestApp(t)
    defer cleanup()

    // Создание тестового пользователя
    createUser(t, app, "test@example.com", "password123")

    // Запрос на логин
    req := httptest.NewRequest("POST", "/auth/login", strings.NewReader(`{
        "email": "test@example.com",
        "password": "password123"
    }`))
    req.Header.Set("Content-Type", "application/json")

    resp, err := app.Test(req)
    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)

    // Проверка cookie
    cookies := resp.Cookies()
    assert.NotEmpty(t, cookies)
    assert.Equal(t, "testgen_token", cookies[0].Name)
}
```

**Test Coverage**:

- Domain entities: 90%
- Use cases: 85%
- Handlers: 80%
- Repositories: 90%
- Overall: ≥80%

**Команды**:

```bash
# Запуск всех тестов
go test ./... -v

# С покрытием
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Конкретный пакет
go test ./internal/domain/entity -v

# Бенчмарки
go test -bench=. -benchmem
```

### 8.2 Frontend Testing (Vitest)

**Unit Tests**:

```typescript
// Пример: тест store аутентификации
describe('authStore', () => {
  it('should login successfully', async () => {
    const store = useAuthStore()
    vi.mocked(authService.login).mockResolvedValue({
      user: mockUser,
      token: 'fake-token'
    })

    await store.login({ email: 'test@example.com', password: 'pass123' })

    expect(store.isAuthenticated).toBe(true)
    expect(store.user).toEqual(mockUser)
  })

  it('should handle login error', async () => {
    const store = useAuthStore()
    vi.mocked(authService.login).mockRejectedValue(new Error('Invalid credentials'))

    await expect(store.login({ email: 'wrong@example.com', password: 'wrong' }))
      .rejects.toThrow('Invalid credentials')

    expect(store.isAuthenticated).toBe(false)
    expect(store.error).toBe('Invalid credentials')
  })
})
```

**Component Tests**:

```typescript
// Пример: тест компонента LoginForm
describe('LoginForm', () => {
  it('should render form fields', () => {
    const wrapper = mount(LoginForm)

    expect(wrapper.find('input[type="email"]').exists()).toBe(true)
    expect(wrapper.find('input[type="password"]').exists()).toBe(true)
    expect(wrapper.find('button[type="submit"]').exists()).toBe(true)
  })

  it('should submit form on valid input', async () => {
    const wrapper = mount(LoginForm)

    await wrapper.find('input[type="email"]').setValue('test@example.com')
    await wrapper.find('input[type="password"]').setValue('password123')
    await wrapper.find('form').trigger('submit')

    expect(authService.login).toHaveBeenCalledWith({
      email: 'test@example.com',
      password: 'password123'
    })
  })
})
```

**Test Statistics**:

- Total tests: 112+
- Auth: 25 tests
- Documents: 21 tests
- Tests: 15 tests
- Users: 15 tests
- Utils: 36 tests
- Services: 20+ tests

**Команды**:

```bash
# Запуск всех тестов
npm run test

# С покрытием
npm run test -- --coverage

# Режим наблюдения
npm run test -- --watch

# Режим UI
npm run test:ui
```

### 8.3 API Testing

**Swagger UI**:

- Интерактивное тестирование API
- Функция пробного запроса
- Примеры запросов/ответов
- Валидация схем

### 8.4 Performance Testing

**k6 Load Tests**:

```javascript
import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  stages: [
    { duration: '1m', target: 10 }, // постепенное увеличение
    { duration: '3m', target: 50 }, // удержание на 50 пользователях
    { duration: '1m', target: 0 },  // постепенное уменьшение
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% < 500ms
    http_req_failed: ['rate<0.01'],   // <1% ошибок
  },
};

export default function () {
  let res = http.get('http://localhost:8080/health');
  check(res, {
    'status is 200': (r) => r.status === 200,
    'response time < 200ms': (r) => r.timings.duration < 200,
  });
  sleep(1);
}
```

**Метрики**:

- Время ответа: p50, p95, p99
- Пропускная способность: запросов/секунду
- Процент ошибок: %
- Одновременные пользователи: максимум
- Подключения к БД: активные

### 8.5 Security Testing

**SAST (Static Application Security Testing)**:

- gosec для Go кода
- npm audit для frontend dependencies
- Snyk для vulnerability scanning

**DAST (Dynamic Application Security Testing)**:

- OWASP ZAP
- Burp Suite
- Тесты на SQL injection
- Тесты на XSS
- Тесты на CSRF
- Тесты обхода аутентификации

**Penetration Testing**:

- Ручное тестирование
- Автоматизированные сканеры
- Оценка уязвимостей
- Попытки эксплуатации
- Генерация отчетов

### 8.6 Критерии успеха

**Functional**:

- All API endpoints work correctly
- File upload supports all formats (PDF, DOCX, PPTX, TXT, MD)
- Parsing accuracy >95%
- LLM generation success rate >90%
- Moodle XML validates correctly
- Authentication/Authorization works

**Performance**:

- API response time <500ms (p95)
- Document parsing <10s (files up to 50MB)
- Test generation <30s (20 questions)
- Support 50+ concurrent users
- Database query time <100ms (p95)

**Security**:

- No critical vulnerabilities (CVSS ≥7.0)
- All inputs validated
- SQL injection protected
- XSS protected
- CSRF protected
- Audit logs working

**Quality**:

- Code coverage ≥80%
- All tests passing
- No linter errors
- TypeScript strict mode
- API documentation complete

## 9. Документация

### 9.1 Technical Documentation

**README.md**:

- Project overview
- Quick start guide
- Prerequisites
- Installation steps
- Configuration
- Running the app
- Testing
- Deployment
- Contributing guidelines

**CLAUDE.md**:

- Architecture overview
- Technology stack
- Project structure
- Design patterns
- Coding conventions
- Development workflow

**API Documentation (Swagger)**:

- Interactive UI: <http://localhost:8080/swagger/index.html>
- JSON spec: `/docs/swagger.json`
- YAML spec: `/docs/swagger.yaml`
- All endpoints documented
- Request/Response examples
- Error codes
- Authentication

**Database Schema (ERD)**:

- Crow's Foot Notation
- Tables: users, roles, documents, tests, questions, answers, activity_logs
- Relationships
- Constraints
- Indexes

**Architecture Decision Records (ADR)**:

- ADR-001: Clean Architecture choice
- ADR-002: Go for backend
- ADR-003: Vue 3 for frontend
- ADR-004: PostgreSQL over NoSQL
- ADR-005: JWT authentication
- ADR-006: HTTP-only cookies
- ADR-007: Strategy Pattern for LLM

### 9.2 User Documentation

**User Manual**:

1. Getting Started
   - Registration
   - Login
   - Dashboard overview

2. Document Management
   - Uploading documents
   - Supported formats
   - Viewing parsed text
   - Deleting documents

3. Test Generation
   - Creating tests manually
   - Generating tests with AI
   - Configuring generation parameters
   - Editing questions
   - Managing answers

4. Moodle Integration
   - Exporting to XML
   - Importing in Moodle
   - Syncing tests
   - Troubleshooting

5. Administration (Admin only)
   - User management
   - Role assignment
   - Activity logs
   - System monitoring

**FAQ**:

- How to upload large documents?
- Which LLM provider is used?
- How to change user role?
- How to export test to Moodle?
- How to view logs?

### 9.3 Deployment Guide

**Docker Deployment**:

```bash
# 1. Клонирование репозитория
git clone <repo-url>
cd testgen

# 2. Настройка окружения
cp .env.example .env
# Отредактируйте .env своими значениями

# 3. Запуск сервисов
docker-compose up -d

# 4. Проверка здоровья
curl http://localhost:8080/health

# 5. Доступ к приложению
# Frontend: http://localhost:5173
# Backend API: http://localhost:8080
# Swagger: http://localhost:8080/swagger/index.html
# Grafana: http://localhost:3000
```

**Manual Deployment**:

```bash
# Backend
cd backend
go mod download
go build -o bin/api cmd/api/main.go
./bin/api

# Frontend
cd frontend
npm install
npm run build
npm run preview

# PostgreSQL
psql -U postgres
CREATE DATABASE testgen_db;
```

**Production Checklist**:

- [ ] Change JWT_SECRET
- [ ] Change ADMIN_PASSWORD
- [ ] Set COOKIE_SECURE=true
- [ ] Set LOG_FORMAT=json
- [ ] Configure CORS for production domain
- [ ] Set up SSL/TLS certificates
- [ ] Configure backup strategy
- [ ] Set up monitoring alerts
- [ ] Configure log rotation
- [ ] Review security settings

### 9.4 Operations (Runbooks)

**Health Checks**:

```bash
# Проверка здоровья Backend
curl http://localhost:8080/health

# Проверка подключения к БД
docker-compose exec postgres psql -U testgen_user -d testgen_db -c "SELECT 1"

# Метрики
curl http://localhost:8080/metrics
```

**Backup & Restore**:

```bash
# Резервное копирование
docker-compose exec postgres pg_dump -U testgen_user testgen_db > backup.sql

# Восстановление
docker-compose exec postgres psql -U testgen_user testgen_db < backup.sql
```

**Logging**:

```bash
# Логи Backend
docker-compose logs -f backend

# Просмотр в Grafana
# Перейдите к источнику данных Loki
# Запрос: {service="backend"} |= "error"
```

**Monitoring**:

```bash
# Цели Prometheus
curl http://localhost:9090/targets

# Дашборды Grafana
# http://localhost:3000/d/testgen-api
# http://localhost:3000/d/testgen-logs
```

**Troubleshooting**:

```bash
# БД не подключается
docker-compose ps
docker-compose logs postgres

# Backend не запускается
docker-compose logs backend
# Проверьте конфигурацию .env

# Frontend не загружается
docker-compose logs frontend
npm run build # пересборка
```

## 10. Итоги и перспективы

### 10.1 Полученные знания

**Технические знания**:

1. **Clean Architecture на практике**
   - Разделение на 4 слоя (Domain, Application, Infrastructure, Interfaces)
   - Принцип инверсии зависимостей
   - Тестируемая архитектура
   - Repository Pattern
   - Use Cases (варианты использования)

2. **Full-stack разработка**
   - Backend: Go + Fiber + GORM
   - Frontend: Vue 3 + TypeScript + Vite
   - Управление состоянием: Pinia
   - Проектирование API: принципы RESTful
   - Обновления в реальном времени

3. **Database Design**
   - Проектирование схемы PostgreSQL
   - Оптимизация индексов
   - Стратегия миграций
   - Мягкое удаление
   - Журналирование аудита

4. **AI/ML Integration**
   - Интеграция LLM API (Perplexity, OpenAI, YandexGPT)
   - Промпт-инжиниринг
   - Strategy Pattern для гибкости
   - Обработка ошибок для AI вызовов
   - Механизмы отката

5. **DevOps & Infrastructure**
   - Контейнеризация Docker
   - Оркестрация Docker Compose
   - Конфигурация Nginx
   - Мониторинг Prometheus
   - Дашборды Grafana
   - Агрегация логов Loki

6. **Security**
   - Аутентификация JWT
   - Авторизация RBAC
   - Хеширование паролей bcrypt
   - Защита от XSS/CSRF/SQL injection
   - Журналирование аудита
   - Безопасная обработка cookies

7. **Testing**
   - Unit тестирование (Go + Vue)
   - Интеграционное тестирование
   - Создание моков зависимостей
   - Анализ покрытия тестами
   - Тестирование производительности (k6)

**Предметная область**:

1. **Образовательные технологии**
   - Типы тестовых вопросов
   - Уровни сложности
   - Педагогический дизайн
   - Качество вопросов

2. **Moodle Integration**
   - Формат Moodle XML
   - Категории вопросов
   - Web Services API
   - Банк вопросов

3. **Автоматизация в образовании**
   - Генерация контента с AI
   - Валидация качества
   - Оптимизация рабочих процессов

### 10.2 Достигнутые результаты

**Функциональность**:

- Полноценная система генерации тестов
- 5 форматов документов (PDF, DOCX, PPTX, TXT, MD)
- 3 LLM провайдера
- Интеграция с Moodle
- RBAC с 3 ролями
- Мониторинг готовый к продакшн
- Всестороннее тестирование

**Метрики**:

- ~100 Go файлов
- ~64 Vue/TS файлов
- 112+ frontend тестов
- 20+ API эндпоинтов
- 7 таблиц базы данных
- 6+ паттернов проектирования
- Полная документация Swagger

**Performance**:

- Ответ API <500ms (p95)
- Парсинг <10s (файлы до 50MB)
- Генерация <30s (20 вопросов)
- Поддержка 50+ одновременных пользователей

### 10.3 Возможные улучшения

**MVP+ (3-6 месяцев)**:

1. **Расширение функциональности**
   - Вопросы на соответствие
   - Эссе вопросы
   - Вопросы с пропусками (заполнение пропусков)
   - Вычисляемые вопросы
   - Редактор форматированного текста (TinyMCE/Quill)
   - Загрузка изображений для вопросов
   - Поддержка формул (MathJax/KaTeX)

2. **Улучшение UX**
   - Перетаскивание для упорядочивания вопросов
   - Массовые операции (удаление, экспорт)
   - Шаблоны вопросов
   - Избранное/закладки
   - Поиск и фильтрация
   - Режим предпросмотра
   - Совместное редактирование

3. **Аналитика**
   - Статистика на дашборде
   - Аналитика использования
   - Метрики качества вопросов
   - Отчеты о активности пользователей
   - Экспорт отчетов (PDF, Excel)

4. **Интеграция**
   - Google Classroom
   - Microsoft Teams
   - Canvas LMS
   - Blackboard
   - Экспорт SCORM

**Среднесрочные (6-12 месяцев)**:

1. **AI Improvements**
   - Дообученные модели для образования
   - Оптимизация промптов (A/B тестирование)
   - Предсказание сложности вопросов
   - Автоматическая пометка по темам
   - Обнаружение плагиата
   - Оценка качества

2. **Масштабирование**
   - Кеширование Redis
   - CDN для статических файлов
   - Реплики базы данных для чтения
   - Горизонтальное масштабирование backend
   - Балансировка нагрузки
   - Улучшения ограничения скорости

3. **Advanced Features**
   - Версионирование вопросов (подобно Git)
   - Обратная связь по вопросам от студентов
   - Адаптивное тестирование
   - Шаблоны тестов
   - Пулы вопросов
   - Правила рандомизации

4. **Collaboration**
   - Командные рабочие пространства
   - Общие банки вопросов
   - Комментарии к вопросам
   - Рабочий процесс проверки
   - Процесс утверждения

**Долгосрочные (1-2 года, Магистерская диссертация)**:

1. **Микросервисы**
   - Сервис парсинга документов
   - Сервис генерации LLM
   - Сервис экспорта
   - Сервис аналитики
   - Событийно-ориентированная архитектура
   - Очередь сообщений (Kafka/RabbitMQ)
   - Service mesh (Istio)

2. **Machine Learning**
   - Кастомные модели (Hugging Face)
   - Локальное развертывание LLM (Ollama)
   - Трансферное обучение
   - Конвейер дообучения моделей
   - Активное обучение из обратной связи

3. **Multi-tenancy SaaS**
   - Управление организациями
   - Изоляция ресурсов
   - Биллинг и подписки
   - Квоты использования
   - Белая метка
   - Кастомные домены

4. **Advanced Analytics**
   - ML для предсказания качества
   - Обнаружение аномалий
   - Рекомендательная система
   - Предиктивная аналитика
   - Дашборды в реальном времени

5. **Mobile & Offline**
   - Приложение React Native
   - Архитектура offline-first
   - Стратегия синхронизации
   - Push уведомления
   - UI оптимизированный для мобильных

6. **Internationalization**
   - Поддержка нескольких языков
   - Языки с письмом справа налево
   - Локализация
   - Поддержка валют
   - Обработка часовых поясов

### 10.4 Бизнес перспективы

**Целевой рынок**:

- Университеты (500+ в России)
- Онлайн-платформы (аналоги Coursera, Udemy)
- Корпоративное обучение
- Школы (частные, государственные)

**Монетизация**:

- Freemium модель (5 тестов/месяц бесплатно)
- Pro тариф ($29/месяц) - неограниченные тесты
- Enterprise тариф ($199/месяц) - командные функции
- Ценообразование по использованию для LLM затрат

**ROI для клиентов**:

- Экономия 3-5 часов на тест
- Стоимость часа преподавателя: $20-50
- ROI: $60-250 на тест
- Окупаемость: 1-2 месяца

### 10.5 Выводы

**Test Generation System** - это успешная реализация информационной системы, которая:

1. **Решает реальную проблему**
   - Автоматизация создания тестов
   - Экономия времени преподавателей
   - Повышение качества вопросов
   - Интеграция с существующими LMS

2. **Демонстрирует профессиональный уровень разработки**
   - Clean Architecture
   - Код готовый к продакшн
   - Всестороннее тестирование
   - Лучшие практики безопасности
   - Наблюдаемость
   - Документация

3. **Использует современные технологии**
   - Go для backend
   - Vue 3 для frontend
   - PostgreSQL для данных
   - AI/ML для генерации
   - Docker для развертывания
   - Prometheus + Grafana для мониторинга

4. **Готова к масштабированию**
   - Дизайн без состояния
   - Архитектура на основе фич
   - Модульная структура
   - Расширяемые паттерны

5. **Имеет коммерческий потенциал**
   - Четкое ценностное предложение
   - Целевой рынок (образование)
   - Стратегия монетизации
   - Потенциал роста
