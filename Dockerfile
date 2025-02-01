# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install required packages
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -o main ./cmd/api

# Final stage
FROM alpine:latest

WORKDIR /app

# Install required runtime packages
RUN apk --no-cache add ca-certificates tzdata

# Copy the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

# Create non-root user
RUN adduser -D -g '' appuser
USER appuser

# Command to run the application
ENTRYPOINT ["./main"]