# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Runtime stage
FROM alpine:latest

# Install CA certificates for secure connections
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary and assets from builder
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Expose API port
EXPOSE 8080

# Environment variables with defaults
ENV PORT=8080

CMD ["./main"]
