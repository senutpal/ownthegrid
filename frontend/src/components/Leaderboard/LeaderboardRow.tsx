import type { LeaderboardEntry } from '../../types/user';
import { useUserStore } from '../../store/userStore';

interface LeaderboardRowProps {
  entry: LeaderboardEntry;
  isOnline?: boolean;
}

export const LeaderboardRow = ({ entry, isOnline = false }: LeaderboardRowProps) => {
  const getTileCount = useUserStore((state) => state.getTileCount);
  const tileCount = getTileCount(entry.userId);

  return (
    <div className="flex items-center justify-between rounded-lg border border-black/10 bg-white px-3 py-2 text-sm text-black">
      <div className="flex items-center gap-2">
        <span className="text-xs text-black/50">#{entry.rank}</span>
        <span
          className="h-2.5 w-2.5 rounded-full"
          style={{ backgroundColor: entry.color }}
        />
        <span>{entry.username}</span>
        {isOnline && (
          <span className="text-green-500 text-xs">(online)</span>
        )}
      </div>
      <span className="text-black">{tileCount}</span>
    </div>
  );
};
