// src/app/books/[id]/page.tsx

import Link from "next/link";
import { getBook } from "@/lib/books";
import { Metadata } from "next";
import EditButton from "@/components/EditButton";
import DeleteButton from "@/components/DeleteButton";
import { notFound } from "next/navigation";

type Props = {
  params: { id: string };
};

function isUuidV4(id: string) {
  return /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i.test(id);
}

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

  if (!isUuidV4(id)) return notFound(); 

  const book = await getBook(id);

  if (!book) return notFound();

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
