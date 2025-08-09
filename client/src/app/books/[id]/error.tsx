"use client";

export default function BookDetailError({
  error,
  reset,
}: {
  error: Error;
  reset: () => void;
}) {
  return (
    <div className="rounded-lg border bg-white p-6 shadow-sm">
      <h2 className="text-lg font-semibold">Can’t load that book</h2>
      <p className="mt-2 text-sm text-gray-600">
        {error.message || "Unknown error"}
      </p>
      <div className="mt-4 flex gap-2">
        <button
          onClick={reset}
          className="rounded-md bg-black px-3 py-1.5 text-white text-sm hover:opacity-90"
        >
          Retry
        </button>
        <a
          href="/books"
          className="rounded-md border px-3 py-1.5 text-sm hover:bg-gray-50"
        >
          Back to Books
        </a>
      </div>
    </div>
  );
}
