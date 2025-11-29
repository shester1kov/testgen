# Single-stage build using go run (no compilation issues)
FROM golang:1.25-bookworm

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Create uploads directory
RUN mkdir -p /app/uploads

EXPOSE 8080

# Run application directly with go run
CMD ["go", "run", "cmd/api/main.go"]
