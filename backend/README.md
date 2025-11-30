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

## Оформление документации по REST-интерфейсу

Документация составлена на основе Swagger-спецификации `backend/docs/swagger.yaml` (базовый префикс: `http://localhost:8080/api/v1`). Все ответы возвращают JSON, если не указано иное. Защищённые маршруты требуют заголовок `Authorization: Bearer <JWT>`.

## GET-запросы

### Получить данные текущего пользователя

**Путь:** `GET /auth/me`

**Тестовые данные:** отправляйте после успешной авторизации (`/auth/login`).

**Тело запроса:** отсутствует.

**Тело ответа (200 OK):**

```json
{
  "id": "c9d9f3c4-1c1c-4c70-9f3a-5b2c9f6a1a1a",
  "email": "admin@example.com",
  "full_name": "Test Admin",
  "role": "admin"
}
```

**Тело ответа (401 Unauthorized):**

```json
{
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Missing or invalid token"
  }
}
```

---

### Получить список документов (с пагинацией)

**Путь:** `GET /documents?page=1&page_size=20`

**Тестовые данные:** параметры по умолчанию `page=1`, `page_size=20`.

**Тело запроса:** отсутствует.

**Тело ответа (200 OK):**

```json
{
  "documents": [
    {
      "id": "a12b34cd-56ef-78ab-90cd-ef1234567890",
      "title": "Лекция 1",
      "file_name": "lecture1.pdf",
      "file_size": 512000,
      "file_type": "pdf",
      "status": "parsed",
      "created_at": "2024-05-01T12:00:00Z",
      "user_name": "Test Admin",
      "user_email": "admin@example.com"
    }
  ],
  "page": 1,
  "page_size": 20,
  "total": 1
}
```

**Тело ответа (401 Unauthorized):** структура ошибки аналогична разделу авторизации.

---

### Получить один документ

**Путь:** `GET /documents/{id}` (пример: `/documents/a12b34cd-56ef-78ab-90cd-ef1234567890`)

**Тело запроса:** отсутствует.

**Тело ответа (200 OK):**

```json
{
  "id": "a12b34cd-56ef-78ab-90cd-ef1234567890",
  "title": "Лекция 1",
  "file_name": "lecture1.pdf",
  "file_size": 512000,
  "file_type": "pdf",
  "status": "parsed",
  "parsed_text": "Вырезка содержимого...",
  "created_at": "2024-05-01T12:00:00Z",
  "user_name": "Test Admin",
  "user_email": "admin@example.com"
}
```

**Тело ответа (404 Not Found):**

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Document not found"
  }
}
```

---

### Получить статистику для дашборда

**Путь:** `GET /stats/dashboard`

**Тело запроса:** отсутствует.

**Тело ответа (200 OK):**

```json
{
  "documents_count": 15,
  "tests_count": 8,
  "questions_count": 120
}
```

---

### Получить список тестов (с пагинацией)

**Путь:** `GET /tests?page=1&page_size=20`

**Тело запроса:** отсутствует.

**Тело ответа (200 OK):**

```json
{
  "tests": [
    {
      "id": "4c5d6e7f-8a9b-0c1d-2e3f-4a5b6c7d8e9f",
      "title": "Тест по лекции 1",
      "description": "10 вопросов по первой лекции",
      "document_id": "a12b34cd-56ef-78ab-90cd-ef1234567890",
      "created_at": "2024-05-02T10:00:00Z"
    }
  ],
  "page": 1,
  "page_size": 20,
  "total": 1
}
```

**Тело ответа (401 Unauthorized):** структура ошибки аналогична разделу авторизации.

---

### Получить один тест

**Путь:** `GET /tests/{id}` (пример: `/tests/4c5d6e7f-8a9b-0c1d-2e3f-4a5b6c7d8e9f`)

**Тело запроса:** отсутствует.

**Тело ответа (200 OK):**

```json
{
  "id": "4c5d6e7f-8a9b-0c1d-2e3f-4a5b6c7d8e9f",
  "title": "Тест по лекции 1",
  "description": "10 вопросов по первой лекции",
  "document_id": "a12b34cd-56ef-78ab-90cd-ef1234567890",
  "questions": [
    {
      "id": "11112222-3333-4444-5555-666677778888",
      "question_text": "Что такое REST?",
      "question_type": "single_choice",
      "difficulty": "easy",
      "points": 1,
      "answers": [
        { "id": "ans1", "answer_text": "Подход к построению API", "is_correct": true, "order_num": 1 },
        { "id": "ans2", "answer_text": "Тип базы данных", "is_correct": false, "order_num": 2 }
      ]
    }
  ]
}
```

**Тело ответа (404 Not Found):** структура ошибки аналогична разделу документов.

---

### Проверить подключение к Moodle

**Путь:** `GET /moodle/validate`

**Тело запроса:** отсутствует.

**Тело ответа (200 OK):**

```json
{
  "connected": true,
  "message": "Connection successful",
  "error": null
}
```

## POST-запросы

### Авторизация пользователя

**Путь:** `POST /auth/login`

**Тестовые данные:** укажите зарегистрированные учётные данные.

**Тело запроса (application/json):**

```json
{
  "email": "admin@example.com",
  "password": "AdminPass123"
}
```

**Тело ответа (200 OK):**

```json
{
  "token": "<jwt-token>",
  "user": {
    "id": "c9d9f3c4-1c1c-4c70-9f3a-5b2c9f6a1a1a",
    "email": "admin@example.com",
    "full_name": "Test Admin",
    "role": "admin"
  }
}
```

**Тело ответа (401 Unauthorized):**

```json
{
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Invalid credentials"
  }
}
```

---

### Регистрация пользователя

**Путь:** `POST /auth/register`

**Тело запроса (application/json):**

```json
{
  "email": "student@example.com",
  "full_name": "Student User",
  "password": "StudentPass123"
}
```

**Тело ответа (201 Created):**

```json
{
  "token": "<jwt-token>",
  "user": {
    "id": "2f4e6a8c-0b1c-2d3e-4f5a-6b7c8d9e0f1a",
    "email": "student@example.com",
    "full_name": "Student User",
    "role": "student"
  }
}
```

**Тело ответа (409 Conflict):**

```json
{
  "error": {
    "code": "CONFLICT",
    "message": "User already exists"
  }
}
```

---

### Загрузка документа

**Путь:** `POST /documents`

**Тело запроса (multipart/form-data):**

- `file` (обязательное поле) – загружаемый файл (`pdf`, `docx`, `pptx`, `txt`, `md`).
- `title` (опционально) – заголовок документа.

**Тело ответа (201 Created):**

```json
{
  "id": "a12b34cd-56ef-78ab-90cd-ef1234567890",
  "title": "Лекция 1",
  "file_name": "lecture1.pdf",
  "file_size": 512000,
  "file_type": "pdf",
  "status": "uploaded",
  "created_at": "2024-05-01T12:00:00Z"
}
```

**Тело ответа (400 Bad Request):**

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input or unsupported file type"
  }
}
```

---

### Запуск парсинга документа

**Путь:** `POST /documents/{id}/parse`

**Тело запроса:** отсутствует.

**Тело ответа (200 OK):**

```json
{
  "id": "a12b34cd-56ef-78ab-90cd-ef1234567890",
  "status": "parsed",
  "parsed_text": "Вырезка содержимого...",
  "text_preview": "Вырезка содержимого..."
}
```

---

### Создание теста вручную

**Путь:** `POST /tests`

**Тело запроса (application/json):**

```json
{
  "title": "Тест по лекции 1",
  "description": "10 вопросов по первой лекции",
  "document_id": "a12b34cd-56ef-78ab-90cd-ef1234567890"
}
```

**Тело ответа (201 Created):** содержит созданный тест с массивом вопросов (может быть пустым на старте).

---

### Генерация теста на основе LLM

**Путь:** `POST /tests/generate`

**Тело запроса (application/json):**

```json
{
  "title": "Тест по лекции 1",
  "difficulty": "medium",
  "document_id": "a12b34cd-56ef-78ab-90cd-ef1234567890",
  "num_questions": 5,
  "llm_provider": "yandexgpt",
  "question_types": ["single_choice", "multiple_choice"]
}
```

**Тело ответа (201 Created):**

```json
{
  "id": "4c5d6e7f-8a9b-0c1d-2e3f-4a5b6c7d8e9f",
  "title": "Тест по лекции 1",
  "description": "Автоматически сгенерированные вопросы",
  "document_id": "a12b34cd-56ef-78ab-90cd-ef1234567890",
  "questions": [ /* массив вопросов с вариантами ответов */ ]
}
```

---

### Синхронизация теста с Moodle

**Путь:** `POST /moodle/sync/{id}`

**Тело запроса (application/json):**

```json
{
  "course_name": "API Basics"
}
```

**Тело ответа (200 OK):**

```json
{
  "course_id": "58",
  "moodle_id": "1024",
  "message": "Test synced successfully"
}
```

## PUT-запросы

### Обновление теста

**Путь:** `PUT /tests/{id}`

**Тело запроса (application/json):**

```json
{
  "title": "Тест по лекции 1 (обновлён)",
  "description": "Расширенный тест"
}
```

**Тело ответа (200 OK):** обновлённый тест с полями, аналогичными `GET /tests/{id}`.

---

### Обновление вопроса

**Путь:** `PUT /tests/{testId}/questions/{questionId}`

**Тело запроса (application/json):** включает текст вопроса, сложность, тип и массив ответов с флагами `is_correct`.

**Тело ответа (200 OK):**

```json
{
  "id": "11112222-3333-4444-5555-666677778888",
  "question_text": "Что такое REST?",
  "question_type": "single_choice",
  "difficulty": "easy",
  "points": 1,
  "answers": [
    { "id": "ans1", "answer_text": "Подход к построению API", "is_correct": true, "order_num": 1 }
  ]
}
```

---

### Смена роли пользователя (администратор)

**Путь:** `PUT /users/{id}/role`

**Тело запроса (application/json):**

```json
{
  "role_name": "teacher"
}
```

**Тело ответа (200 OK):**

```json
{
  "id": "2f4e6a8c-0b1c-2d3e-4f5a-6b7c8d9e0f1a",
  "email": "student@example.com",
  "full_name": "Student User",
  "role": "teacher"
}
```

## DELETE-запросы

### Удаление документа

**Путь:** `DELETE /documents/{id}`

**Тело ответа (200 OK):**

```json
{
  "message": "Operation completed successfully"
}
```

### Удаление теста

**Путь:** `DELETE /tests/{id}`

**Тело ответа (200 OK):**

```json
{
  "message": "Operation completed successfully"
}
```

## Тестирование

```bash
# Запустить все тесты
go test ./... -v

# С покрытием
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 📁 Структура

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

## ✅ Что реализовано (MVP)

- ✅ Clean Architecture (Domain, Application, Infrastructure, Interfaces)
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

## 🔜 TODO

- [ ] Document upload handler
- [ ] Document parsers (PDF, DOCX, PPTX, TXT)
- [ ] LLM integration (Perplexity API)
- [ ] Test generation logic
- [ ] Moodle XML export
- [ ] Prometheus metrics
- [ ] Unit & Integration tests
- [ ] Wire dependency injection
