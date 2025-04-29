# 🛒 Purchase Service

Сервис для управления пользовательскими покупками. Поддерживает регистрацию запланированных покупок, их обновление, удаление, статистику, напоминания и метки.

## 📦 Основной функционал:
- Создание покупок
- Получение всех покупок пользователя
- Обновление информации о покупке
- Пометка покупки как завершённой
- Деактивация покупки (мягкое удаление)
- Добавление напоминания
- Привязка меток (тегов)
- Получение общей статистики

---

## 📘 API Роуты

### `POST /purchases`
Создать список новых покупок

**Headers:**  
`X-User-ID: <uuid>`

**Body:**
```json
[
  {
    "name": "Хлеб",
    "category": "Продукты",
    "estimated_price": 80,
    "deadline": "2025-05-10T12:00:00Z"
  },
  {
    "name": "Мыло",
    "category": "Бытовая химия"
  }
]
```

### `GET /purchases`
Получить все покупки пользователя

**Headers:**
`X-User-ID: <uuid>`

**Response:** `200 OK`

```
[
  {
    "id": "UUID",
    "name": "Хлеб",
    "category": "Продукты",
    "estimated_price": 80.0,
    "actual_price": null,
    "created_at": "2025-04-01T12:00:00Z",
    "deadline": "2025-05-10T12:00:00Z",
    "is_active": true,
    "is_purchased": false
  }
]
```

### `PATCH /purchases/{id}/purchase`
Отметить покупку как завершённую

**Headers:**
`X-User-ID: <uuid>`

**Response:** `200 OK`
```
{
  "status": "marked as purchased"
}
```

### `PATCH /purchases/{id}/deactivate`
Деактивировать (мягко удалить) покупку

**Headers:**
`X-User-ID: <uuid>`

**Response:** `200 OK`
```
{
  "status": "purchase deactivated"
}
```

### `PUT /purchases/{id}`
Обновить поля покупки (только переданные будут обновлены)

**Headers:**
`X-User-ID: <uuid>`

**Body:**
```
{
  "name": "Хлеб цельнозерновой",
  "estimated_price": 120
}
```
**Response:** `200 OK`
```
{
  "status": "purchase updated"
}
```

### `POST /purchases/{id}/reminder`
Добавить напоминание о покупке

**Headers:**
`X-User-ID: <uuid>`

**Body:**
```
{
  "remind_at": "2025-05-08T08:00:00Z",
  "message": "Купить до завтра!"
}
```

**Response:** `201 Created`

```
{
  "status": "reminder added"
}

```

### `POST /purchases/{id}/tags`
Привязать метки к покупке

**Headers:**
`X-User-ID: <uuid>`

**Body:**
```
{
  "tag_ids": ["tag-1", "tag-2"]
}
```
**Response:** `201 Created`
```
{
  "status": "tags added"
}
```

### `GET /purchases/stats`
Получить агрегированную статистику покупок

**Headers:**
`X-User-ID: <uuid>`

**Response:** `200 OK`

```
{
  "total": 12,
  "active": 7,
  "purchased": 5,
  "estimated_total": 2500.00,
  "actual_total": 1900.00
}
```

## 🔒 Аутентификация
- На текущий момент авторизация выполняется через заголовок X-User-ID

- В будущем планируется интеграция с централизованным auth-сервисом (через JWT)

## 🛠 Используемые технологии
- Golang net/http
- PostgreSQL (через pgx)
- Redis (опционально, для напоминаний/кеша)
- Jaeger (трейсинг через OTLP)
- Elasticsearch + Filebeat (сбор и анализ логов)
- slog — structured logging
