# frontend

react 19 + typescript + tailwind + motion client for ownthegrid. renders the board, handles optimistic tile claims, and syncs state over websocket.

## env

copy `frontend/.env.example` to `frontend/.env`.

required:

- `VITE_API_URL`
- `VITE_WS_URL`
- `VITE_GRID_WIDTH`
- `VITE_GRID_HEIGHT`

## run locally

1. `pnpm install`
2. `pnpm dev`

## build

`pnpm build`

## data flow

- initial board load via `GET /api/board`
- websocket connects to `/ws` and receives `INIT_BOARD`
- tile claims: optimistic ui -> `CLAIM_TILE` -> `TILE_CLAIMED` or `CLAIM_REJECTED`

## state stores

- `boardStore`: tiles, optimistic updates, grid sizing
- `userStore`: current user, online users, leaderboard
- `wsStore`: connection state and send function
- `uiStore`: toasts and onboarding modal

## styling

- tailwind for layout and typography
- css modules for grid/tile rendering
