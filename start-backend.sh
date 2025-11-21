#!/bin/bash

# Цвета для вывода
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}🚀 Starting Test Generation System Backend${NC}"
echo ""

# Проверка Go
echo -e "${YELLOW}[1/7]${NC} Checking Go installation..."
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed. Please install Go 1.21+${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Go $(go version | awk '{print $3}')${NC}"
echo ""

# Проверка Docker
echo -e "${YELLOW}[2/7]${NC} Checking Docker..."
if ! command -v docker &> /dev/null; then
    echo -e "${RED}❌ Docker is not installed${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Docker is installed${NC}"
echo ""

# Установка зависимостей
echo -e "${YELLOW}[3/7]${NC} Installing Go dependencies..."
cd backend
go mod download
go mod tidy
echo -e "${GREEN}✅ Dependencies installed${NC}"
echo ""

# Проверка .env
echo -e "${YELLOW}[4/7]${NC} Checking .env file..."
cd ..
if [ ! -f .env ]; then
    echo -e "${YELLOW}⚠️  .env file not found. Creating from .env.example...${NC}"
    if [ -f .env.example ]; then
        cp .env.example .env
        echo -e "${GREEN}✅ .env file created${NC}"
    else
        echo -e "${RED}❌ .env.example not found${NC}"
        exit 1
    fi
else
    echo -e "${GREEN}✅ .env file exists${NC}"
fi
echo ""

# Запуск PostgreSQL
echo -e "${YELLOW}[5/7]${NC} Starting PostgreSQL..."
docker-compose up -d postgres
sleep 3
echo -e "${GREEN}✅ PostgreSQL started${NC}"
echo ""

# Применение миграций
echo -e "${YELLOW}[6/7]${NC} Applying database migrations..."
echo "Waiting for PostgreSQL to be ready..."
sleep 5

# Попытка применить миграции
if command -v migrate &> /dev/null; then
    cd backend
    migrate -path internal/infrastructure/persistence/migrations \
      -database "postgres://testgen_user:testgen_password@localhost:5432/testgen_db?sslmode=disable" up
    cd ..
    echo -e "${GREEN}✅ Migrations applied${NC}"
else
    echo -e "${YELLOW}⚠️  migrate not found. Using Docker...${NC}"
    docker run --rm \
      -v "$(pwd)/backend:/app" \
      --network host \
      migrate/migrate \
      -path /app/internal/infrastructure/persistence/migrations \
      -database "postgres://testgen_user:testgen_password@localhost:5432/testgen_db?sslmode=disable" up
    echo -e "${GREEN}✅ Migrations applied${NC}"
fi
echo ""

# Генерация Swagger
echo -e "${YELLOW}[7/7]${NC} Generating Swagger documentation..."
cd backend
if [ ! -d "docs" ]; then
    if command -v swag &> /dev/null; then
        swag init -g cmd/api/main.go -o docs
    else
        echo -e "${YELLOW}⚠️  Installing swag...${NC}"
        go install github.com/swaggo/swag/cmd/swag@latest
        export PATH=$PATH:$(go env GOPATH)/bin
        swag init -g cmd/api/main.go -o docs
    fi
    echo -e "${GREEN}✅ Swagger documentation generated${NC}"
else
    echo -e "${GREEN}✅ Swagger docs already exist${NC}"
fi
cd ..
echo ""

# Создание директории uploads
mkdir -p uploads

echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}✅ Backend setup complete!${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${YELLOW}Starting server...${NC}"
echo ""

# Запуск сервера
cd backend
go run cmd/api/main.go
