# =============================================================================
# Stage 1: Builder - собираем production build
# =============================================================================
FROM node:24.11.1-alpine AS builder

# Метаданные
LABEL maintainer="testgen-team"
LABEL description="TestGen Frontend Vue.js Application"

WORKDIR /app

# Копируем package files для кэширования зависимостей
# Docker кэширует этот слой, если package.json не изменился
COPY package*.json ./

# Устанавливаем зависимости
# npm ci - быстрее и надёжнее чем npm install для CI/CD
RUN npm ci --silent

# Копируем исходный код
COPY . .

# Build arguments для Vite environment variables
ARG VITE_API_BASE_URL=http://localhost:8080/api/v1
ARG VITE_MAX_FILE_SIZE=52428800
ARG VITE_SUPPORTED_FORMATS=pdf,docx,pptx,txt,md

# Устанавливаем переменные окружения для build
ENV VITE_API_BASE_URL=${VITE_API_BASE_URL}
ENV VITE_MAX_FILE_SIZE=${VITE_MAX_FILE_SIZE}
ENV VITE_SUPPORTED_FORMATS=${VITE_SUPPORTED_FORMATS}
ENV NODE_ENV=production

# Собираем production build с type checking
# Результат: статические файлы в dist/
# Type checking включен для гарантии качества кода при деплое
RUN npm run build

# =============================================================================
# Stage 2: Production - минимальный образ с Nginx
# =============================================================================
FROM nginx:1.27-alpine

# Метаданные
LABEL maintainer="testgen-team"
LABEL description="TestGen Frontend (Production)"

WORKDIR /usr/share/nginx/html

# Удаляем дефолтный nginx контент
RUN rm -rf ./*

# Копируем собранные статические файлы из builder stage
COPY --from=builder /app/dist .

# Копируем кастомную nginx конфигурацию для SPA
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Создаём непривилегированного пользователя для безопасности
# В production НИКОГДА не запускаем процессы от root
RUN addgroup -g 101 -S nginx || true && \
    adduser -u 101 -D -S -G nginx nginx || true

# Создаём необходимые директории и устанавливаем права
RUN chown -R nginx:nginx /usr/share/nginx/html && \
    chown -R nginx:nginx /var/cache/nginx && \
    chown -R nginx:nginx /var/log/nginx && \
    chown -R nginx:nginx /etc/nginx/conf.d && \
    touch /var/run/nginx.pid && \
    chown -R nginx:nginx /var/run/nginx.pid

# Переключаемся на непривилегированного пользователя
USER nginx

# Порт для Nginx
EXPOSE 80

# Health check (встроенный в образ, будет работать вместе с docker-compose)
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://127.0.0.1:80/health || exit 1

# Запускаем Nginx в foreground режиме
# daemon off - nginx остаётся в foreground для Docker
CMD ["nginx", "-g", "daemon off;"]
