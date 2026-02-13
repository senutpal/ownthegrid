import { create } from 'zustand';
import { immer } from 'zustand/middleware/immer';

type WSStatus = 'Connecting' | 'Connected' | 'Disconnected';
type Sender = <T>(type: string, payload: T) => void;

interface WSState {
  status: WSStatus;
  setStatus: (status: WSStatus) => void;
  send: Sender;
  setSender: (sender: Sender | null) => void;
}

export const useWSStore = create<WSState>()(
  immer((set) => ({
    status: 'Disconnected',
    setStatus: (status) =>
      set((state) => {
        state.status = status;
      }),
    send: () => undefined,
    setSender: (sender) =>
      set((state) => {
        state.send = sender ?? (() => undefined);
      }),
  }))
);
