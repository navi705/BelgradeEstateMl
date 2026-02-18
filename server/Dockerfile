# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod ./
# Copy sum file if it exists, otherwise ignore error
COPY go.su[m] ./ 
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/main .

EXPOSE 8000

ENV PORT=8000
ENV ML_ENGINE_URL=http://localhost:8080

CMD ["./main"]
