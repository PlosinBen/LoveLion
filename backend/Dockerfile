# Production Dockerfile for Go backend
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy source and build
COPY . .
RUN go mod tidy && go build -o bin/server main.go

# Runtime image
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/bin/server .

# Install migrate tool for production migrations
RUN apk add --no-cache ca-certificates

EXPOSE 8080

CMD ["./server"]
