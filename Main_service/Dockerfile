# Stage 1: Build
FROM golang:1.24-alpine AS builder

ENV GOTOOLCHAIN=local

RUN apk add --no-cache tzdata ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
COPY vendor/ ./vendor/
COPY . .

# Собираем бинарник (добавлен флаг -ldflags для уменьшения размера)
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -ldflags="-s -w" -o main ./cmd/main.go

# Stage 2: Final
FROM alpine:3.21

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./main"]