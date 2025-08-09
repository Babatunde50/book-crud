import Link from "next/link";

import { AddButton } from "@/app/components/AddButton";
import { Metadata } from "next";
import { getBooks } from "@/lib/books";

export const metadata: Metadata = {
  title: "Books",
  description: "Books page",
};

export default async function BooksPage() {
  const books = await getBooks();

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Library</h1>
        <AddButton />
      </div>

      <div className="text-sm text-gray-500">
        {books.length} books available
      </div>

      <ul className="divide-y border rounded-md">
        {books.map((book) => (
          <li key={book.id} className="p-4 flex justify-between">
            <div>
              <p className="font-semibold">{book.title}</p>
              <p className="text-gray-500 text-sm">
                {book.author} · {book.year}
              </p>
            </div>
            <Link
              href={`/books/${book.id}`}
              className="text-sm text-indigo-600 hover:underline"
            >
              View
            </Link>
          </li>
        ))}
      </ul>
    </div>
  );
}
