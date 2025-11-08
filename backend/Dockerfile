FROM golang:1.21-alpine AS development

WORKDIR /app

# Install development tools
RUN apk add --no-cache git make curl

# Install goose and sqlc
RUN go install github.com/pressly/goose/v3/cmd/goose@latest && \
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Expose port
EXPOSE 8080

# Default command (can be overridden in docker-compose)
CMD ["go", "run", "cmd/server/main.go"]

# Production build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server cmd/server/main.go

# Final production stage
FROM alpine:latest AS production

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/server .
COPY --from=builder /app/config ./config

# Create uploads directory
RUN mkdir -p /app/uploads

EXPOSE 8080

CMD ["./server"]
