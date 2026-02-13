import type { LeaderboardEntry } from '../../types/user';
import { LeaderboardRow } from './LeaderboardRow';
import { AnimatePresence, motion } from 'motion/react';
import { useUserStore } from '../../store/userStore';
import { useMemo } from 'react';

interface LeaderboardProps {
  entries: LeaderboardEntry[];
}

export const Leaderboard = ({ entries }: LeaderboardProps) => {
  const onlineUsers = useUserStore((state) => state.onlineUsers);

  const onlineUserIds = useMemo(
    () => new Set(onlineUsers.map((u) => u.id)),
    [onlineUsers]
  );

  if (!entries.length) {
    return (
      <div className="rounded-xl border border-dashed border-black/30 p-4 text-sm text-black/60">
        No leaderboard data yet.
      </div>
    );
  }

  return (
    <motion.div
      layout
      className="max-h-[300px] space-y-2 overflow-y-auto pr-1"
    >
      <AnimatePresence initial={false}>
        {entries.map((entry) => (
          <motion.div
            key={entry.userId}
            layout
            initial={{ opacity: 0, y: 6 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -6 }}
            transition={{ duration: 0.2 }}
          >
            <LeaderboardRow
              entry={entry}
              isOnline={onlineUserIds.has(entry.userId)}
            />
          </motion.div>
        ))}
      </AnimatePresence>
    </motion.div>
  );
};
