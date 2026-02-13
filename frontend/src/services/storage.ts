const USER_KEY = 'ownthegrid:user';

export interface StoredUser {
  id: string;
  username: string;
  color: string;
  token?: string;
}

export const saveUser = (user: StoredUser): void => {
  localStorage.setItem(USER_KEY, JSON.stringify(user));
};

export const loadUser = (): StoredUser | null => {
  const raw = localStorage.getItem(USER_KEY);
  if (!raw) return null;
  try {
    const parsed = JSON.parse(raw) as StoredUser;
    if (!parsed?.id || !parsed?.username) return null;
    return parsed;
  } catch {
    return null;
  }
};

export const clearUser = (): void => {
  localStorage.removeItem(USER_KEY);
};
