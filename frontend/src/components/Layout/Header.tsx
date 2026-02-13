import { useWSStore } from '../../store/wsStore';
import { useUserStore } from '../../store/userStore';
import { UserBadge } from '../UI/UserBadge';
import { OnlineIndicator } from '../UI/OnlineIndicator';

export const Header = () => {
  const status = useWSStore((state) => state.status);
  const currentUser = useUserStore((state) => state.currentUser);
  const onlineUsers = useUserStore((state) => state.onlineUsers);

  return (
    <header className="flex flex-col gap-4 px-6 py-6 mb-0 md:mb-4 lg:flex-row lg:items-center lg:justify-between">
      <div>
        <h1 className="text-3xl font-bold tracking-tight text-black lg:text-4xl mb-0">OwnTheGrid</h1>
        <p className="ui-text text-sm text-black/60 ml-0.5">Status: {status}</p>
      </div>
      {status !== 'Connected' && (
        <div className="ui-text rounded-full border border-black/20 bg-white px-4 py-2 text-xs tracking-[0.2em] text-black">
          Reconnecting
        </div>
      )}
      <div className="flex items-center gap-4">
        <OnlineIndicator count={onlineUsers.length} />
        {currentUser ? <UserBadge username={currentUser.username} color={currentUser.color} /> : null}
      </div>
    </header>
  );
};
