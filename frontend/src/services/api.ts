import axios from 'axios';

import type { Tile } from '../types/tile';
import type { LeaderboardEntry, User } from '../types/user';

const apiBase = (import.meta.env.VITE_API_URL as string | undefined) ?? window.location.origin;

const api = axios.create({
  baseURL: apiBase,
  withCredentials: true,
});

export interface RegisterResponse extends User {
  token: string;
}

export const registerUser = async (username: string): Promise<RegisterResponse> => {
  const response = await api.post('/api/users/register', { username });
  return response.data as RegisterResponse;
};

export const fetchBoard = async (): Promise<{
  tiles: Tile[];
  gridWidth: number;
  gridHeight: number;
  totalTiles: number;
  claimedTiles: number;
}> => {
  const response = await api.get('/api/board');
  return response.data as {
    tiles: Tile[];
    gridWidth: number;
    gridHeight: number;
    totalTiles: number;
    claimedTiles: number;
  };
};

export const fetchStats = async (): Promise<{
  totalTiles: number;
  claimedTiles: number;
  unclaimedTiles: number;
  onlineUsers: number;
  totalUsers: number;
  lastActivity: string | null;
}> => {
  const response = await api.get('/api/board/stats');
  return response.data as {
    totalTiles: number;
    claimedTiles: number;
    unclaimedTiles: number;
    onlineUsers: number;
    totalUsers: number;
    lastActivity: string | null;
  };
};

export const fetchLeaderboard = async (): Promise<LeaderboardEntry[]> => {
  const response = await api.get('/api/users/leaderboard');
  return (response.data as { leaderboard: LeaderboardEntry[] }).leaderboard;
};

export const fetchOnlineUsers = async (): Promise<User[]> => {
  const response = await api.get('/api/users/online');
  return (response.data as { users: User[] }).users;
};
