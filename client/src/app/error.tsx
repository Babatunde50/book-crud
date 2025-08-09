"use client";

import { useEffect } from "react";

export default function GlobalError({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  useEffect(() => {
    console.error(error);
  }, [error]);

  return (
    <html>
      <body className="min-h-screen grid place-items-center p-6">
        <div className="max-w-md w-full rounded-lg border bg-white p-6 shadow-sm">
          <h1 className="text-lg font-semibold">Something went wrong</h1>
          <p className="mt-2 text-sm text-gray-600">
            We hit a snag rendering this page. Try again or go home.
          </p>

          <div className="mt-4 flex gap-2">
            <button
              onClick={reset}
              className="rounded-md bg-black px-3 py-1.5 text-white text-sm hover:opacity-90"
            >
              Try again
            </button>
            <a
              href="/books"
              className="rounded-md border px-3 py-1.5 text-sm hover:bg-gray-50"
            >
              Back to Books
            </a>
          </div>
        </div>
      </body>
    </html>
  );
}
