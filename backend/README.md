# Test Generation System - Backend

Go backend для системы генерации тестовых заданий на основе документов с интеграцией в Moodle.

## 📚 Документация API

После запуска сервера, Swagger UI доступен по адресу:
**http://localhost:8080/swagger/index.html**

Для генерации Swagger документации:
```bash
make swagger
```

## 🏗️ Архитектурные паттерны

1. **Repository Pattern** - Абстракция доступа к данным
2. **Factory Pattern** - Создание парсеров документов
3. **Strategy Pattern** - Выбор LLM провайдера
4. **Dependency Injection** - Wire для автоматического внедрения зависимостей
5. **Middleware Chain** - Обработка сквозной функциональности

## 🚀 Быстрый старт

### Локальная разработка

```bash
# 1. Установить зависимости
go mod download

# 2. Создать .env файл
cp ../.env.example ../.env

# 3. Запустить PostgreSQL (через Docker)
docker-compose -f ../docker-compose.yml up -d postgres

# 4. Применить миграции
migrate -path internal/infrastructure/persistence/migrations \
        -database "postgres://testgen_user:testgen_password@localhost:5432/testgen_db?sslmode=disable" up

# 5. Запустить сервер
go run cmd/api/main.go
```

### С Docker

```bash
# Запустить все сервисы
cd ..
docker-compose up -d backend
```

## 📚 API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Регистрация пользователя
- `POST /api/v1/auth/login` - Вход в систему
- `GET /api/v1/auth/me` - Получить текущего пользователя (требует авторизации)

### Health Check
- `GET /health` - Проверка работоспособности

## 🧪 Тестирование

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
- ✅ PostgreSQL + GORM
- ✅ Database migrations
- ✅ JWT Authentication
- ✅ User Registration & Login
- ✅ Role-Based Access Control (RBAC)
- ✅ Error handling
- ✅ Structured logging
- ✅ Configuration management
- ✅ CORS middleware
- ✅ Health check endpoint

## 🔜 TODO

- [ ] Document upload handler
- [ ] Document parsers (PDF, DOCX, PPTX, TXT)
- [ ] LLM integration (Perplexity API)
- [ ] Test generation logic
- [ ] Moodle XML export
- [ ] Prometheus metrics
- [ ] Unit & Integration tests
- [ ] Wire dependency injection
