# Multi-stage build for optimized Go binary
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with CGO enabled for MySQL driver compatibility
RUN CGO_ENABLED=1 GOOS=linux go build -o main cmd/server/main.go

# Final stage - minimal runtime image
FROM alpine:latest 

# Install runtime dependencies
RUN apk --no-cache add ca-certificates mysql-client bash curl

WORKDIR /root/

# Copy binary and files from builder
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/.env .
COPY --from=builder /app/wait-for-service.sh .

# Make the wait script executable
RUN chmod +x wait-for-service.sh

# Expose port
EXPOSE 8085

# Run the application with MySQL wait
CMD ["./wait-for-service.sh", "mysql", "3306", "./main"]