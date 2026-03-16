# Stage 1: Build
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Настройка прокси для прохождения таймаутов в КЗ
ENV GOPROXY=https://goproxy.io,direct

# Сначала копируем только файлы модулей для кэширования слоев
COPY go.mod go.sum ./
RUN go mod download || (echo "Сетевая ошибка, попробуй включить VPN или проверить go.mod" && exit 1)

# Копируем остальной код
COPY . .

# Собираем бинарник (добавлен флаг -ldflags для уменьшения размера)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main ./cmd/main.go

# Stage 2: Final
FROM alpine:latest

# Установка часового пояса и сертификатов (важно для работы с API и БД)
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Копируем только то, что нужно для работы
COPY --from=builder /app/main .
# Проверь, что папка migrations существует в корне проекта
COPY --from=builder /app/migrations ./migrations
# Если у тебя есть .env файл, его тоже нужно прокинуть (либо через docker-compose)

EXPOSE 8080

CMD ["./main"]