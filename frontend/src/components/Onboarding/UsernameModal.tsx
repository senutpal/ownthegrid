import { useState } from 'react';

import { useUser } from '../../hooks/useUser';
import { useUIStore } from '../../store/uiStore';

export const UsernameModal = () => {
  const { register, isLoading } = useUser();
  const show = useUIStore((state) => state.showUsernameModal);
  const [username, setUsername] = useState('');

  if (!show) return null;

  return (
    <div className="fixed inset-0 z-40 flex items-center justify-center bg-black/80 px-4">
      <div className="w-full max-w-md rounded-2xl border border-black/10 bg-white p-6 shadow-xl">
        <p className="text-xs uppercase tracking-[0.3em] text-black/60">Welcome</p>
        <h2 className="mt-2 text-2xl font-semibold text-black">Pick a username</h2>
        <p className="mt-2 text-sm text-black/60">
          This will be visible on the grid. Keep it short.
        </p>

        <input
          value={username}
          onChange={(event) => setUsername(event.target.value)}
          placeholder="PixelKing"
          className="mt-4 w-full rounded-xl border border-black/20 bg-white px-4 py-3 text-sm text-black placeholder:text-black/40 focus:border-black focus:outline-none"
        />

        <button
          type="button"
          onClick={() => register(username)}
          disabled={isLoading || username.trim().length < 2}
          className="mt-4 w-full rounded-xl bg-black px-4 py-3 text-sm font-semibold text-white transition hover:bg-black/80 disabled:cursor-not-allowed disabled:bg-black/40"
        >
          {isLoading ? 'Creating...' : 'Enter the grid'}
        </button>
      </div>
    </div>
  );
};
