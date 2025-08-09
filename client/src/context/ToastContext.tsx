// src/context/ToastContext.tsx
"use client";

import {
  createContext,
  useCallback,
  useContext,
  useMemo,
  useRef,
  useState,
} from "react";

type Variant = "success" | "error" | "info";

type Toast = {
  id: string;
  message: string;
  variant: Variant;
};

type ShowOptions = {
  variant?: Variant;
  duration?: number; // ms
};

type ToastContextValue = {
  show(message: string, opts?: ShowOptions): string;
  success(message: string, duration?: number): string;
  error(message: string, duration?: number): string;
  info(message: string, duration?: number): string;
  dismiss(id: string): void;
  clear(): void;
};

const ToastContext = createContext<ToastContextValue | null>(null);

export function ToastProvider({ children }: { children: React.ReactNode }) {
  const [toasts, setToasts] = useState<Toast[]>([]);
  const timers = useRef<Record<string, number>>({});

  const dismiss = useCallback((id: string) => {
    setToasts((t) => t.filter((x) => x.id !== id));
    if (timers.current[id]) {
      window.clearTimeout(timers.current[id]);
      delete timers.current[id];
    }
  }, []);

  const show = useCallback(
    (message: string, opts?: ShowOptions) => {
      const id = `${Date.now()}-${Math.random().toString(36).slice(2)}`;
      const variant = opts?.variant ?? "info";
      const duration = opts?.duration ?? 4000;

      setToasts((t) => [{ id, message, variant }, ...t]);

      timers.current[id] = window.setTimeout(() => dismiss(id), duration);

      return id;
    },
    [dismiss]
  );

  const success = useCallback(
    (m: string, d?: number) => show(m, { variant: "success", duration: d }),
    [show]
  );
  const error = useCallback(
    (m: string, d?: number) => show(m, { variant: "error", duration: d }),
    [show]
  );
  const info = useCallback(
    (m: string, d?: number) => show(m, { variant: "info", duration: d }),
    [show]
  );

  const clear = useCallback(() => {
    Object.values(timers.current).forEach((t) => window.clearTimeout(t));
    timers.current = {};
    setToasts([]);
  }, []);

  const value = useMemo(
    () => ({ show, success, error, info, dismiss, clear }),
    [show, success, error, info, dismiss, clear]
  );

  return (
    <ToastContext.Provider value={value}>
      {children}

      {/* Toast viewport */}
      <div className="pointer-events-none fixed top-4 right-4 z-50 flex w-[calc(100%-2rem)] max-w-sm flex-col gap-2">
        {toasts.map((t) => (
          <div
            key={t.id}
            className={[
              "pointer-events-auto rounded-lg border px-4 py-3 shadow-sm transition-all",
              t.variant === "success" &&
                "border-green-200 bg-green-50 text-green-900",
              t.variant === "error" && "border-red-200 bg-red-50 text-red-900",
              t.variant === "info" && "border-sky-200 bg-sky-50 text-sky-900",
            ].join(" ")}
            role="status"
            aria-live="polite"
          >
            <div className="flex items-start justify-between gap-3">
              <p className="text-sm">{t.message}</p>
              <button
                onClick={() => dismiss(t.id)}
                className="shrink-0 rounded p-1 text-gray-500 hover:bg-black/5"
                aria-label="Dismiss"
              >
                ×
              </button>
            </div>
          </div>
        ))}
      </div>
    </ToastContext.Provider>
  );
}

export function useToast() {
  const ctx = useContext(ToastContext);
  if (!ctx) throw new Error("useToast must be used within a ToastProvider");
  return ctx;
}
