import { create } from 'zustand';
import { immer } from 'zustand/middleware/immer';
import { enableMapSet } from "immer";

import type { Tile, TileMap } from '../types/tile';
import type { TileClaimedPayload } from '../types/ws';

enableMapSet();

interface BoardState {
  tiles: TileMap;
  optimisticTiles: Map<number, string>;
  gridWidth: number;
  gridHeight: number;
  isLoaded: boolean;
  previousTiles: Map<number, Tile>;
  pendingClaims: Set<number>;
  claimedTiles: number;
  optimisticClaimedDelta: number;
  lastActivity: string | null;
  initBoard: (tiles: Tile[], gridWidth: number, gridHeight: number) => void;
  applyTileClaimed: (payload: TileClaimedPayload) => void;
  optimisticallyClaimTile: (tileId: number, userId: string, color: string) => void;
  revertOptimisticClaim: (tileId: number) => void;
  getClaimedTiles: () => number;
  getUnclaimedTiles: () => number;
  getLastActivity: () => string | null;
}

export const useBoardStore = create<BoardState>()(
  immer((set, get) => ({
    tiles: new Map(),
    optimisticTiles: new Map(),
    gridWidth: Number(import.meta.env.VITE_GRID_WIDTH ?? 50),
    gridHeight: Number(import.meta.env.VITE_GRID_HEIGHT ?? 40),
    isLoaded: false,
    previousTiles: new Map(),
    pendingClaims: new Set(),
    claimedTiles: 0,
    optimisticClaimedDelta: 0,
    lastActivity: null,

    initBoard: (tiles, gridWidth, gridHeight) =>
      set((state) => {
        state.tiles = new Map(tiles.map((tile) => [tile.id, tile]));
        state.gridWidth = gridWidth;
        state.gridHeight = gridHeight;
        state.isLoaded = true;
        state.claimedTiles = tiles.filter((t) => t.ownerId !== null && t.ownerId !== undefined).length;
        const tilesWithClaimedAt = tiles.filter((t) => t.claimedAt);
        if (tilesWithClaimedAt.length > 0) {
          const sorted = tilesWithClaimedAt.sort(
            (a, b) => new Date(b.claimedAt!).getTime() - new Date(a.claimedAt!).getTime()
          );
          state.lastActivity = sorted[0].claimedAt;
        }
      }),

    applyTileClaimed: (payload) =>
      set((state) => {
        const existing = state.tiles.get(payload.tileId);
        if (existing) {
          if (!existing.ownerId) {
            state.claimedTiles += 1;
          }
          existing.ownerId = payload.userId;
          existing.ownerUsername = payload.username;
          existing.ownerColor = payload.color;
          existing.claimedAt = payload.claimedAt;
          state.lastActivity = payload.claimedAt;
        }
        state.optimisticTiles.delete(payload.tileId);
        state.optimisticClaimedDelta = 0;
        state.previousTiles.delete(payload.tileId);
        const newPending = new Set(state.pendingClaims);
        newPending.delete(payload.tileId);
        state.pendingClaims = newPending;
      }),

    optimisticallyClaimTile: (tileId, userId, color) =>
      set((state) => {
        const newPending = new Set(state.pendingClaims);
        newPending.add(tileId);
        state.pendingClaims = newPending;
        if (!state.previousTiles.has(tileId)) {
          const snapshot = state.tiles.get(tileId);
          if (snapshot) {
            state.previousTiles.set(tileId, { ...snapshot });
          }
        }
        state.optimisticTiles.set(tileId, userId);
        const tile = state.tiles.get(tileId);
        if (tile) {
          if (!tile.ownerId) {
            state.optimisticClaimedDelta += 1;
          }
          tile.ownerId = userId;
          tile.ownerColor = color;
          tile.claimedAt = new Date().toISOString();
          state.lastActivity = tile.claimedAt;
        }
      }),

    revertOptimisticClaim: (tileId) =>
      set((state) => {
        state.optimisticTiles.delete(tileId);
        const newPending = new Set(state.pendingClaims);
        newPending.delete(tileId);
        state.pendingClaims = newPending;
        const previous = state.previousTiles.get(tileId);
        if (previous) {
          if (!previous.ownerId && state.optimisticClaimedDelta > 0) {
            state.optimisticClaimedDelta -= 1;
          }
          state.tiles.set(tileId, previous);
          state.previousTiles.delete(tileId);
          const tiles = Array.from(state.tiles.values());
          const tilesWithClaimedAt = tiles.filter((t) => t.claimedAt);
          if (tilesWithClaimedAt.length > 0) {
            const sorted = tilesWithClaimedAt.sort(
              (a, b) => new Date(b.claimedAt!).getTime() - new Date(a.claimedAt!).getTime()
            );
            state.lastActivity = sorted[0].claimedAt || null;
          } else {
            state.lastActivity = null;
          }
        }
      }),
    getClaimedTiles: () => {
      const state = get();
      return state.claimedTiles + state.optimisticClaimedDelta;
    },
    getUnclaimedTiles: () => {
      const state = get();
      const total = state.gridWidth * state.gridHeight;
      return total - (state.claimedTiles + state.optimisticClaimedDelta);
    },
    getLastActivity: () => {
      const state = get();
      return state.lastActivity;
    },
  }))
);
