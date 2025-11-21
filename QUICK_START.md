# 🚀 Быстрый старт бэкенда

Пошаговая инструкция для запуска бэкенда Test Generation System.

## Предварительные требования

- Go 1.23+ установлен ([скачать](https://golang.org/dl/))
- Docker Desktop установлен ([скачать](https://www.docker.com/products/docker-desktop))
- Git Bash или WSL (для Windows)

Проверьте установку:
```bash
go version      # должно показать go1.21 или выше
docker --version
```

---

## Шаг 1: Инициализация Go модуля

```bash
cd "c:\Users\shest\Desktop\course work\backend"

# Скачать все зависимости
go mod download

# Если есть ошибки, попробуйте:
go mod tidy
```

**Ожидаемый результат**: Все зависимости скачаны без ошибок.

---

## Шаг 2: Создать .env файл

Создайте файл `.env` в корне проекта (`c:\Users\shest\Desktop\course work\.env`):

```bash
cd "c:\Users\shest\Desktop\course work"

# Windows (PowerShell)
Copy-Item .env.example .env

# Или создайте вручную
```

Содержимое `.env`:
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
JWT_SECRET=my-super-secret-jwt-key-change-in-production
JWT_EXPIRATION=24h

# File Upload
MAX_FILE_SIZE=52428800
UPLOAD_DIR=./uploads

# LLM (можно оставить пустыми для тестирования)
PERPLEXITY_API_KEY=
OPENAI_API_KEY=
YANDEX_GPT_API_KEY=
LLM_PROVIDER=perplexity

# Moodle (можно оставить пустыми для тестирования)
MOODLE_URL=
MOODLE_TOKEN=
```

---

## Шаг 3: Запустить PostgreSQL через Docker

```bash
cd "c:\Users\shest\Desktop\course work"

# Запустить только PostgreSQL
docker-compose up -d postgres

# Проверить, что контейнер запустился
docker ps

# Вы должны увидеть контейнер с именем postgres
```

**Проверка подключения**:
```bash
# Подключиться к PostgreSQL
docker exec -it course-work-postgres-1 psql -U testgen_user -d testgen_db

# В psql консоли выполните:
\dt  # Показать таблицы (пока будет пусто)
\q   # Выйти
```

---

## Шаг 4: Установить golang-migrate (для миграций)

### Windows:
```bash
# Через scoop (если установлен)
scoop install migrate

# Или скачать бинарник вручную:
# https://github.com/golang-migrate/migrate/releases
# Скачайте migrate.windows-amd64.tar.gz
# Распакуйте и добавьте в PATH
```

### Альтернатива - использовать Docker:
```bash
# Создайте alias для migrate через Docker
alias migrate="docker run --rm -v \"c:\Users\shest\Desktop\course work\backend:/app\" --network host migrate/migrate"
```

---

## Шаг 5: Применить миграции БД

```bash
cd "c:\Users\shest\Desktop\course work\backend"

# Вариант 1: Если migrate установлен локально
migrate -path internal/infrastructure/persistence/migrations \
  -database "postgres://testgen_user:testgen_password@localhost:5432/testgen_db?sslmode=disable" up

# Вариант 2: Через Docker
docker run --rm \
  -v "c:\Users\shest\Desktop\course work\backend:/app" \
  --network host \
  migrate/migrate \
  -path /app/internal/infrastructure/persistence/migrations \
  -database "postgres://testgen_user:testgen_password@localhost:5432/testgen_db?sslmode=disable" up
```

**Проверка**:
```bash
# Подключитесь к БД и проверьте таблицы
docker exec -it course-work-postgres-1 psql -U testgen_user -d testgen_db -c "\dt"

# Должны увидеть таблицы: users, documents, tests, questions, answers, activity_logs
```

---

## Шаг 6: Создать директорию для загрузок

```bash
cd "c:\Users\shest\Desktop\course work"
mkdir -p uploads
```

---

## Шаг 7: Сгенерировать Swagger документацию

```bash
cd "c:\Users\shest\Desktop\course work\backend"

# Установить swag
go install github.com/swaggo/swag/cmd/swag@latest

# Сгенерировать документацию
swag init -g cmd/api/main.go -o docs

# Или через Makefile
cd "c:\Users\shest\Desktop\course work"
make swagger
```

**Ожидаемый результат**: Создана папка `backend/docs/` с файлами `docs.go`, `swagger.json`, `swagger.yaml`.

---

## Шаг 8: Запустить бэкенд сервер

```bash
cd "c:\Users\shest\Desktop\course work\backend"

# Запуск
go run cmd/api/main.go

# Или через Makefile
cd "c:\Users\shest\Desktop\course work"
make backend-run
```

**Ожидаемый вывод**:
```
Server started on port 8080
```

---

## Шаг 9: Проверить работу API

Откройте браузер и перейдите:

### 1. Health Check
```
http://localhost:8080/health
```
Должны увидеть:
```json
{
  "status": "ok",
  "message": "Test Generation System API is running"
}
```

### 2. API Info
```
http://localhost:8080/api/v1
```
Должны увидеть:
```json
{
  "message": "Test Generation System API v1",
  "endpoints": {
    "auth": "/api/v1/auth",
    "documents": "/api/v1/documents",
    "tests": "/api/v1/tests",
    "moodle": "/api/v1/moodle"
  }
}
```

### 3. Swagger UI 🎉
```
http://localhost:8080/swagger/index.html
```
Должны увидеть интерактивную документацию API!

---

## Шаг 10: Протестировать API

### Тест 1: Регистрация пользователя

**Через Swagger UI:**
1. Откройте http://localhost:8080/swagger/index.html
2. Найдите `POST /api/v1/auth/register`
3. Нажмите "Try it out"
4. Введите данные:
```json
{
  "email": "teacher@test.com",
  "password": "password123",
  "full_name": "Test Teacher",
  "role": "teacher"
}
```
5. Нажмите "Execute"

**Через curl:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "teacher@test.com",
    "password": "password123",
    "full_name": "Test Teacher",
    "role": "teacher"
  }'
```

**Ожидаемый ответ**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "uuid-here",
    "email": "teacher@test.com",
    "full_name": "Test Teacher",
    "role": "teacher"
  }
}
```

### Тест 2: Вход в систему

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "teacher@test.com",
    "password": "password123"
  }'
```

### Тест 3: Получить текущего пользователя (с авторизацией)

```bash
# Замените YOUR_TOKEN_HERE на токен из предыдущего ответа
curl -X GET http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

---

## Проблемы и решения

### Ошибка: "cannot find package"
```bash
cd backend
go mod tidy
go mod download
```

### Ошибка: "connection refused" (PostgreSQL)
```bash
# Проверьте, что PostgreSQL запущен
docker ps | grep postgres

# Перезапустите контейнер
docker-compose restart postgres

# Проверьте логи
docker-compose logs postgres
```

### Ошибка: "migration failed"
```bash
# Откатите миграции
migrate -path internal/infrastructure/persistence/migrations \
  -database "postgres://testgen_user:testgen_password@localhost:5432/testgen_db?sslmode=disable" down

# Примените заново
migrate -path internal/infrastructure/persistence/migrations \
  -database "postgres://testgen_user:testgen_password@localhost:5432/testgen_db?sslmode=disable" up
```

### Ошибка: "docs package not found"
```bash
cd backend
swag init -g cmd/api/main.go -o docs
go run cmd/api/main.go
```

### Порт 8080 занят
Измените в `.env`:
```env
PORT=8081
```

---

## Остановка сервисов

```bash
# Остановить бэкенд: Ctrl+C

# Остановить PostgreSQL
docker-compose stop postgres

# Или остановить все сервисы
docker-compose down
```

---

## Полезные команды

```bash
# Посмотреть логи PostgreSQL
docker-compose logs -f postgres

# Подключиться к базе данных
docker exec -it course-work-postgres-1 psql -U testgen_user -d testgen_db

# Посмотреть все таблицы
docker exec -it course-work-postgres-1 psql -U testgen_user -d testgen_db -c "\dt"

# Посмотреть пользователей
docker exec -it course-work-postgres-1 psql -U testgen_user -d testgen_db -c "SELECT * FROM users;"

# Сбросить базу данных
docker-compose down -v
docker-compose up -d postgres
# Затем снова применить миграции
```

---

## Следующие шаги

1. ✅ Бэкенд запущен и работает
2. 📝 Протестируйте все endpoints через Swagger UI
3. 🧪 Запустите тесты: `go test ./... -v`
4. 📦 Следующий шаг: Разработка frontend на Vue 3

---

## Структура проекта (что уже реализовано)

✅ **Authentication** - Регистрация, вход, JWT токены
✅ **Database** - PostgreSQL с миграциями
✅ **Clean Architecture** - Domain, Application, Infrastructure, Interfaces
✅ **Design Patterns** - Repository, Factory, Strategy, DI, Middleware
✅ **API Documentation** - Swagger UI
✅ **Validation** - go-playground/validator
✅ **Error Handling** - Структурированные ошибки
✅ **Configuration** - Environment variables

🚧 **In Progress** (заготовки созданы):
- Document Upload & Parsing (handlers готовы, парсеры - TODO)
- Test Generation with LLM (handlers готовы, интеграция - TODO)
- Moodle XML Export (логика готова, тестирование - TODO)

---

Готово! Бэкенд запущен и готов к использованию! 🎉
