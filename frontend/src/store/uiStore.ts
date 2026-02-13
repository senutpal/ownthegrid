import { create } from 'zustand';
import { immer } from 'zustand/middleware/immer';

export interface Toast {
  id: string;
  message: string;
  type: 'success' | 'error' | 'info';
}

interface UIState {
  showUsernameModal: boolean;
  toasts: Toast[];
  lastToastTime: Map<string, number>;
  setShowUsernameModal: (show: boolean) => void;
  addToast: (toast: Toast) => void;
  removeToast: (id: string) => void;
}

export const useUIStore = create<UIState>()(
  immer((set) => ({
    showUsernameModal: true,
    toasts: [],
    lastToastTime: new Map(),
    setShowUsernameModal: (show) =>
      set((state) => {
        state.showUsernameModal = show;
      }),
    addToast: (toast) =>
      set((state) => {
        const key = `${toast.type}:${toast.message}`;
        const now = Date.now();
        const lastShown = state.lastToastTime.get(key) ?? 0;
        if (now - lastShown < 1500) {
          return;
        }
        state.lastToastTime.set(key, now);
        state.toasts.push(toast);
      }),
    removeToast: (id) =>
      set((state) => {
        state.toasts = state.toasts.filter((toast) => toast.id !== id);
      }),
  }))
);
