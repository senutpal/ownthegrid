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
  initBoard: (tiles: Tile[], gridWidth: number, gridHeight: number) => void;
  applyTileClaimed: (payload: TileClaimedPayload) => void;
  optimisticallyClaimTile: (tileId: number, userId: string, color: string) => void;
  revertOptimisticClaim: (tileId: number) => void;
}

export const useBoardStore = create<BoardState>()(
  immer((set) => ({
    tiles: new Map(),
    optimisticTiles: new Map(),
    gridWidth: Number(import.meta.env.VITE_GRID_WIDTH ?? 50),
    gridHeight: Number(import.meta.env.VITE_GRID_HEIGHT ?? 40),
    isLoaded: false,
    previousTiles: new Map(),
    pendingClaims: new Set(),

    initBoard: (tiles, gridWidth, gridHeight) =>
      set((state) => {
        state.tiles = new Map(tiles.map((tile) => [tile.id, tile]));
        state.gridWidth = gridWidth;
        state.gridHeight = gridHeight;
        state.isLoaded = true;
      }),

    applyTileClaimed: (payload) =>
      set((state) => {
        const existing = state.tiles.get(payload.tileId);
        if (existing) {
          existing.ownerId = payload.userId;
          existing.ownerUsername = payload.username;
          existing.ownerColor = payload.color;
          existing.claimedAt = payload.claimedAt;
        }
        state.optimisticTiles.delete(payload.tileId);
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
          tile.ownerId = userId;
          tile.ownerColor = color;
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
          state.tiles.set(tileId, previous);
          state.previousTiles.delete(tileId);
        }
      }),
  }))
);
