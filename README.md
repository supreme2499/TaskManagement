# tasks
___
API для управления заметками с поддержкой PostgreSQL, Redis и Kafka.
___
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
   cd tasks
   ```

2. Настройте .env файл(пример):
   ```
   ENV=local

   POSTGRES_STORAGE_URL=postgres://user:password@postgres:5432/database?sslmode=disable
   POSTGRES_MIGRATIONS_PATH=migrations

   REDIS_ADDRESS=redis:6379
   REDIS_DB=0

   HTTP_SERVER_ADDRESS=localhost:8000
   HTTP_SERVER_TIMEOUT=4s
   HTTP_SERVER_IDLE_TIMEOUT=60s
   HTTP_SERVER_WITH_TIMEOUT=10s
   
   KAFKA_ADDRESSES="kafka1:29091, kafka2:29092, kafka3:29093"
   ```
3. Запустите сервисы:
   ```
   docker-compose build
   docker-compose up -d
   ```
## Документация API

### Эндпоинты

## 1. Создать новую задачу
**POST** `/task`

**Параметры запроса**
- **Body**:
```json
{
  "task_text": "Название задачи",
  "description": "Описание задачи",
  "deadline": "2023-12-31T23:59:59Z"
}
```

**Ответ**
- Успешный ответ:
```json
{
  "response": {
    "status": "OK"
  },
  "task_id": 1
}
```

- Ошибка:
```json
{
  "response": {
    "status": "Error",
    "error": "описание ошибки"
  }
}
```

---

## 2. Прикрепить пользователя к задаче
**POST** `/adduser`

**Параметры запроса**
- **Body**:
```json
{
  "user_id": 1,
  "task_id": 1
}
```

**Ответ**
- Успешный ответ:
```json
{
  "response": {
    "status": "OK"
  }
}
```

- Ошибка:
```json
{
  "response": {
    "status": "Error",
    "error": "описание ошибки"
  }
}
```

---

## 3. Получить список всех пользователей работающих над задачей
**GET** `/users`

**Параметры запроса**
- **Body**:
```json
{
  "task_id": 1
}
```

**Ответ**
- Успешный ответ:
```json
{
  "users": [
    {
      "id": 1,
      "name": "Имя пользователя"
    }
  ],
  "response": {
    "status": "OK"
  }
}
```

- Ошибка:
```json
{
  "response": {
    "status": "Error",
    "error": "описание ошибки"
  }
}
```

---

## 4. Получить список всех задач над которыми работает конкретный пользователь
**GET** `/tasks`

**Параметры запроса**
- **Body**:
```json
{
  "user_id": 1
}
```

**Ответ**
- Успешный ответ:
```json
{
  "tasks": [
    {
      "id": 1,
      "name_task": "Название задачи",
      "description": "Описание задачи",
      "deadline": "2023-12-31T23:59:59Z",
      "status": "В процессе"
    }
  ],
  "response": {
    "status": "OK"
  }
}
```

- Ошибка:
```json
{
  "response": {
    "status": "Error",
    "error": "описание ошибки"
  }
}
```

---

## 5. Получить задачи с коротким дедлайном
**GET** `/shortdeadline`

**Параметры запроса**
- **Body**:
```json
{
  "user_id": 1
}
```

**Ответ**
- Успешный ответ:
```json
{
  "tasks": [
    {
      "id": 1,
      "name_task": "Название задачи",
      "description": "Описание задачи",
      "deadline": "2023-12-31T23:59:59Z",
      "status": "В процессе"
    }
  ],
  "response": {
    "status": "OK"
  }
}
```

- Ошибка:
```json
{
  "response": {
    "status": "Error",
    "error": "описание ошибки"
  }
}
```

---

## 6. Получить задачу по её ID
**GET** `/taskbyid`

**Параметры запроса**
- **Body**:
```json
{
  "task_id": 1
}
```

**Ответ**
- Успешный ответ:
```json
{
  "task": {
    "id": 1,
    "name_task": "Название задачи",
    "description": "Описание задачи",
    "deadline": "2023-12-31T23:59:59Z",
    "status": "В процессе"
  },
  "response": {
    "status": "OK"
  }
}
```

- Ошибка:
```json
{
  "response": {
    "status": "Error",
    "error": "описание ошибки"
  }
}
```

---

## 7. Обновить статус задачи
**PUT** `/status`

**Параметры запроса**
- **Body**:
```json
{
  "task_id": 1,
  "new_status": "Завершено"
}
```

**Ответ**
- Успешный ответ:
```json
{
  "response": {
    "status": "OK"
  }
}
```

- Ошибка:
```json
{
  "response": {
    "status": "Error",
    "error": "описание ошибки"
  }
}
```

---

## 8. Удалить задачу
**DELETE** `/task`

**Параметры запроса**
- **Body**:
```json
{
  "task_id": 1
}
```

**Ответ**
- Успешный ответ:
```json
{
  "response": {
    "status": "OK"
  }
}
```

- Ошибка:
```json
{
  "response": {
    "status": "Error",
    "error": "описание ошибки"
  }
}
```

---

## 9. Снять пользователя с задачи
**DELETE** `/user`

**Параметры запроса**
- **Body**:
```json
{
  "user_id": 1,
  "task_id": 1
}
```

**Ответ**
- Успешный ответ:
```json
{
  "response": {
    "status": "OK"
  }
}
```

- Ошибка:
```json
{
  "response": {
    "status": "Error",
    "error": "описание ошибки"
  }
}
```

---



   