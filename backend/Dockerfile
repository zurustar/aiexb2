# Backend Dockerfile
# Go development environment with hot reload

FROM golang:1.24-alpine

WORKDIR /app

# Install dependencies required for go install
RUN apk add --no-cache git build-base

# Configure Go proxy
ENV GOPROXY=https://proxy.golang.org,direct

# Install Air for hot reload (v1.61.1 is the last version compatible with Go 1.23)
RUN go install github.com/air-verse/air@v1.61.1

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
