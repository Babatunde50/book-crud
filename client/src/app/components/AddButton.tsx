"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import BookFormModal, { BookFormData } from "@/app/components/BookFormModal";
import { parseApiError } from "@/lib/api-error";

export function AddButton() {
  const [open, setOpen] = useState(false);
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});
  const [formError, setFormError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const router = useRouter();

  const onCreate = async (data: BookFormData) => {
    setFormError(null);
    setFieldErrors({});
    setSubmitting(true);

    try {
      const controller = new AbortController();
      const timeout = setTimeout(() => controller.abort(), 10000);

      const res = await fetch("/api/books", {
        method: "POST",
        headers: { "content-type": "application/json" },
        body: JSON.stringify({
          title: data.title,
          author: data.author,
          year: Number(data.year),
        }),
        signal: controller.signal,
      });

      clearTimeout(timeout);

      if (!res.ok) {
        const parsed = await parseApiError(res);
        setFormError(parsed.formError);
        setFieldErrors(parsed.fieldErrors);
        setSubmitting(false);
        return;
      }
      setOpen(false);
      router.refresh();
    } catch (err) {
      const parsed = await parseApiError(err);
      setFormError(parsed.formError);
      setFieldErrors(parsed.fieldErrors);
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
        Add Book
      </button>

      <BookFormModal
        isOpen={open}
        onClose={() => setOpen(false)}
        onSubmit={onCreate}
        fieldErrors={fieldErrors}
        formError={formError}
        submitting={submitting}
      />
    </>
  );
}
