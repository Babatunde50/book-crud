"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { parseApiError } from "@/lib/api-error";
import { Book } from "@/lib/types";
import BookFormModal, { BookFormData } from "@/components/BookFormModal";

type Props = {
  book: Book;
};

export default function EditButton({ book }: Props) {
  const [open, setOpen] = useState(false);
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});
  const [formError, setFormError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const router = useRouter();

  const initial: BookFormData = {
    title: book.title,
    author: book.author,
    year: String(book.year),
  };

  const onEdit = async (data: BookFormData) => {
    setSubmitting(true);
    setFieldErrors({});
    setFormError(null);

    try {
      // Build a partial payload with only changed fields
      const payload: Partial<{ title: string; author: string; year: number }> =
        {};
      if (data.title !== book.title) payload.title = data.title;
      if (data.author !== book.author) payload.author = data.author;

      const newYear = Number(data.year);
      if (!Number.isNaN(newYear) && newYear !== book.year)
        payload.year = newYear;

      // If nothing changed, just close
      if (Object.keys(payload).length === 0) {
        setOpen(false);
        return;
      }

      const controller = new AbortController();
      const timeout = setTimeout(() => controller.abort(), 10000);

      const res = await fetch(`/api/books/${book.id}`, {
        method: "PUT",
        headers: { "content-type": "application/json" },
        body: JSON.stringify(payload),
        signal: controller.signal,
      });

      clearTimeout(timeout);

      if (!res.ok) {
        const { fieldErrors, formError } = await parseApiError(res);
        setFieldErrors(fieldErrors);
        setFormError(formError);
        return;
      }

      setOpen(false);
      router.refresh();
    } catch (err) {
      const { fieldErrors, formError } = await parseApiError(err);
      setFieldErrors(fieldErrors);
      setFormError(formError);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <>
      <button
        onClick={() => {
          setOpen(true);
          setFieldErrors({});
          setFormError(null);
        }}
        className="cursor-pointer rounded-md bg-indigo-600 text-white px-3 py-1.5 text-sm hover:bg-indigo-700"
      >
        Edit
      </button>

      <BookFormModal
        isOpen={open}
        onClose={() => setOpen(false)}
        onSubmit={onEdit}
        initialData={initial}
        fieldErrors={fieldErrors}
        formError={formError}
        submitting={submitting}
      />
    </>
  );
}
