import { useCallback, useEffect, useState } from 'react';

import { registerUser } from '../services/api';
import { loadUser, saveUser } from '../services/storage';
import { useUserStore } from '../store/userStore';
import { useUIStore } from '../store/uiStore';
import type { User } from '../types/user';

export const useUser = () => {
  const { currentUser, setUser } = useUserStore();
  const { setShowUsernameModal, addToast } = useUIStore();
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
      const stored = loadUser();
      if (stored) {
        const user: User = {
          id: stored.id,
          username: stored.username,
          color: stored.color,
          createdAt: new Date().toISOString(),
          lastSeen: new Date().toISOString(),
          token: stored.token,
        };
        setUser(user);
        setShowUsernameModal(false);
      } else {
        setShowUsernameModal(true);
      }
    }, [setUser, setShowUsernameModal]);

  const register = useCallback(
    async (username: string) => {
      setIsLoading(true);
      try {
        const response = await registerUser(username);
        const user: User = {
          id: response.id,
          username: response.username,
          color: response.color,
          createdAt: response.createdAt,
          lastSeen: response.lastSeen,
          token: response.token,
        };
        saveUser({
          id: response.id,
          username: response.username,
          color: response.color,
          token: response.token,
        });
        setUser(user);
        setShowUsernameModal(false);
        addToast({
          id: `welcome-${Date.now()}`,
          message: `Welcome, ${user.username}!`,
          type: 'success',
        });
      } catch {
        addToast({
          id: `register-${Date.now()}`,
          message: 'Username already taken or invalid',
          type: 'error',
        });
      } finally {
        setIsLoading(false);
      }
    },
    [addToast, setShowUsernameModal, setUser]
  );

  return { currentUser, isLoading, register };
};
