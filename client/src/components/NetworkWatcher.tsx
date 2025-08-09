"use client";

import { useEffect } from "react";
import { useToast } from "@/context/ToastContext";

export default function NetworkWatcher() {
  const { error, info } = useToast();

  useEffect(() => {
    const onOffline = () => error("You’re offline. Actions may fail.");
    const onOnline = () => info("Back online.");

    window.addEventListener("offline", onOffline);
    window.addEventListener("online", onOnline);
    return () => {
      window.removeEventListener("offline", onOffline);
      window.removeEventListener("online", onOnline);
    };
  }, [error, info]);

  return null;
}
