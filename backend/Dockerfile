# =============================================================================
# Stage 1: Builder - компилируем приложение
# =============================================================================
FROM golang:1.25-bookworm AS builder

WORKDIR /app

# Копируем go.mod и go.sum для кэширования зависимостей
# Docker кэширует этот слой, если файлы не изменились
COPY go.mod go.sum ./

# Скачиваем зависимости (кэшируется отдельно от кода)
RUN go mod download && go mod verify

RUN go install github.com/google/wire/cmd/wire@latest

# Копируем весь исходный код
COPY . .

RUN cd cmd/api && wire

RUN ls -la cmd/api/wire_gen.go

# Компилируем бинарник
# CGO_ENABLED=0 - статическая компиляция без зависимостей от libc
# -ldflags="-w -s" - уменьшаем размер бинарника (удаляем debug info и symbol table)
# -trimpath - удаляем абсолютные пути из бинарника
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -trimpath \
    -o /app/bin/api \
    ./cmd/api

# =============================================================================
# Stage 2: Production - минимальный образ для запуска
# =============================================================================
FROM debian:bookworm-slim

# Метаданные образа
LABEL maintainer="testgen-team"
LABEL description="TestGen Backend API"

WORKDIR /app

# Устанавливаем минимальный набор утилит:
# - wget для healthcheck
# - ca-certificates для HTTPS запросов (LLM API, Moodle)
# - tzdata для правильной работы с часовыми поясами
RUN apt-get update && apt-get install -y --no-install-recommends \
    wget \
    ca-certificates \
    tzdata \
    && rm -rf /var/lib/apt/lists/*

# Создаём непривилегированного пользователя для безопасности
RUN useradd -m -u 1000 -s /bin/bash appuser

# Копируем скомпилированный бинарник из builder stage
COPY --from=builder /app/bin/api /app/api

# Копируем миграции БД (нужны для автоматического запуска миграций)
COPY --from=builder /app/internal/infrastructure/persistence/migrations /app/internal/infrastructure/persistence/migrations

# Переключаемся на непривилегированного пользователя
USER appuser

# Порт приложения
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Запускаем скомпилированный бинарник
CMD ["/app/api"]
