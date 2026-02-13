import { create } from 'zustand';
import { immer } from 'zustand/middleware/immer';

import type { User, LeaderboardEntry } from '../types/user';

interface UserState {
  currentUser: User | null;
  onlineUsers: User[];
  leaderboard: LeaderboardEntry[];
  setUser: (user: User | null) => void;
  setOnlineUsers: (users: User[] | ((prev: User[]) => User[])) => void;
  setLeaderboard: (entries: LeaderboardEntry[]) => void;
  optimisticallyAddUser: (user: User) => void;
  optimisticallyIncrementTileCount: (userId: string) => void;
  confirmTileCount: (userId: string) => void;
  revertTileCount: (userId: string, shouldDecrementBase?: boolean) => void;
  getTileCount: (userId: string) => number;
}

export const useUserStore = create<UserState>()(
  immer((set, get) => ({
    currentUser: null,
    onlineUsers: [],
    leaderboard: [],
    setUser: (user) =>
      set((state) => {
        state.currentUser = user;
      }),
    setOnlineUsers: (users) =>
      set((state) => {
        state.onlineUsers = typeof users === 'function' ? users(state.onlineUsers) : users;
      }),
    setLeaderboard: (entries) =>
      set((state) => {
        state.leaderboard = entries;
      }),
    optimisticallyAddUser: (user) =>
      set((state) => {
        if (!state.onlineUsers.some((u) => u.id === user.id)) {
          state.onlineUsers.push(user);
        }
      }),
    optimisticallyIncrementTileCount: (userId) =>
      set((state) => {
        const entry = state.leaderboard.find((e) => e.userId === userId);
        if (entry) {
          entry.optimisticDelta = (entry.optimisticDelta || 0) + 1;
        }
      }),
    confirmTileCount: (userId) =>
      set((state) => {
        const entry = state.leaderboard.find((e) => e.userId === userId);
        if (entry) {
          entry.tileCount += 1;
          entry.optimisticDelta = 0;
        }
      }),
    revertTileCount: (userId, shouldDecrementBase = false) =>
      set((state) => {
        const entry = state.leaderboard.find((e) => e.userId === userId);
        if (entry) {
          const currentDelta = entry.optimisticDelta || 0;
          if (currentDelta > 0) {
            entry.optimisticDelta = currentDelta - 1;
          }
          if (shouldDecrementBase && entry.tileCount > 0) {
            entry.tileCount -= 1;
          }
        }
      }),
    getTileCount: (userId) => {
      const state = get();
      const entry = state.leaderboard.find((e) => e.userId === userId);
      const base = entry?.tileCount || 0;
      const optimistic = entry?.optimisticDelta || 0;
      return base + optimistic;
    },
  }))
);
