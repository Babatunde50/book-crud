export default function NotFound() {
  return (
    <div className="min-h-[60vh] grid place-items-center p-6">
      <div className="text-center space-y-2">
        <h1 className="text-2xl font-bold">Page not found</h1>
        <p className="text-gray-600">
          The page you’re looking for doesn’t exist.
        </p>
        <a
          href="/books"
          className="inline-block mt-2 rounded-md border px-3 py-1.5 text-sm hover:bg-gray-50"
        >
          Go to Books
        </a>
      </div>
    </div>
  );
}
