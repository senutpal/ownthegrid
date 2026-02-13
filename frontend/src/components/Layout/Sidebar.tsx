import { useUserStore } from '../../store/userStore';
import { useBoardStore } from '../../store/boardStore';
import { useEffect, useState } from 'react';
import { fetchStats } from '../../services/api';
import { formatRelativeTime } from '../../utils/time';
import { Leaderboard } from '../Leaderboard/Leaderboard';

interface BoardStats {
  totalTiles: number;
  claimedTiles: number;
  unclaimedTiles: number;
  onlineUsers: number;
  totalUsers: number;
  lastActivity: string | null;
}

export const Sidebar = () => {
  const leaderboard = useUserStore((state) => state.leaderboard);
  const gridWidth = useBoardStore((state) => state.gridWidth);
  const gridHeight = useBoardStore((state) => state.gridHeight);
  const [stats, setStats] = useState<BoardStats | null>(null);

  // Load stats on mount and refresh every 10 seconds
  useEffect(() => {
    let cancelled = false;
    
    const loadStats = async () => {
      const response = await fetchStats();
      if (!cancelled) {
        setStats(response);
      }
    };

    loadStats();
    const interval = setInterval(loadStats, 10000);
    
    return () => {
      cancelled = true;
      clearInterval(interval);
    };
  }, []);

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
            <p className="text-lg font-semibold text-black">{stats?.claimedTiles ?? '—'}</p>
          </div>
          <div className="rounded-xl border border-black/10 bg-white p-3">
            <p className="text-xs text-black/60">Unclaimed</p>
            <p className="text-lg font-semibold text-black">{stats?.unclaimedTiles ?? '—'}</p>
          </div>
          <div className="rounded-xl border border-black/10 bg-white p-3">
            <p className="text-xs text-black/60">Last activity</p>
            <p className="text-lg font-semibold text-black">{formatRelativeTime(stats?.lastActivity ?? null)}</p>
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
