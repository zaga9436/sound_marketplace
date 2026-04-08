# SoundMarket Backend

Backend для платформы SoundMarket, которая связывает заказчиков и звукоинженеров: публикация `Offer` и `Request`, сделки через внутренний баланс, заказы, чат, уведомления и работа с медиа.

## Стек

- Go
- REST API + WebSocket
- Ent ORM
- PostgreSQL
- Redis
- S3-compatible storage
- FFmpeg worker
- Docker / Docker Compose

## Что уже заложено

- модульный монолит с разделением по слоям;
- конфигурация через `.env`;
- каркас auth, profiles, cards, bids, orders, payments;
- Ent-схемы основных сущностей;
- OpenAPI-описание первого этапа;
- задел под storage, worker, notifications и WebSocket-чат.

## Локальный запуск

1. Скопировать `.env.example` в `.env`.
2. Заполнить значения переменных окружения.
3. Убедиться, что установлены Go и Docker.
4. Запустить `docker compose up --build`.

## Полезные команды

```bash
go run ./cmd/api
go test ./...
docker compose up --build
```

## Сервисы

- `api` — HTTP/WebSocket backend
- `postgres` — основная база данных
- `redis` — кэш, Pub/Sub, вспомогательная инфраструктура

## Переменные окружения

Все чувствительные настройки должны храниться только в `.env`. Для старта обязательны:

- `APP_PORT`
- `POSTGRES_*`
- `REDIS_*`
- `JWT_SECRET`
- `JWT_TTL`
- `S3_*`
- `YOOKASSA_*`

## Git-подход

Изменения стоит коммитить отдельными этапами:

1. foundation и структура проекта
2. docker-compose и config
3. auth и users/profiles
4. cards
5. bids
6. orders
7. payments
8. chat/websocket
9. storage/worker
10. reviews/disputes/moderation
