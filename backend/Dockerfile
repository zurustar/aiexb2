# Backend Dockerfile
# Go development environment with hot reload

FROM golang:1.21-alpine

WORKDIR /app

# Install Air for hot reload
RUN go install github.com/cosmtrek/air@latest

# Install migrate tool
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build
RUN go build -o /app/api cmd/api/main.go
RUN go build -o /app/worker cmd/worker/main.go

# Expose port
EXPOSE 8080

# Start with Air for hot reload
CMD ["air", "-c", ".air.toml"]
