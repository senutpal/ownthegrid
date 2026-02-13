interface OnlineIndicatorProps {
  count: number;
}

export const OnlineIndicator = ({ count }: OnlineIndicatorProps) => {
  return (
    <div className="flex items-center gap-2 rounded-full border border-black/20 bg-white px-3 py-1 text-xs text-black">
      <span className="relative flex h-2 w-2">
        <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-green-500/40 opacity-60" />
        <span className="relative inline-flex h-2 w-2 rounded-full bg-green-500" />
      </span>
      {count} online
    </div>
  );
};
