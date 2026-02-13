import { useEffect } from 'react';

import { Grid } from './components/Grid';
import { Header, Sidebar } from './components/Layout';
import { AppShell } from './components/AppShell';
import { Toast } from './components/UI/Toast';
import { UsernameModal } from './components/Onboarding';
import { useUserStore } from './store/userStore';
import { useWebSocket } from './hooks/useWebSocket';
import { fetchLeaderboard, fetchOnlineUsers } from './services/api';

const App = () => {
  const setOnlineUsers = useUserStore((state) => state.setOnlineUsers);
  const setLeaderboard = useUserStore((state) => state.setLeaderboard);
  useWebSocket();

  useEffect(() => {
    const load = async () => {
      const [online, leaderboard] = await Promise.all([
        fetchOnlineUsers(),
        fetchLeaderboard(),
      ]);
      setOnlineUsers(online);
      setLeaderboard(leaderboard);
    };
    load();
  }, [setLeaderboard, setOnlineUsers]);

  return (
    <AppShell>
      <Header />
      <main className="mx-auto flex w-full max-w-6xl flex-col gap-8 px-6 pb-12 lg:flex-row">
        <div className="flex-1 rounded-3xl border border-black/10 bg-white p-6 shadow-[0_30px_80px_-60px_rgba(0,0,0,0.35)]">
          <Grid />
        </div>
        <Sidebar />
      </main>
      <UsernameModal />
      <Toast />
    </AppShell>
  );
};

export default App;
