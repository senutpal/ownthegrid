import type { Tile } from './tile';
import type { User, LeaderboardEntry } from './user';

export type WSMessageType =
  | 'INIT_BOARD'
  | 'TILE_CLAIMED'
  | 'CLAIM_REJECTED'
  | 'USER_JOINED'
  | 'USER_LEFT'
  | 'LEADERBOARD_UPDATE'
  | 'ERROR'
  | 'PING'
  | 'PONG';

export interface WSMessage<T = unknown> {
  type: WSMessageType;
  payload: T;
  timestamp?: string;
}

export interface TileClaimedPayload {
  tileId: number;
  x: number;
  y: number;
  userId: string;
  username: string;
  color: string;
  claimedAt: string;
  previousOwner: string | null;
}

export interface ClaimRejectedPayload {
  tileId: number;
  reason: 'ALREADY_CLAIMED' | 'INVALID_TILE' | 'SERVER_ERROR';
  retryAfterMs?: number;
}

export interface InitBoardPayload {
  tiles: Tile[];
  user: User;
  onlineCount: number;
  gridWidth?: number;
  gridHeight?: number;
}

export interface UserJoinedPayload {
  userId: string;
  username: string;
  color: string;
  onlineCount: number;
}

export interface UserLeftPayload {
  userId: string;
  username: string;
  onlineCount: number;
}

export interface LeaderboardUpdatePayload {
  leaderboard: LeaderboardEntry[];
}

export interface ErrorPayload {
  code: string;
  message: string;
}

export interface PongPayload {}
