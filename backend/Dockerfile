# Production Dockerfile for Go backend
FROM golang:1.25.9-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bin/server main.go && \
    go build -o bin/migrate cmd/migrate/main.go

# Runtime image
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/bin/server .
COPY --from=builder /app/bin/migrate .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./server"]
