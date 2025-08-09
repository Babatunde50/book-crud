"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { parseApiError } from "@/lib/api-error";
import { useToast } from "@/context/ToastContext";

type Props = {
  id: string;
  className?: string;
  label?: string;
  redirectTo?: string; // default: "/books"
};

export default function DeleteButton({
  id,
  className,
  label = "Delete",
  redirectTo = "/books",
}: Props) {
  const [open, setOpen] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();
  const { success, error: toastError } = useToast();
  const onDelete = async () => {
    setSubmitting(true);
    setError(null);

    try {
      const controller = new AbortController();
      const timeout = setTimeout(() => controller.abort(), 10000);

      const res = await fetch(`/api/books/${id}`, {
        method: "DELETE",
        signal: controller.signal,
      });

      clearTimeout(timeout);

      if (!res.ok) {
        const { formError } = await parseApiError(res);
        setError(formError ?? `Delete failed (${res.status})`);
        return;
      }

      setOpen(false);
      router.push(redirectTo);
      router.refresh();
      success("Book deleted");
    } catch (err) {
      const { formError } = await parseApiError(err);
      setError(formError);
      toastError(formError ?? "Something went wrong");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <>
      <button
        onClick={() => setOpen(true)}
        className={
          className ??
          "rounded-md bg-red-600 text-white px-3 py-1.5 text-sm hover:bg-red-700"
        }
      >
        {label}
      </button>

      {open && (
        <>
          {/* Backdrop */}
          <div
            className="fixed inset-0 z-40 bg-black/40"
            onClick={() => !submitting && setOpen(false)}
          />
          {/* Dialog */}
          <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div
              role="dialog"
              aria-modal="true"
              className="w-full max-w-sm rounded-lg bg-white shadow-xl border border-gray-200"
            >
              <div className="p-5">
                <h3 className="text-base font-semibold text-gray-900">
                  Delete this book?
                </h3>
                <p className="mt-2 text-sm text-gray-600">
                  This action cannot be undone.
                </p>

                {error && <p className="mt-3 text-sm text-red-600">{error}</p>}

                <div className="mt-5 flex items-center justify-end gap-2">
                  <button
                    onClick={() => setOpen(false)}
                    disabled={submitting}
                    className="rounded-md border border-gray-300 bg-white px-3 py-1.5 text-sm text-gray-700 hover:bg-gray-50"
                  >
                    Cancel
                  </button>
                  <button
                    onClick={onDelete}
                    disabled={submitting}
                    className="rounded-md bg-red-600 px-3 py-1.5 text-sm text-white hover:bg-red-700 disabled:opacity-60"
                  >
                    {submitting ? "Deleting…" : "Delete"}
                  </button>
                </div>
              </div>
            </div>
          </div>
        </>
      )}
    </>
  );
}
