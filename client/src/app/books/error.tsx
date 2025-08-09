"use client";

export default function BooksError({
  error,
  reset,
}: {
  error: Error;
  reset: () => void;
}) {
  return (
    <div className="rounded-lg border bg-white p-6 shadow-sm">
      <h2 className="text-lg font-semibold">Can’t load books</h2>
      <p className="mt-2 text-sm text-gray-600">
        {error.message || "Unknown error"}
      </p>
      <button
        onClick={reset}
        className="mt-4 rounded-md bg-black px-3 py-1.5 text-white text-sm hover:opacity-90"
      >
        Retry
      </button>
    </div>
  );
}
