# Используем официальный образ Go
FROM golang:1.23.2 AS builder

# Устанавливаем зависимости для библиотеки librdkafka
RUN apt-get update && \
    apt-get install -y \
    build-essential \
    librdkafka-dev \
    git

# Устанавливаем рабочую директорию в контейнере
WORKDIR /app

# Копируем файл go.mod и go.sum в контейнер
COPY go.mod go.sum ./

# Загружаем все зависимости
RUN go mod tidy

# Копируем весь исходный код в контейнер
COPY . .

# Переходим в директорию с main.go
WORKDIR /app/cmd/app

# Компилируем приложение для Linux архитектуры
RUN GOOS=linux GOARCH=amd64 go build -o /app/myapp .

# Используем более новый образ для выполнения приложения
FROM debian:bookworm-slim

# Устанавливаем необходимые библиотеки для работы с Kafka
RUN apt-get update && \
    apt-get install -y \
    librdkafka1

# Копируем скомпилированный файл из предыдущего этапа
COPY --from=builder /app/myapp /myapp

# Открываем порт, если необходимо
EXPOSE 8000

# Устанавливаем команду для запуска приложения
CMD ["/myapp"]
