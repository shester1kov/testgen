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
- **Swagger UI**: http://localhost:8080/swagger/index.html
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000 (admin/admin123)

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

# Сгенерировать Swagger документацию
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/api/main.go -o docs

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

## API Documentation

### Swagger UI

Полная интерактивная документация API доступна через Swagger UI:

- **URL**: <http://localhost:8080/swagger/index.html>
- **Формат**: OpenAPI 2.0 (Swagger)
- **Файлы**: `backend/docs/swagger.json`, `backend/docs/swagger.yaml`

### Основные эндпоинты

#### Authentication (`/auth`)

- `POST /auth/register` - Регистрация (автоматически назначается роль student)
- `POST /auth/login` - Вход (JWT токен в HTTP-only cookie)
- `POST /auth/logout` - Выход (очистка cookie)
- `GET /auth/me` - Текущий пользователь

#### User Management (`/users`) - Admin only

- `GET /users` - Список пользователей с пагинацией
- `PUT /users/{id}/role` - Изменение роли пользователя

#### Documents (`/documents`)

- `POST /documents` - Загрузка документа (PDF, DOCX, PPTX, TXT)
- `GET /documents` - Список документов с пагинацией
- `GET /documents/{id}` - Детали документа
- `POST /documents/{id}/parse` - Парсинг текста из документа
- `DELETE /documents/{id}` - Удаление документа

#### Tests (`/tests`)

- `POST /tests` - Создание теста
- `POST /tests/generate` - Генерация вопросов с помощью LLM
- `GET /tests` - Список тестов с пагинацией
- `GET /tests/{id}` - Детали теста
- `DELETE /tests/{id}` - Удаление теста

#### Moodle Integration (`/moodle`)

- `GET /tests/{id}/export-xml` - Экспорт теста в Moodle XML формат
- `POST /tests/{id}/sync-moodle` - Синхронизация теста с Moodle
- `GET /moodle/courses` - Получение списка курсов Moodle
- `GET /moodle/validate` - Проверка подключения к Moodle

**Все эндпоинты (кроме `/auth/register` и `/auth/login`) требуют аутентификации через JWT токен.**

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

## Мониторинг и Логирование

### Доступ к интерфейсам

- **Grafana Dashboard**: http://localhost:3000 (admin/admin123)
- **Prometheus (метрики)**: http://localhost:9090
- **Loki (логи)**: http://localhost:3100
- **API Metrics**: http://localhost:8080/metrics
- **Health Check**: http://localhost:8080/health

### Готовый Dashboard в Grafana

После запуска `docker-compose up -d` дашборд **TestGen API Dashboard** автоматически загружен в Grafana!

**Что показывает:**
- Request Rate (запросы/сек)
- Total Requests (общее количество)
- Response Time Percentiles (p50, p95, p99)
- Requests In Progress (активные запросы)
- HTTP Status Codes (2xx, 4xx, 5xx)
- Requests by Endpoint (pie chart)
- Requests by Method (pie chart)

### Собираемые метрики

Автоматически через `fiberprometheus` middleware:

- `http_requests_total` - счётчик всех HTTP запросов (method, path, status_code)
- `http_request_duration_seconds` - время выполнения (histogram для p50/p95/p99)
- `http_requests_in_progress_total` - активные запросы

### Примеры PromQL запросов

```promql
# Request rate за 5 минут
rate(http_requests_total{job="testgen-backend"}[5m])

# 95th percentile response time
histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le, path))

# Количество ошибок 5xx
sum(rate(http_requests_total{status_code=~"5.."}[5m]))
```

### Тестирование мониторинга

Сгенерируйте тестовый трафик для проверки графиков:

```bash
# Простой способ
for i in {1..100}; do curl -s http://localhost/health > /dev/null & done

# Нагрузочное тестирование с ab (Apache Bench)
ab -n 1000 -c 10 http://localhost/health
```

Затем откройте Grafana и посмотрите красивые графики!

### Логирование с Loki

**Loki** - система агрегации логов от Grafana Labs, работающая как "Prometheus для логов".

#### Архитектура логирования:

```
Go Backend (JSON logs) → Docker stdout → Promtail → Loki → Grafana
```

#### Компоненты:

- **Loki** - хранилище и индексация логов
- **Promtail** - агент для сбора логов из Docker контейнеров
- **Grafana** - визуализация и поиск логов

#### Структурированные логи (JSON):

Backend автоматически логирует в JSON формате с полями:
- `level` - уровень логирования (debug, info, warn, error)
- `timestamp` - время события
- `message` - текст сообщения
- `caller` - место в коде
- `method`, `path`, `status`, `duration` - для HTTP запросов
- `user_id`, `request_id` - контекстная информация

#### Просмотр логов в Grafana:

1. Откройте Grafana: <http://localhost:3000>
2. Перейдите в **Explore** (компас слева)
3. Выберите datasource **Loki**
4. Используйте LogQL запросы:

```logql
# Все логи backend
{service="backend"}

# Только ошибки
{service="backend"} |= "level=error"

# Логи конкретного пользователя
{service="backend"} | json | user_id="123e4567-e89b-12d3-a456-426614174000"

# HTTP запросы со статусом 500
{service="backend"} | json | status="500"

# Логи за последние 5 минут с фильтром
{service="backend"} | json | level="error" | line_format "{{.message}}"
```

#### Настройка формата логов:

В `.env` файле:
```bash
LOG_LEVEL=info      # debug, info, warn, error
LOG_FORMAT=json     # console (dev) или json (production)
```

**Важно**: Для корректной работы Promtail backend должен логировать в **JSON формате**!

## Безопасность

- **JWT аутентификация** - HTTP-only cookies + Authorization header
- **RBAC** (Role-Based Access Control) - система управления ролями с таблицей в БД
- **Bcrypt hashing** для паролей
- **HTTPS/TLS** шифрование
- **Валидация** входных данных
- **SQL injection** защита (GORM prepared statements)
- **XSS защита** (HTTP-only cookies, CSP headers)

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

Система управления ролями реализована с использованием таблицы `roles` в БД:

- **Admin**: Управление пользователями (назначение ролей), просмотр всех данных
- **Teacher**: Загрузка документов, генерация и редактирование тестов, экспорт в Moodle
- **Student**: Роль по умолчанию при регистрации (для будущего расширения)

**Важно**: При регистрации пользователь автоматически получает роль `student`. Только администратор может изменять роли через API `/api/v1/users/:id/role`.

## Основные возможности

### Аутентификация

- Регистрация с автоматическим назначением роли student
- Вход с JWT токеном в HTTP-only cookie (защита от XSS)
- Поддержка двух способов авторизации: Cookie (приоритет) и Authorization header
- Выход с очисткой cookie

### Управление пользователями (Admin)

- Просмотр списка всех пользователей с пагинацией
- Изменение ролей пользователей (admin/teacher/student)
- Защита endpoint'ов middleware для проверки роли администратора

### Генерация тестов

- Загрузка документов (PDF, DOCX, PPTX, TXT)
- Автоматическая генерация вопросов с использованием LLM
- Редактирование и экспорт тестов в формат Moodle XML
