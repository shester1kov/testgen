@echo off
echo.
echo ===================================
echo   Starting Backend Server
echo ===================================
echo.

:: Проверка Go
echo [1/7] Checking Go installation...
where go >nul 2>nul
if %errorlevel% neq 0 (
    echo ERROR: Go is not installed!
    exit /b 1
)
go version
echo.

:: Проверка Docker
echo [2/7] Checking Docker...
where docker >nul 2>nul
if %errorlevel% neq 0 (
    echo ERROR: Docker is not installed!
    exit /b 1
)
echo Docker is installed
echo.

:: Установка зависимостей
echo [3/7] Installing Go dependencies...
cd backend
go mod download
go mod tidy
cd ..
echo Dependencies installed
echo.

:: Проверка .env
echo [4/7] Checking .env file...
if not exist .env (
    echo .env file not found. Creating from .env.example...
    copy .env.example .env
)
echo .env file OK
echo.

:: Запуск PostgreSQL
echo [5/7] Starting PostgreSQL...
docker-compose up -d postgres
timeout /t 5 /nobreak >nul
echo PostgreSQL started
echo.

:: Миграции будут применены автоматически при запуске backend
echo [6/7] Database migrations will be applied automatically on backend start...
echo.

:: Генерация Swagger
echo [7/7] Generating Swagger documentation...
cd backend
if not exist docs (
    echo Installing swag...
    go install github.com/swaggo/swag/cmd/swag@latest
    swag init -g cmd/api/main.go -o docs
)
cd ..
echo Swagger documentation generated
echo.

:: Создание директории uploads
if not exist uploads mkdir uploads

echo.
echo ===================================
echo   Backend setup complete!
echo ===================================
echo.
echo Starting server...
echo.
echo Swagger UI will be available at:
echo http://localhost:8080/swagger/index.html
echo.

:: Запуск сервера
cd backend
go run cmd/api/main.go
