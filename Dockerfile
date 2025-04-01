# Стадия сборки
FROM golang:alpine AS builder

# Установка зависимостей
RUN apk add --no-cache git

# Рабочая директория
WORKDIR /app

# Копирование go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения из поддиректории cmd/event_service
RUN CGO_ENABLED=0 GOOS=linux go build -o event_service ./cmd/event_service

# Финальный образ
FROM alpine:latest

# Установка tzdata для работы с временными зонами
RUN apk add --no-cache tzdata

# Рабочая директория
WORKDIR /app

# Копирование бинарника из стадии сборки
COPY --from=builder /app/event_service .
COPY --from=builder /app/config .

# Команда запуска
CMD ["./event_service"]