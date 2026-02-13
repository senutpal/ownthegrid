import { useCallback, useEffect } from 'react';

import { fetchBoard } from '../services/api';
import { useBoardStore } from '../store/boardStore';
import { useUserStore } from '../store/userStore';
import { useWSStore } from '../store/wsStore';

export const useBoard = () => {
  const initBoard = useBoardStore((state) => state.initBoard);
  const pendingClaims = useBoardStore((state) => state.pendingClaims);
  const optimisticallyClaimTile = useBoardStore((state) => state.optimisticallyClaimTile);
  const { currentUser, optimisticallyIncrementTileCount } = useUserStore();
  const send = useWSStore((state) => state.send);

  useEffect(() => {
    const load = async () => {
      const data = await fetchBoard();
      initBoard(data.tiles, data.gridWidth, data.gridHeight);
    };
    load();
  }, [initBoard]);

  const handleClaim = useCallback(
    (tileId: number) => {
      if (!currentUser) return;
      if (pendingClaims.has(tileId)) return;
      optimisticallyClaimTile(tileId, currentUser.id, currentUser.color);
      optimisticallyIncrementTileCount(currentUser.id);
      send('CLAIM_TILE', { tileId });
    },
    [currentUser, pendingClaims, optimisticallyClaimTile, optimisticallyIncrementTileCount, send]
  );

  return { handleClaim };
};
