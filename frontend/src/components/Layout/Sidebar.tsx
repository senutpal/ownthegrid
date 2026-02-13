import { useUserStore } from '../../store/userStore';
import { useBoardStore } from '../../store/boardStore';
import { formatRelativeTime } from '../../utils/time';
import { Leaderboard } from '../Leaderboard/Leaderboard';
import { useMemo } from 'react';

export const Sidebar = () => {
  const leaderboard = useUserStore((state) => state.leaderboard);
  const tiles = useBoardStore((state) => state.tiles);
  const gridWidth = useBoardStore((state) => state.gridWidth);
  const gridHeight = useBoardStore((state) => state.gridHeight);
  const lastActivity = useBoardStore((state) => state.lastActivity);

  const claimedTiles = useMemo(
    () => Array.from(tiles.values()).filter((t) => t.ownerId !== null && t.ownerId !== undefined).length,
    [tiles]
  );
  const unclaimedTiles = gridWidth * gridHeight - claimedTiles;

  return (
    <aside className="ui-text w-full max-w-sm space-y-6 rounded-2xl border border-black/10 bg-white p-5 shadow-[0_20px_60px_-40px_rgba(0,0,0,0.3)]">
      <div className="space-y-2">
        <p className="text-xs uppercase tracking-[0.28em] text-black/60">Board stats</p>
        <div className="grid grid-cols-2 gap-3">
          <div className="rounded-xl border border-black/10 bg-white p-3">
            <p className="text-xs text-black/60">Grid</p>
            <p className="text-lg font-semibold text-black">{gridWidth} × {gridHeight}</p>
          </div>
          <div className="rounded-xl border border-black/10 bg-white p-3">
            <p className="text-xs text-black/60">Claimed</p>
            <p className="text-lg font-semibold text-black">{claimedTiles}</p>
          </div>
          <div className="rounded-xl border border-black/10 bg-white p-3">
            <p className="text-xs text-black/60">Unclaimed</p>
            <p className="text-lg font-semibold text-black">{unclaimedTiles}</p>
          </div>
          <div className="rounded-xl border border-black/10 bg-white p-3">
            <p className="text-xs text-black/60">Last activity</p>
            <p className="text-lg font-semibold text-black">{formatRelativeTime(lastActivity)}</p>
          </div>
        </div>
      </div>

      <div className="space-y-2">
        <p className="text-xs uppercase tracking-[0.28em] text-black/60">Leaderboard</p>
        <Leaderboard entries={leaderboard} />
      </div>
    </aside>
  );
};
