# SoundMarket Backend

Backend for SoundMarket marketplace. Users register as `customer` or `engineer`, publish `request` and `offer` cards, send bids, create orders, exchange messages, upload media and deliverables, resolve disputes, leave verified reviews, and work through an internal balance with external payment provider integration.

## Stack

- Go
- PostgreSQL
- Redis
- Chi
- JWT + bcrypt
- Docker / Docker Compose
- S3-compatible object storage
- YooKassa provider adapter

## Current Scope

Implemented backend flows:

- auth: register, login, me
- profiles: public profile, self update, public cards, public reviews
- cards: create, list, get, update, search, filters, sorting, pagination
- bids: create, list
- orders: create from offer or bid, list, get, lifecycle transitions
- ledger/payments: deposits, holds, releases, refunds, YooKassa create payment, webhook, dev-only sync
- disputes: open, get, close with resolution
- reviews: verified review after completed order, profile rating aggregation
- media: preview upload, full upload, signed download
- deliverables: upload, list, signed download, versioning
- chat: conversations, messages, mark as read, websocket endpoints
- notifications: list, mark read, websocket endpoint
- admin: users, cards, disputes, moderation actions

Not in focus for this pass:

- highload tuning
- analytics
- observability stack
- frontend

## Project Layout

- [cmd/api/main.go](C:\Users\User\Desktop\project\cmd\api\main.go): application entrypoint
- [internal/app/app.go](C:\Users\User\Desktop\project\internal\app\app.go): bootstrap and wiring
- [internal/http/router/router.go](C:\Users\User\Desktop\project\internal\http\router\router.go): route registration
- [internal/service](C:\Users\User\Desktop\project\internal\service): business logic
- [internal/repository](C:\Users\User\Desktop\project\internal\repository): PostgreSQL repositories
- [internal/storage](C:\Users\User\Desktop\project\internal\storage): S3-compatible storage adapter
- [migrations](C:\Users\User\Desktop\project\migrations): SQL migrations
- [docs/openapi.yaml](C:\Users\User\Desktop\project\docs\openapi.yaml): API contract

## Environment

1. Copy `.env.example` to `.env`.
2. Fill values for PostgreSQL, Redis, S3, JWT and payment provider.
3. Do not use production secrets in local development.

Important env groups:

- app:
  - `APP_ENV`
  - `APP_PORT`
  - `AUTO_APPLY_MIGRATIONS`
  - `MIGRATIONS_DIR`
- admin bootstrap:
  - `ADMIN_BOOTSTRAP_ENABLED`
  - `ADMIN_BOOTSTRAP_EMAIL`
  - `ADMIN_BOOTSTRAP_PASSWORD`
- postgres:
  - `POSTGRES_HOST`
  - `POSTGRES_PORT`
  - `POSTGRES_DB`
  - `POSTGRES_USER`
  - `POSTGRES_PASSWORD`
  - `POSTGRES_SSLMODE`
- redis:
  - `REDIS_HOST`
  - `REDIS_PORT`
  - `REDIS_PASSWORD`
- auth:
  - `JWT_SECRET`
  - `JWT_TTL`
- payments:
  - `PAYMENT_PROVIDER=mock|yookassa`
  - `YOOKASSA_SHOP_ID`
  - `YOOKASSA_SECRET_KEY`
  - `YOOKASSA_RETURN_URL`
- storage:
  - `S3_ENDPOINT`
  - `S3_REGION`
  - `S3_BUCKET`
  - `S3_ACCESS_KEY`
  - `S3_SECRET_KEY`
  - `S3_USE_SSL`
  - `S3_FORCE_PATH_STYLE`
  - `SIGNED_URL_TTL`

## Local Run With Docker

Start everything:

```bash
docker compose up --build
```

Reset database volumes:

```bash
docker compose down -v
docker compose up --build
```

Expected result after startup:

- PostgreSQL is healthy
- Redis is healthy
- API applies migrations automatically in dev/docker
- `GET /health` returns `200 OK`

## Local Run Without Docker

1. Start PostgreSQL and Redis.
2. Keep `AUTO_APPLY_MIGRATIONS=true` in `.env`.
3. Run:

```bash
go mod tidy
go run ./cmd/api
```

## Migrations

Migrations live in [migrations](C:\Users\User\Desktop\project\migrations).

Current behavior:

- in docker/dev, the API applies migrations on startup when `AUTO_APPLY_MIGRATIONS=true`
- in production, migrations can be disabled in app startup and run as a separate deploy step
- applied migrations are tracked in `schema_migrations`

## Dev-only Flows

### Bootstrap Admin

Self-register for `admin` stays disabled.

For local development:

```env
ADMIN_BOOTSTRAP_ENABLED=true
ADMIN_BOOTSTRAP_EMAIL=admin@example.com
ADMIN_BOOTSTRAP_PASSWORD=admin12345
```

When enabled outside production, admin is created on startup if missing. After that you can use normal `POST /api/v1/auth/login`.

### Manual Payment Sync

`POST /api/v1/payments/sync` is a dev-only fallback for local YooKassa testing without a public webhook URL.

Typical flow:

1. create payment through `POST /api/v1/payments/deposits`
2. pay it in YooKassa test mode
3. call `POST /api/v1/payments/sync`
4. verify internal balance through `GET /api/v1/payments/balance`

In production, the main path remains payment provider webhook processing.

## Frontend-ready API Areas

Stable list/detail flows already prepared for frontend:

- auth:
  - `POST /api/v1/auth/register`
  - `POST /api/v1/auth/login`
  - `GET /api/v1/auth/me`
- profile:
  - `GET /api/v1/profiles/{id}`
  - `GET /api/v1/profiles/{id}/cards`
  - `GET /api/v1/profiles/{id}/reviews`
  - `GET /api/v1/profiles/me`
  - `PUT /api/v1/profiles/me`
- catalog:
  - `GET /api/v1/cards`
  - `GET /api/v1/cards/{id}`
  - `POST /api/v1/cards`
  - `PUT /api/v1/cards/{id}`
- bids:
  - `GET /api/v1/requests/{id}/bids`
  - `POST /api/v1/requests/{id}/bids`
- orders:
  - `GET /api/v1/orders`
  - `GET /api/v1/orders/{id}`
  - `POST /api/v1/orders/from-offer`
  - `POST /api/v1/orders/from-bid`
  - `PATCH /api/v1/orders/{id}/status`
- disputes:
  - `POST /api/v1/orders/{id}/dispute`
  - `GET /api/v1/orders/{id}/dispute`
  - `POST /api/v1/orders/{id}/dispute/close`
- reviews:
  - `POST /api/v1/orders/{id}/reviews`
- media:
  - `POST /api/v1/cards/{id}/media/preview`
  - `POST /api/v1/cards/{id}/media/full`
  - `GET /api/v1/cards/{id}/download`
- deliverables:
  - `POST /api/v1/orders/{id}/deliverables`
  - `GET /api/v1/orders/{id}/deliverables`
  - `GET /api/v1/orders/{id}/deliverables/{deliverable_id}/download`
- chat:
  - `GET /api/v1/chats`
  - `GET /api/v1/orders/{id}/messages`
  - `POST /api/v1/orders/{id}/messages`
  - `POST /api/v1/orders/{id}/messages/read`
- notifications:
  - `GET /api/v1/notifications`
  - `POST /api/v1/notifications/read`
- payments:
  - `POST /api/v1/payments/deposits`
  - `POST /api/v1/payments/webhook`
  - `POST /api/v1/payments/sync`
  - `GET /api/v1/payments/balance`
- admin:
  - `GET /api/v1/admin/users`
  - `GET /api/v1/admin/cards`
  - `GET /api/v1/admin/disputes`
  - moderation actions endpoints

## Tests

Run the current smoke-level backend tests:

```bash
go test ./...
```

The stabilization pass includes service-level tests for:

- order lifecycle
- dispute close flow
- payment sync idempotency
- deliverables versioning and access
- admin suspend/hide enforcement

## OpenAPI

The current API contract is in [openapi.yaml](C:\Users\User\Desktop\project\docs\openapi.yaml).

Frontend work should use that file as the main source of truth for:

- request shapes
- list pagination responses
- auth expectations
- admin endpoints
