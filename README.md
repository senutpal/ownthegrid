# ownthegrid

real-time collaborative grid where users claim tiles on a shared 50x40 board. updates are delivered over websockets, persisted in postgres, and broadcast through redis pub/sub.

## requirements

- go 1.22+
- node 20+
- docker + docker compose
- postgres 18, redis 8 (via compose)

## quick start

1. `docker-compose up postgres redis -d`
2. `psql $DATABASE_URL -f backend/migrations/001_init.up.sql`
3. `cd backend && cp .env.example .env`
4. `cd frontend && cp .env.example .env && pnpm install && pnpm dev`
5. (in another terminal) `cd backend && go run cmd/server/main.go`

## full docker

`docker-compose up --build`

frontend will be available on `http://localhost:5173`, backend on `http://localhost:8080`.

## repo layout

```
backend/   go server, websocket hub, redis pub/sub, postgres
frontend/  react + typescript + tailwind + motion
```

## websocket contract

all messages use the envelope `{ "type": "...", "payload": { ... } }`.

client -> server

- `CLAIM_TILE` `{ tileId: number }`
- `PING` `{}`

server -> client

- `INIT_BOARD`
- `TILE_CLAIMED`
- `CLAIM_REJECTED`
- `USER_JOINED`
- `USER_LEFT`
- `LEADERBOARD_UPDATE`
- `ERROR`
- `PONG`

## env vars

see `backend/.env.example` and `frontend/.env.example`.

## testing

- backend: `go test ./...`
- frontend: `pnpm lint && pnpm build`
