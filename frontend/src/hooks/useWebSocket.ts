import { useCallback, useEffect, useRef } from 'react';

import { WebSocketService } from '../services/websocket';
import { useBoardStore } from '../store/boardStore';
import { useUserStore } from '../store/userStore';
import { useWSStore } from '../store/wsStore';
import { useUIStore } from '../store/uiStore';
import type { User } from '../types/user';
import tileStyles from '../components/Grid/Tile.module.css';
import type {
  ClaimRejectedPayload,
  InitBoardPayload,
  LeaderboardUpdatePayload,
  PongPayload,
  TileClaimedPayload,
  UserJoinedPayload,
  UserLeftPayload,
  ErrorPayload,
} from '../types/ws';

const claimSoundUrl = '/assets/sounds/claim.mp3';

export const useWebSocket = () => {
  const wsRef = useRef<WebSocketService | null>(null);
  const { initBoard, applyTileClaimed, revertOptimisticClaim } = useBoardStore();
  const { currentUser, setLeaderboard, setOnlineUsers, optimisticallyAddUser, confirmTileCount, revertTileCount } = useUserStore();
  const { setStatus, setSender } = useWSStore();
  const { addToast } = useUIStore();

  useEffect(() => {
    if (!currentUser?.id) return;

    const wsBase = (import.meta.env.VITE_WS_URL as string | undefined) ?? window.location.origin;
    const normalizedBase = wsBase.startsWith('http') ? wsBase.replace(/^http/, 'ws') : wsBase;
    const wsUrl = new URL('/ws', normalizedBase);
    wsUrl.searchParams.set('userId', currentUser.id);
    if (currentUser.token) {
      wsUrl.searchParams.set('token', currentUser.token);
    }
    const ws = new WebSocketService(wsUrl.toString(), setStatus);
    wsRef.current = ws;

    if (currentUser) {
      optimisticallyAddUser({
        id: currentUser.id,
        username: currentUser.username,
        color: currentUser.color,
        createdAt: new Date().toISOString(),
        lastSeen: new Date().toISOString(),
      });
    }

    ws.on<InitBoardPayload>('INIT_BOARD', ({ payload }) => {
      initBoard(
        payload.tiles,
        payload.gridWidth ?? Number(import.meta.env.VITE_GRID_WIDTH ?? 50),
        payload.gridHeight ?? Number(import.meta.env.VITE_GRID_HEIGHT ?? 40)
      );
    });

    ws.on<TileClaimedPayload>('TILE_CLAIMED', ({ payload }) => {
      applyTileClaimed(payload);
      confirmTileCount(payload.userId);
      const tile = document.querySelector(`[data-tile-id="${payload.tileId}"]`);
      if (tile) {
        tile.classList.add(tileStyles.justClaimed);
        window.setTimeout(() => tile.classList.remove(tileStyles.justClaimed), 280);
      }
      if (payload.userId === currentUser.id) {
        const audio = new Audio(claimSoundUrl);
        audio.volume = 0.4;
        void audio.play().catch(() => undefined);
      }
    });

    ws.on<ClaimRejectedPayload>('CLAIM_REJECTED', ({ payload }) => {
      revertOptimisticClaim(payload.tileId);
      if (currentUser) {
        revertTileCount(currentUser.id, false);
      }
      const message =
        payload.reason === 'ALREADY_CLAIMED'
          ? 'Too slow! Someone else got there first.'
          : payload.reason === 'INVALID_TILE'
            ? 'Invalid tile.'
            : 'Server error while claiming tile.';
      addToast({
        id: `${payload.tileId}-${Date.now()}`,
        message,
        type: 'error',
      });
    });

    ws.on<UserJoinedPayload>('USER_JOINED', ({ payload }) => {
      if (payload.userId !== currentUser?.id) {
        addToast({
          id: `join-${payload.userId}-${Date.now()}`,
          message: `${payload.username} joined the board`,
          type: 'info',
        });
      }
      setOnlineUsers((prev) => {
        if (prev.some((user) => user.id === payload.userId)) return prev;
        const nextUser: User = {
          id: payload.userId,
          username: payload.username,
          color: payload.color,
          createdAt: new Date().toISOString(),
          lastSeen: new Date().toISOString(),
        };
        return [...prev, nextUser];
      });
    });

    ws.on<UserLeftPayload>('USER_LEFT', ({ payload }) => {
      addToast({
        id: `left-${payload.userId}-${Date.now()}`,
        message: `${payload.username} left`,
        type: 'info',
      });
      setOnlineUsers((prev) => prev.filter((user) => user.id !== payload.userId));
    });

    ws.on<LeaderboardUpdatePayload>('LEADERBOARD_UPDATE', ({ payload }) => {
      const sorted = [...payload.leaderboard].sort((a, b) => b.tileCount - a.tileCount);
      const ranked = sorted.map((entry, index) => ({ ...entry, rank: index + 1 }));
      setLeaderboard(ranked);
    });

    ws.on<ErrorPayload>('ERROR', ({ payload }) => {
      addToast({
        id: `err-${payload.code}-${Date.now()}`,
        message: payload.message,
        type: 'error',
      });
    });

    ws.on<PongPayload>('PONG', () => {});

    ws.connect();
    setSender(ws.send.bind(ws));

    return () => {
      ws.disconnect();
      setSender(null);
    };
  }, [
    currentUser?.id,
    currentUser?.token,
    currentUser,
    initBoard,
    applyTileClaimed,
    revertOptimisticClaim,
    setLeaderboard,
    addToast,
    setStatus,
    setOnlineUsers,
    setSender,
    optimisticallyAddUser,
    confirmTileCount,
    revertTileCount,
  ]);

  const claimTile = useCallback((tileId: number) => {
    wsRef.current?.send('CLAIM_TILE', { tileId });
  }, []);

  return { claimTile };
};
