import { useEffect } from 'react';
import { useUIStore } from '../../store/uiStore';

export const Toast = () => {
  const toasts = useUIStore((state) => state.toasts);
  const removeToast = useUIStore((state) => state.removeToast);

  useEffect(() => {
    if (toasts.length === 0) return;
    const timers = toasts.map((toast) =>
      setTimeout(() => removeToast(toast.id), 3500)
    );
    return () => {
      timers.forEach((timer) => clearTimeout(timer));
    };
  }, [toasts, removeToast]);

  return (
    <div className="fixed bottom-6 right-6 z-50 flex w-72 flex-col gap-3">
      {toasts.map((toast) => (
        <div
          key={toast.id}
          className={`rounded-xl border px-4 py-3 text-sm shadow-lg ${
            toast.type === 'error'
              ? 'border-black/40 bg-white text-black'
              : toast.type === 'success'
                ? 'border-black/40 bg-white text-black'
                : 'border-black/20 bg-white text-black'
          }`}
        >
          {toast.message}
        </div>
      ))}
    </div>
  );
};
