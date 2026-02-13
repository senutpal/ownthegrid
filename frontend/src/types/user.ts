export interface User {
  id: string;
  username: string;
  color: string;
  createdAt: string;
  lastSeen: string;
  token?: string;
}

export interface LeaderboardEntry {
  userId: string;
  username: string;
  color: string;
  tileCount: number;
  rank: number;
  optimisticDelta?: number;
}
