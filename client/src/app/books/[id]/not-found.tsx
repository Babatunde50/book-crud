export default function BookNotFound() {
  return (
    <div className="rounded-lg border bg-white p-6 shadow-sm">
      <h2 className="text-lg font-semibold">Book not found</h2>
      <p className="mt-2 text-sm text-gray-600">
        We couldn’t find that book. It may have been deleted.
      </p>
      <a
        href="/books"
        className="mt-4 inline-block rounded-md border px-3 py-1.5 text-sm hover:bg-gray-50"
      >
        Back to Books
      </a>
    </div>
  );
}
