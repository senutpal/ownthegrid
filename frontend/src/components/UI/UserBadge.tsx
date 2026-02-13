interface UserBadgeProps {
  username: string;
  color?: string;
}

export const UserBadge = ({ username, color }: UserBadgeProps) => {
  return (
    <div className="flex items-center gap-2 rounded-full border border-black/20 bg-white px-3 py-1 text-sm text-black">
      <span
        className="h-2.5 w-2.5 rounded-full"
        style={{ backgroundColor: color || '#000000' }}
      />
      {username}
    </div>
  );
};
