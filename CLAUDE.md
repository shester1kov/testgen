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
- **Testing**: Vitest, Vue Test Utils
- **Logging**: Custom logger utility with DEBUG/INFO/WARN/ERROR levels
- **Code style**: ESLint + Prettier

### ML/AI

- **LLM API**:
  - YandexGPT
- **Local models** (опционально):
  - Ollama (llama2, mistral)
  - Hugging Face Transformers (для экспериментов)

### Database

- **Primary DB**: PostgreSQL 15
- **Caching**: Redis

### Infrastructure

- **Containerization**: Docker + Docker Compose
- **Load Balancer**: Nginx + SSL (prod)
- **Monitoring**: Prometheus + Grafana
- **Logging**: Structured logging (zerolog/zap) + Loki + Promtail

### Storage

- **S3**: MinIO

### Integration

- **Export formats**: JSON, CSV, Moodle XML

## Structured Logging

### Backend Logging (Go)

#### Features

1. **Structured Logs**: JSON и console форматы
2. **Log Levels**: debug, info, warn, error
3. **Request Tracking**: Автоматическое добавление request ID к каждому запросу
4. **Context Fields**: Поддержка добавления полей для контекста
5. **HTTP Middleware**: Автоматическое логирование всех HTTP запросов

## Core Features (MVP)

### 1. Аутентификация и авторизация

### 2. Управление документами

### 3. Генерация тестов (LLM)

### 4. Управление тестами

### 5. Интеграция с Moodle

- Экспорт тестов в Moodle XML формат

### 6. Мониторинг и логирование

## Design Patterns

### 1. Repository Pattern

**Где**: `internal/domain/repository/` + `internal/infrastructure/persistence/`

**Назначение**: Абстракция доступа к данным, разделение бизнес-логики и слоя данных

### 2. Factory Pattern

**Где**: `internal/infrastructure/parser/`

**Назначение**: Создание парсеров документов в зависимости от типа файла

### 3. Strategy Pattern

**Где**: `internal/infrastructure/llm/`

**Назначение**: Выбор LLM провайдера (OpenAI, YandexGPT)

### 4. Adapter Pattern (Wrapper Pattern)

**Где**: `pkg/logger/logger.go`

**Назначение**: Обертка над zap.Logger

## API Documentation

### Swagger/OpenAPI

- **URL**: `http://localhost:8080/swagger/index.html`

## User Roles

### Admin

- **Полный доступ к системе:**
  - Управление пользователями
  - Загрузка и управление документами
  - Создание, редактирование и удаление тестов
  - Экспорт в Moodle
  - Просмотр всех документов и тестов
  - Доступ к логам и метрикам

### Teacher

- **Создание контента:**
  - Загрузка документов
  - Генерация и редактирование тестов
  - Экспорт в Moodle
  - Просмотр своих тестов и документов
- **Управление пользователями (ограниченное):**
  - Просмотр списка пользователей
  - Назначение тестов студентам (в будущем)
  - **НЕ МОЖЕТ** изменять роли пользователей

### Student

- **Только чтение и прохождение:**
  - Просмотр назначенных тестов
  - **НЕ ИМЕЕТ** доступа к:
    - Загрузке документов
    - Созданию тестов
    - Редактированию тестов
    - Списку пользователей

## Security Requirements (ГОСТ Р ИСО/МЭК 27001-2012)

1. **Аутентификация**: JWT tokens с истечением срока действия
2. **Авторизация**: Role-based access control (RBAC)
3. **Шифрование**: HTTPS/TLS для всех соединений
4. **Хранение паролей**: bcrypt hashing
5. **Валидация входных данных**: защита от SQL injection, XSS
6. **Журналирование**: аудит всех действий пользователей
7. **Ограничение файлов**: проверка типов и размеров файлов

## Testing Strategy

### Backend Tests

- Unit tests: Domain entities, services
- Repository tests: Database operations (with testcontainers)
- Handler tests: HTTP endpoints (mocked dependencies)
- Integration tests: Full API workflow

### Frontend Tests

**Test Structure:** All tests follow co-location pattern in `__tests__/` directories with `.spec.ts` naming convention.

#### Key Testing Patterns

- **Mock-based testing** with Vitest `vi.mock()`
- **User registration without role** - role assigned by backend (default: student)
- **AuthResponse includes token field** - required for JWT authentication
- **Comprehensive negative testing** - negative test cases for security and edge cases
- **localStorage handling** - tests account for undefined/null differences in test environment

#### Coverage Target

- Backend: >80%
- Frontend: >80%

## Important Notes

- Весь код должен быть на английском (комментарии, переменные, функции)
- Коммиты в Git на английском
- Документация может быть на русском (для курсовой)
- Следовать SOLID принципам
- Избегать OWASP Top 10 уязвимостей
- Код должен проходить линтеры (golangci-lint, ESLint)
- Обязательное покрытие тестами критичных функций
- Никогда не хардкодь пароли, токены, URL
- Не хвали код

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
