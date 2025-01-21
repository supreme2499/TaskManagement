# TaskManagement

API для управления заметками с поддержкой PostgreSQL, Redis и Kafka.

## Возможности
- Создание и удаление заметок.
- Прикрепление/снятие пользователя к задаче.
- Получение задачи по её ID.
- Получение списка задач по userID.
- Изменение статуса задачи.
- Уведомление через Kafka при изменении статуса задачи.

## Технологии
- **Backend**: Go
- **Database**: PostgreSQL
- **Cache**: Redis
- **Message Broker**: Kafka
- **Докеризация**: Docker Compose

## Установка
1. Клонируйте репозиторий:
   ```
   https://github.com/supreme2499/TaskManagement.git
   cd notes-management-api
   ```

2. Настройте .env файл(пример):
   ```
   ENV=local

   POSTGRES_STORAGE_URL=postgres://user:password@address:port/database?sslmode=disable
   POSTGRES_MIGRATIONS_PATH=migrations

   REDIS_ADDRESS=localhost:6379
   REDIS_DB=0

   HTTP_SERVER_ADDRESS=localhost:8080
   HTTP_SERVER_TIMEOUT=4s
   HTTP_SERVER_IDLE_TIMEOUT=60s
   HTTP_SERVER_WITH_TIMEOUT=10s
   
   KAFKA_ADDRESSES="localhost:9091, localhost:9092, localhost:9093"
   ```
3. Запустите сервисы:
   ```
   docker-compose build
   docker-compose up -d

   go run ./cmd/migrator/main.go
   go run ./cmd/app/main.go 
   ```



   