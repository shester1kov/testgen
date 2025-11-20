# Test Generation System

Система автоматической генерации тестовых заданий на основе документов с интеграцией в Moodle.

## Описание проекта

Информационная система для автоматического создания тестовых вопросов из учебных материалов (PDF, DOCX, PPTX, TXT) с использованием LLM моделей и интеграцией с платформой Moodle.

**Предметная область**: Образовательные технологии
**Методология**: Waterfall (каскадная модель)
**Архитектура**: Распределенный монолит (Frontend + Backend REST API)

## Архитектура

```
┌─────────────┐
│   Browser   │
└──────┬──────┘
       │ HTTPS
┌──────▼──────────────────┐
│   Nginx (Load Balancer) │
└──────┬──────────────────┘
   ┌───┴────┐
┌──▼───┐ ┌─▼────────────┐
│ Vue  │ │ Go Backend   │
│ SPA  │ │ REST API     │
└──────┘ └──┬───────────┘
            │
    ┌───────┼────────┐
┌───▼────┐ ┌─▼───┐ ┌──▼─────┐
│Postgres│ │Redis│ │LLM API │
└────────┘ └─────┘ └────────┘
```

## Технологии

### Backend
- Go 1.23.1
- Fiber v2 (Web framework)
- GORM (PostgreSQL ORM)
- Wire (Dependency Injection)
- golang-migrate (DB migrations)
- Prometheus (Metrics)

### Frontend
- Vue 3 (Composition API)
- TypeScript
- Vite
- Pinia (State management)
- Vue Router 4
- Tailwind 4
- Axios (HTTP client)

### Infrastructure
- Docker & Docker Compose
- PostgreSQL 15
- Redis 7
- Nginx (Load balancer)
- Prometheus + Grafana (Monitoring)

### AI/ML
- Perplexity API / OpenAI API / YandexGPT

## Установка и запуск

### Требования
- Docker & Docker Compose
- Go 1.23+ (для локальной разработки)
- Node.js 22+ (для локальной разработки)

### Быстрый старт с Docker

1. **Клонировать репозиторий**
```bash
git clone <repository-url>
cd course-work
```

2. **Настроить переменные окружения**
```bash
cp .env.example .env
# Отредактировать .env файл (настроить DB_PASSWORD, JWT_SECRET, LLM API keys)
```

3. **Запустить все сервисы**
```bash
docker-compose up -d
```

**Миграции БД применяются автоматически при запуске backend!** ✅

4. **Открыть приложение**
- Frontend: http://localhost:5173
- Backend API: http://localhost:8080
- Swagger UI: http://localhost:8080/swagger/index.html
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3001 (admin/admin)

### Быстрый старт для Windows (без Docker)

Используйте готовый батник:

```bash
start-backend.bat
```

Батник автоматически:
- Проверит Go и Docker
- Установит зависимости
- Запустит PostgreSQL в Docker
- **Применит миграции автоматически** при старте backend
- Сгенерирует Swagger документацию
- Запустит backend сервер

### Локальная разработка

#### Backend

```bash
cd backend

# Установить зависимости
go mod download

# Сгенерировать Wire DI
go install github.com/google/wire/cmd/wire@latest
wire

# Запустить сервер
go run cmd/api/main.go

# Запустить тесты
go test ./... -v -cover
```

#### Frontend

```bash
cd frontend

# Установить зависимости
npm install

# Запустить dev сервер
npm run dev

# Запустить тесты
npm run test
```

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Регистрация
- `POST /api/v1/auth/login` - Вход
- `POST /api/v1/auth/logout` - Выход
- `GET /api/v1/auth/me` - Текущий пользователь

### Documents
- `POST /api/v1/documents` - Загрузка документа
- `GET /api/v1/documents` - Список документов
- `GET /api/v1/documents/:id` - Детали документа
- `DELETE /api/v1/documents/:id` - Удаление документа

### Tests
- `POST /api/v1/tests` - Создание теста
- `GET /api/v1/tests` - Список тестов
- `GET /api/v1/tests/:id` - Детали теста
- `PUT /api/v1/tests/:id` - Обновление теста
- `DELETE /api/v1/tests/:id` - Удаление теста
- `POST /api/v1/tests/:id/generate` - Генерация вопросов (LLM)
- `POST /api/v1/tests/:id/export` - Экспорт в Moodle XML

Полная документация: см. [CLAUDE.md](CLAUDE.md)

## Тестирование

```bash
# Backend tests
cd backend
go test ./... -v -cover
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Frontend tests
cd frontend
npm run test
npm run test:ui
```

## Мониторинг

- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000 (admin/admin)
- **Метрики**: http://localhost:8080/metrics
- **Health Check**: http://localhost:8080/health

## Безопасность

- JWT аутентификация
- RBAC (Role-Based Access Control)
- Bcrypt hashing для паролей
- HTTPS/TLS шифрование
- Валидация входных данных
- SQL injection защита (GORM)
- XSS защита

## Структура проекта

```
course-work/
├── backend/              # Go backend
│   ├── cmd/             # Entry points
│   ├── internal/        # Private code
│   │   ├── domain/      # Domain layer
│   │   ├── application/ # Use cases
│   │   ├── infrastructure/ # External dependencies
│   │   └── interfaces/  # HTTP handlers
│   └── pkg/             # Shared packages
├── frontend/            # Vue frontend
│   ├── src/
│   │   ├── components/  # Shared components
│   │   ├── features/    # Feature modules
│   │   ├── layouts/     # Page layouts
│   │   ├── router/      # Vue Router
│   │   ├── services/    # API services
│   │   └── stores/      # Pinia stores
├── nginx/               # Nginx configuration
├── prometheus/          # Prometheus config
├── docs/                # Documentation
├── docker-compose.yml   # Docker orchestration
├── CLAUDE.md            # Architecture guide
└── README.md            # This file
```

## Соответствие требованиям

- ✅ Waterfall методология (5 фаз)
- ✅ Распределенный монолит (Frontend + Backend)
- ✅ 5 паттернов проектирования
- ✅ UML диаграммы (Use Case, Class, Activity, Sequence)
- ✅ ERD диаграмма (Crow's Foot Notation)
- ✅ ГОСТ Р ИСО 9241 (эргономика, UI/UX)
- ✅ ГОСТ Р ИСО/МЭК 27001 (безопасность)
- ✅ Тестирование (Backend + Frontend)

## Роли пользователей

- **Admin**: Управление пользователями, просмотр всех данных
- **Teacher**: Загрузка документов, генерация и редактирование тестов
- **Student**: Прохождение тестов (для будущего расширения)
