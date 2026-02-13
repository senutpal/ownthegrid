# backend

go 1.22 api + websocket server for ownthegrid. exposes rest endpoints for board/user data and a websocket gateway for real-time tile claims.

## services

- http api: `:8080`
- websocket: `/ws`
- postgres 18
- redis 8

## env

copy `backend/.env.example` to `backend/.env` and set values.

required:

- `DATABASE_URL`
- `REDIS_URL`
- `JWT_SECRET` (min 32 chars)

## run locally

1. `docker-compose up postgres redis -d`
2. `psql $DATABASE_URL -f migrations/001_init.up.sql`
3. `go run cmd/server/main.go`

## routes

rest:

- `POST /api/users/register`
- `GET /api/users/{id}`
- `GET /api/users/online`
- `GET /api/users/leaderboard`
- `GET /api/board`
- `GET /api/board/stats`

websocket:

- `GET /ws?userId=<id>&token=<jwt>`

## messaging rules

- ws messages are always `{ "type": "...", "payload": { ... } }`
- db writes never happen in websocket handlers
- broadcasts always go through redis pub/sub (`board:events`)
- first-write-wins enforced in sql (`WHERE owner_id IS NULL`)

## build

`go test ./...`
