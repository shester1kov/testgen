# Test Generation System - Backend

Go backend для системы генерации тестовых заданий на основе документов с интеграцией в Moodle.

## Документация API

После запуска сервера, Swagger UI доступен по адресу:
**<http://localhost:8080/swagger/index.html>**

## Архитектурные паттерны

1. **Repository Pattern** - Абстракция доступа к данным
2. **Factory Pattern** - Создание парсеров документов
3. **Strategy Pattern** - Выбор LLM провайдера
4. **Dependency Injection** - Wire для автоматического внедрения зависимостей
5. **Middleware Chain** - Обработка сквозной функциональности

## Быстрый старт

### Локальная разработка

```bash
# 1. Установить зависимости
go mod download

# 2. Создать .env файл
cp ../.env.example ../.env

# 3. Запустить PostgreSQL (через Docker)
docker-compose -f ../docker-compose.yml up -d postgres

# 4. Запустить сервер
go run cmd/api/main.go
```

### С Docker

```bash
# Запустить все сервисы
cd ..
docker-compose up -d backend
```

## API Endpoints

### Аутентификация

#### POST /api/v1/auth/register
Регистрация нового пользователя с автоматическим назначением роли student.

**Тело запроса:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "full_name": "Иван Иванов"
}
```

**Ответ (201 Created):**
```json
{
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "full_name": "Иван Иванов",
    "role": "student"
  },
  "token": "jwt-token"
}
```

**Возможные ошибки:**
- 400: Некорректные данные
- 409: Email уже зарегистрирован
- 500: Внутренняя ошибка сервера

---

#### POST /api/v1/auth/login
Вход пользователя и получение JWT токена.

**Тело запроса:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Ответ (200 OK):**
```json
{
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "full_name": "Иван Иванов",
    "role": "student"
  },
  "token": "jwt-token"
}
```

**Примечание:** Токен также устанавливается в HTTP-only cookie `testgen_token`.

**Возможные ошибки:**
- 400: Некорректные данные
- 401: Неверный email или пароль
- 500: Внутренняя ошибка сервера

---

#### POST /api/v1/auth/logout
Выход пользователя (очистка cookie).

**Ответ (200 OK):**
```json
{
  "message": "Logged out successfully"
}
```

---

#### GET /api/v1/auth/me
Получение информации о текущем аутентифицированном пользователе.

**Заголовки:**
```
Authorization: Bearer <jwt-token>
```

**Ответ (200 OK):**
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "full_name": "Иван Иванов",
  "role": "student"
}
```

**Возможные ошибки:**
- 401: Не авторизован
- 404: Пользователь не найден

---

### Документы

#### POST /api/v1/documents
Загрузка документа для парсинга.

**Заголовки:**
```
Authorization: Bearer <jwt-token>
Content-Type: multipart/form-data
```

**Параметры формы:**
- `file` (обязательно): Файл документа (PDF, DOCX, PPTX, TXT, MD)
- `title` (опционально): Название документа

**Ответ (201 Created):**
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "title": "Лекция по математике",
  "file_name": "lecture.pdf",
  "file_type": "pdf",
  "file_size": 1024000,
  "status": "uploaded",
  "parsed_text": null,
  "created_at": "2024-01-20T15:04:05Z"
}
```

**Возможные ошибки:**
- 400: Некорректный файл или неподдерживаемый формат
- 401: Не авторизован
- 500: Внутренняя ошибка сервера

---

#### GET /api/v1/documents
Получение списка документов с пагинацией.

**Заголовки:**
```
Authorization: Bearer <jwt-token>
```

**Query параметры:**
- `page` (по умолчанию: 1): Номер страницы
- `page_size` (по умолчанию: 20): Размер страницы

**Ответ (200 OK):**
```json
{
  "documents": [
    {
      "id": "uuid",
      "user_id": "uuid",
      "title": "Лекция 1",
      "file_name": "lecture1.pdf",
      "file_type": "pdf",
      "file_size": 1024000,
      "status": "parsed",
      "created_at": "2024-01-20T15:04:05Z",
      "user_name": "Иван Иванов",
      "user_email": "user@example.com"
    }
  ],
  "total": 100,
  "page": 1,
  "page_size": 20
}
```

**Примечание:** Admin видит все документы с информацией о владельце, остальные пользователи видят только свои документы.

**Возможные ошибки:**
- 401: Не авторизован
- 500: Внутренняя ошибка сервера

---

#### GET /api/v1/documents/:id
Получение информации о конкретном документе.

**Заголовки:**
```
Authorization: Bearer <jwt-token>
```

**Ответ (200 OK):**
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "title": "Лекция по математике",
  "file_name": "lecture.pdf",
  "file_type": "pdf",
  "file_size": 1024000,
  "status": "parsed",
  "parsed_text": "Текст документа...",
  "created_at": "2024-01-20T15:04:05Z"
}
```

**Возможные ошибки:**
- 400: Некорректный ID документа
- 401: Не авторизован
- 403: Доступ запрещен
- 404: Документ не найден

---

#### POST /api/v1/documents/:id/parse
Парсинг текста из загруженного документа.

**Заголовки:**
```
Authorization: Bearer <jwt-token>
```

**Ответ (200 OK):**
```json
{
  "id": "uuid",
  "status": "parsed",
  "parsed_text": "Полный текст документа...",
  "text_preview": "Первые 200 символов..."
}
```

**Возможные ошибки:**
- 400: Некорректный ID или неподдерживаемый формат
- 401: Не авторизован
- 403: Доступ запрещен
- 404: Документ не найден
- 500: Ошибка парсинга

---

#### DELETE /api/v1/documents/:id
Удаление документа и связанного файла.

**Заголовки:**
```
Authorization: Bearer <jwt-token>
```

**Ответ (200 OK):**
```json
{
  "message": "Document deleted successfully"
}
```

**Возможные ошибки:**
- 400: Некорректный ID документа
- 401: Не авторизован
- 403: Доступ запрещен
- 404: Документ не найден
- 500: Внутренняя ошибка сервера

---

### Тесты

#### POST /api/v1/tests
Создание нового теста.

**Заголовки:**
```
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Тело запроса:**
```json
{
  "title": "Тест по математике",
  "description": "Контрольная работа №1",
  "document_id": "uuid"
}
```

**Ответ (201 Created):**
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "title": "Тест по математике",
  "description": "Контрольная работа №1",
  "total_questions": 0,
  "status": "draft",
  "moodle_synced": false,
  "questions": [],
  "created_at": "2024-01-20T15:04:05Z"
}
```

**Возможные ошибки:**
- 400: Некорректные данные
- 401: Не авторизован
- 500: Внутренняя ошибка сервера

---

#### POST /api/v1/tests/generate
Генерация вопросов теста с помощью LLM.

**Заголовки:**
```
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Тело запроса:**
```json
{
  "title": "Тест по математике",
  "document_id": "uuid",
  "num_questions": 20,
  "difficulty": "medium",
  "question_types": ["single_choice"],
  "llm_provider": "perplexity"
}
```

**Параметры:**
- `title` (обязательно): Название теста (минимум 3 символа)
- `document_id` (обязательно): UUID документа
- `num_questions` (обязательно): Количество вопросов (1-50)
- `difficulty` (обязательно): Сложность - `easy`, `medium`, `hard`
- `question_types` (опционально): Типы вопросов
- `llm_provider` (опционально): Провайдер LLM - `perplexity`, `openai`, `yandexgpt`

**Ответ (201 Created):**
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "title": "Тест по математике",
  "description": "",
  "total_questions": 20,
  "status": "draft",
  "moodle_synced": false,
  "questions": [
    {
      "id": "uuid",
      "question_text": "Чему равна производная x²?",
      "question_type": "single_choice",
      "difficulty": "medium",
      "points": 1.0,
      "order_num": 1,
      "answers": [
        {
          "id": "uuid",
          "answer_text": "2x",
          "is_correct": true,
          "order_num": 1
        },
        {
          "id": "uuid",
          "answer_text": "x",
          "is_correct": false,
          "order_num": 2
        }
      ]
    }
  ],
  "created_at": "2024-01-20T15:04:05Z"
}
```

**Возможные ошибки:**
- 400: Некорректные данные или документ не распарсен
- 401: Не авторизован
- 404: Документ не найден
- 500: Ошибка генерации или ошибка БД

---

#### GET /api/v1/tests
Получение списка тестов с пагинацией.

**Заголовки:**
```
Authorization: Bearer <jwt-token>
```

**Query параметры:**
- `page` (по умолчанию: 1): Номер страницы
- `page_size` (по умолчанию: 20): Размер страницы

**Ответ (200 OK):**
```json
{
  "tests": [
    {
      "id": "uuid",
      "user_id": "uuid",
      "title": "Тест по математике",
      "description": "Контрольная работа",
      "total_questions": 20,
      "status": "draft",
      "moodle_synced": false,
      "created_at": "2024-01-20T15:04:05Z",
      "user_name": "Иван Иванов",
      "user_email": "user@example.com"
    }
  ],
  "total": 50,
  "page": 1,
  "page_size": 20
}
```

**Примечание:** Admin видит все тесты, остальные пользователи видят только свои тесты.

**Возможные ошибки:**
- 401: Не авторизован
- 500: Внутренняя ошибка сервера

---

#### GET /api/v1/tests/:id
Получение информации о тесте с вопросами и ответами.

**Заголовки:**
```
Authorization: Bearer <jwt-token>
```

**Ответ (200 OK):**
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "title": "Тест по математике",
  "description": "Контрольная работа",
  "total_questions": 20,
  "status": "draft",
  "moodle_synced": false,
  "questions": [
    {
      "id": "uuid",
      "question_text": "Вопрос 1",
      "question_type": "single_choice",
      "difficulty": "medium",
      "points": 1.0,
      "order_num": 1,
      "answers": [...]
    }
  ],
  "created_at": "2024-01-20T15:04:05Z"
}
```

**Возможные ошибки:**
- 400: Некорректный ID теста
- 401: Не авторизован
- 403: Доступ запрещен
- 404: Тест не найден

---

#### PUT /api/v1/tests/:id
Обновление названия и описания теста.

**Заголовки:**
```
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Тело запроса:**
```json
{
  "title": "Новое название",
  "description": "Новое описание"
}
```

**Ответ (200 OK):**
```json
{
  "id": "uuid",
  "title": "Новое название",
  "description": "Новое описание",
  ...
}
```

**Возможные ошибки:**
- 400: Некорректные данные
- 401: Не авторизован
- 403: Доступ запрещен
- 404: Тест не найден
- 500: Внутренняя ошибка сервера

---

#### DELETE /api/v1/tests/:id
Удаление теста и всех связанных вопросов.

**Заголовки:**
```
Authorization: Bearer <jwt-token>
```

**Ответ (200 OK):**
```json
{
  "message": "Test deleted successfully"
}
```

**Возможные ошибки:**
- 400: Некорректный ID теста
- 401: Не авторизован
- 403: Доступ запрещен
- 404: Тест не найден
- 500: Внутренняя ошибка сервера

---

#### PUT /api/v1/tests/:testId/questions/:questionId
Обновление вопроса (текст, тип, сложность, баллы, ответы).

**Заголовки:**
```
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Тело запроса:**
```json
{
  "question_text": "Обновленный текст вопроса",
  "question_type": "single_choice",
  "difficulty": "hard",
  "points": 2.0,
  "answers": [
    {
      "id": "uuid",
      "answer_text": "Вариант 1",
      "is_correct": true,
      "order_num": 1
    },
    {
      "answer_text": "Вариант 2",
      "is_correct": false,
      "order_num": 2
    }
  ]
}
```

**Примечание:** Если у ответа нет `id`, будет создан новый ответ.

**Ответ (200 OK):**
```json
{
  "id": "uuid",
  "question_text": "Обновленный текст вопроса",
  "question_type": "single_choice",
  "difficulty": "hard",
  "points": 2.0,
  "order_num": 1,
  "answers": [...]
}
```

**Возможные ошибки:**
- 400: Некорректные данные
- 401: Не авторизован
- 403: Доступ запрещен
- 404: Вопрос не найден
- 500: Внутренняя ошибка сервера

---

### Экспорт тестов

#### GET /api/v1/tests/:id/export/json
Экспорт теста в JSON формате для скачивания.

**Заголовки:**
```
Authorization: Bearer <jwt-token>
```

**Ответ (200 OK):**
```json
{
  "id": "uuid",
  "title": "Тест по математике",
  "questions": [...]
}
```

**Возможные ошибки:**
- 400: Некорректный ID теста
- 401: Не авторизован
- 403: Доступ запрещен
- 404: Тест не найден
- 500: Внутренняя ошибка сервера

---

#### GET /api/v1/tests/:id/export/xml
Экспорт теста в формате Moodle XML для скачивания.

**Заголовки:**
```
Authorization: Bearer <jwt-token>
```

**Ответ (200 OK):**
```xml
<?xml version="1.0" encoding="UTF-8"?>
<quiz>
  <question type="multichoice">
    <name><text>Вопрос 1</text></name>
    <questiontext format="html">
      <text><![CDATA[Чему равна производная x²?]]></text>
    </questiontext>
    <answer fraction="100">
      <text>2x</text>
    </answer>
    ...
  </question>
</quiz>
```

**Возможные ошибки:**
- 400: Некорректный ID или тест без вопросов
- 401: Не авторизован
- 403: Доступ запрещен
- 404: Тест не найден
- 500: Ошибка экспорта

---

### Moodle интеграция

#### GET /api/v1/moodle/validate
Проверка подключения к серверу Moodle.

**Заголовки:**
```
Authorization: Bearer <jwt-token>
```

**Ответ (200 OK):**
```json
{
  "connected": true,
  "message": "Moodle connection successful"
}
```

**Ответ (503 Service Unavailable):**
```json
{
  "connected": false,
  "message": "Connection failed",
  "error": "Timeout"
}
```

**Возможные ошибки:**
- 401: Не авторизован
- 503: Ошибка подключения

---

#### GET /api/v1/moodle/courses
Получение списка доступных курсов из Moodle.

**Заголовки:**
```
Authorization: Bearer <jwt-token>
```

**Ответ (200 OK):**
```json
{
  "courses": [
    {
      "id": "1",
      "name": "Математический анализ",
      "short_name": "MATH101"
    },
    {
      "id": "2",
      "name": "Линейная алгебра",
      "short_name": "MATH102"
    }
  ]
}
```

**Возможные ошибки:**
- 401: Не авторизован
- 500: Не удалось получить курсы

---

#### GET /api/v1/moodle/tests/:id/export
Экспорт теста в Moodle XML (альтернативный эндпоинт).

**Заголовки:**
```
Authorization: Bearer <jwt-token>
```

**Ответ (200 OK):**
```xml
<?xml version="1.0" encoding="UTF-8"?>
<quiz>...</quiz>
```

**Возможные ошибки:**
- 400: Некорректный ID или тест без вопросов
- 401: Не авторизован
- 404: Тест не найден
- 500: Ошибка экспорта

---

#### POST /api/v1/moodle/tests/:id/sync
Синхронизация теста с Moodle (загрузка как quiz в курс).

**Заголовки:**
```
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Тело запроса:**
```json
{
  "course_name": "MATH101"
}
```

**Ответ (200 OK):**
```json
{
  "message": "Test synced successfully",
  "moodle_id": "quiz_12345",
  "course_id": "1"
}
```

**Возможные ошибки:**
- 400: Некорректные данные или тест без вопросов
- 401: Не авторизован
- 404: Тест не найден
- 500: Ошибка синхронизации

---

### Пользователи (Admin only)

#### GET /api/v1/users
Получение списка всех пользователей (только для admin).

**Заголовки:**
```
Authorization: Bearer <jwt-token>
```

**Query параметры:**
- `limit` (по умолчанию: 10): Количество пользователей
- `offset` (по умолчанию: 0): Смещение

**Ответ (200 OK):**
```json
{
  "users": [
    {
      "id": "uuid",
      "email": "user@example.com",
      "full_name": "Иван Иванов",
      "role": "student"
    }
  ],
  "total": 100,
  "limit": 10,
  "offset": 0
}
```

**Возможные ошибки:**
- 401: Не авторизован
- 403: Доступ запрещен (не admin)
- 500: Внутренняя ошибка сервера

---

#### PUT /api/v1/users/:id/role
Изменение роли пользователя (только для admin).

**Заголовки:**
```
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Тело запроса:**
```json
{
  "role_name": "teacher"
}
```

**Допустимые роли:** `admin`, `teacher`, `student`

**Ответ (200 OK):**
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "full_name": "Иван Иванов",
  "role": "teacher"
}
```

**Возможные ошибки:**
- 400: Некорректная роль
- 401: Не авторизован
- 403: Доступ запрещен (не admin)
- 404: Пользователь не найден
- 500: Внутренняя ошибка сервера

---

### Статистика

#### GET /api/v1/stats/dashboard
Получение статистики для dashboard (количество документов, тестов, вопросов).

**Заголовки:**
```
Authorization: Bearer <jwt-token>
```

**Ответ (200 OK):**
```json
{
  "documents_count": 15,
  "tests_count": 8,
  "questions_count": 120
}
```

**Возможные ошибки:**
- 401: Не авторизован
- 500: Внутренняя ошибка сервера

---

### Мониторинг

#### GET /health
Проверка состояния сервера.

**Ответ (200 OK):**
```json
{
  "status": "ok",
  "timestamp": "2024-01-20T15:04:05Z"
}
```

---

#### GET /metrics
Метрики Prometheus для мониторинга.

**Ответ (200 OK):**
```
# HELP http_requests_total Total number of HTTP requests
# TYPE http_requests_total counter
http_requests_total{method="GET",endpoint="/api/v1/tests"} 1234
...
```

---

## Тестирование

```bash
# Запустить все тесты
go test ./... -v

# С покрытием
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Структура

```
backend/
├── cmd/api/                          # Entry point
├── internal/
│   ├── domain/                       # Бизнес-логика
│   │   ├── entity/                   # Сущности
│   │   └── repository/               # Интерфейсы репозиториев
│   ├── application/                  # Use cases
│   │   └── dto/                      # Data Transfer Objects
│   ├── infrastructure/               # Внешние зависимости
│   │   └── persistence/
│   │       ├── postgres/             # PostgreSQL реализация
│   │       └── migrations/           # SQL миграции
│   └── interfaces/                   # HTTP layer
│       └── http/
│           ├── handler/              # Request handlers
│           ├── middleware/           # Middleware
│           └── router/               # Routes setup
└── pkg/                              # Shared packages
    ├── config/                       # Configuration
    ├── logger/                       # Logging
    ├── errors/                       # Error handling
    └── utils/                        # Utilities (JWT)
```

## Что реализовано (MVP)

- Clean Architecture (Domain, Application, Infrastructure, Interfaces)
- PostgreSQL + GORM
- Database migrations
- JWT Authentication
- User Registration & Login
- Role-Based Access Control (RBAC)
- Error handling
- Structured logging
- Configuration management
- CORS middleware
- Health check endpoint

## TODO

- [ ] Document upload handler
- [ ] Document parsers (PDF, DOCX, PPTX, TXT)
- [ ] LLM integration (Perplexity API)
- [ ] Test generation logic
- [ ] Moodle XML export
- [ ] Prometheus metrics
- [ ] Unit & Integration tests
- [ ] Wire dependency injection
