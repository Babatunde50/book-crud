// src/app/books/[id]/page.tsx

import Link from "next/link";
import { getBook } from "@/lib/books";
import { Metadata } from "next";
import EditButton from "@/app/components/EditButton";
import DeleteButton from "@/app/components/DeleteButton";

type Props = {
  params: { id: string };
};

export async function generateMetadata({
  params,
}: {
  params: Promise<{ id: string }>;
}): Promise<Metadata> {
  const { id } = await params;
  return {
    title: `Book #${id} | byFood Library`,
    description: `Details for book ID ${id}`,
  };
}

export default async function BookDetailPage({ params }: Props) {
  const { id } = await params;

  const book = await getBook(id);

  return (
    <div className="space-y-6">
      {/* Back button */}
      <div>
        <Link
          href="/books"
          className="inline-flex items-center gap-2 text-sm text-gray-600 hover:text-black"
        >
          ← Back to books
        </Link>
      </div>

      {/* Book Card */}
      <div className="rounded-lg border border-gray-200 p-6 shadow-sm bg-white space-y-4">
        <div className="flex items-start justify-between">
          <div>
            <h1 className="text-2xl font-bold">{book.title}</h1>
            <p className="text-gray-600">By {book.author}</p>
            <p className="text-sm text-gray-400">{book.year}</p>
          </div>

          <div className="flex gap-2">
            <EditButton book={book} />
            <DeleteButton id={book.id} />
          </div>
        </div>
      </div>
    </div>
  );
}
