# SoundMarket Backend

Backend для SoundMarket: пользователи регистрируются как `customer` или `engineer`, публикуют `Request` и `Offer`, отправляют `Bid`, создают `Order` и работают через внутренний баланс.

## Стек

- Go
- PostgreSQL
- Redis
- Chi
- JWT + bcrypt
- Ent schemas
- Docker / Docker Compose

## Что уже реализовано

- единый env loader с fail-fast проверкой;
- bootstrap приложения с проверкой подключения к PostgreSQL и Redis;
- автоматическое применение SQL-миграций при старте API в dev/docker;
- SQL repositories для `auth`, `profiles`, `cards`, `bids`, `orders`, `payments`, `notifications`;
- бизнес-правила для `Offer`, `Request`, `Bid`, `Order`;
- внутренний ledger через `transactions`;
- моковый YooKassa adapter;
- OpenAPI для рабочих endpoint-ов.

## Что пока scaffold

- реальная S3-интеграция загрузки файлов;
- FFmpeg / worker processing pipeline;
- полноценный WebSocket чат;
- reviews / disputes / moderation endpoints;
- Redis Pub/Sub слой для live-событий.

## Где лежат миграции

- SQL-миграции лежат в [migrations/001_init.sql](C:\Users\User\Desktop\project\migrations\001_init.sql)
- трекинг примененных миграций хранится в таблице `schema_migrations`

## Кто и когда запускает миграции

- В `docker compose`-сценарии миграции запускает сам `api` контейнер при старте приложения.
- В обычном локальном dev-режиме миграции тоже запускает приложение, если `AUTO_APPLY_MIGRATIONS=true`.
- В production это поведение можно отключить через `AUTO_APPLY_MIGRATIONS=false` и запускать миграции отдельным шагом деплоя.

## Почему раньше схема не создавалась

Раньше SQL лежал в `migrations/`, но `docker compose up --build` только поднимал `postgres`, `redis` и `api`.  
Ни `docker-compose.yml`, ни bootstrap приложения не выполняли `migrations/001_init.sql`, поэтому база была пустой и первый запрос к `/api/v1/cards` падал с `pq: relation "cards" does not exist`.

## Env

1. Скопировать `.env.example` в `.env`.
2. Заполнить значения без реальных боевых секретов.
3. Ключевые переменные:
   - `APP_PORT`
   - `AUTO_APPLY_MIGRATIONS`
   - `MIGRATIONS_DIR`
   - `ADMIN_BOOTSTRAP_ENABLED`
   - `ADMIN_BOOTSTRAP_EMAIL`
   - `ADMIN_BOOTSTRAP_PASSWORD`
   - `POSTGRES_*`
   - `REDIS_*`
   - `JWT_SECRET`
   - `JWT_TTL`
   - `S3_*`
   - `YOOKASSA_*`

`.env.example` соответствует текущему коду.

## Dev Admin

Self-register для `admin` намеренно запрещен.

Для локальной/dev проверки можно включить bootstrap admin через env:

```env
ADMIN_BOOTSTRAP_ENABLED=true
ADMIN_BOOTSTRAP_EMAIL=admin@example.com
ADMIN_BOOTSTRAP_PASSWORD=admin12345
```

При старте приложения в non-production окружении такой пользователь будет создан автоматически, если его еще нет. После этого можно получить JWT обычным `POST /api/v1/auth/login`.

## Запуск через Docker

После этого шага `/health` и `/api/v1/cards` должны работать сразу, без ручного `psql -f`.

```bash
docker compose up --build
```

Если нужно пересоздать БД с нуля:

```bash
docker compose down -v
docker compose up --build
```

## Запуск локально без Docker

1. Поднять PostgreSQL и Redis.
2. Убедиться, что в `.env` стоит `AUTO_APPLY_MIGRATIONS=true`.
3. Запустить:

```bash
go mod tidy
go run ./cmd/api
```

Если автоприменение выключено, тогда миграцию нужно выполнить вручную:

```bash
psql "host=localhost port=5432 dbname=soundmarket user=soundmarket password=soundmarket sslmode=disable" -f migrations/001_init.sql
```

## Полезные endpoint-ы

- `GET /health`
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `GET /api/v1/auth/me`
- `GET /api/v1/profiles/{id}`
- `GET /api/v1/profiles/me`
- `PUT /api/v1/profiles/me`
- `GET /api/v1/cards`
- `GET /api/v1/cards/{id}`
- `POST /api/v1/cards`
- `PUT /api/v1/cards/{id}`
- `GET /api/v1/requests/{id}/bids`
- `POST /api/v1/requests/{id}/bids`
- `POST /api/v1/orders/from-offer`
- `POST /api/v1/orders/from-bid`
- `GET /api/v1/orders/{id}`
- `PATCH /api/v1/orders/{id}/status`
- `POST /api/v1/payments/deposits`
- `POST /api/v1/payments/webhook`
- `GET /api/v1/payments/balance`

## Коммиты по этапам

1. `foundation: project structure and runtime config`
2. `infra: docker compose and init migration`
3. `auth: registration login and profiles`
4. `cards: offers requests and search`
5. `bids: request bids flow`
6. `orders: transactional order creation and status flow`
7. `payments: internal ledger and mock yookassa`
8. `docs: openapi and readme refresh`
