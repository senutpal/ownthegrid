import type { ReactNode } from 'react';

interface AppShellProps {
  children: ReactNode;
}

export const AppShell = ({ children }: AppShellProps) => {
  return (
    <div className="relative min-h-screen bg-white text-black">
      <div className="relative">{children}</div>
    </div>
  );
};
