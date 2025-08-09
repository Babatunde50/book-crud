// src/components/BookFormModal.tsx
"use client";

import { useEffect, useState } from "react";

export type BookFormData = {
  title: string;
  author: string;
  year: string;
};

type FieldErrors = Partial<Record<keyof BookFormData, string>>;

export default function BookFormModal({
  isOpen,
  onClose,
  onSubmit,
  initialData,
  fieldErrors,
  formError,
  submitting = false,
}: {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: BookFormData) => void | Promise<void>;
  initialData?: BookFormData;
  fieldErrors?: FieldErrors;
  formError?: string | null;
  submitting?: boolean;
}) {
  const [data, setData] = useState<BookFormData>({
    title: initialData?.title ?? "",
    author: initialData?.author ?? "",
    year: initialData?.year ?? "",
  });

  useEffect(() => {
    if (initialData) {
      setData(initialData);
    }
  }, [initialData, isOpen]);

  if (!isOpen) return null;

  const inputBase =
    "block w-full rounded-md bg-white py-1.5 px-3 text-gray-900 outline-1 -outline-offset-1 outline-gray-300 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6";
  const errorInput =
    "block w-full rounded-md bg-white py-1.5 pr-10 pl-3 text-red-900 outline-1 -outline-offset-1 outline-red-300 placeholder:text-red-300 focus:outline-2 focus:-outline-offset-2 focus:outline-red-600 sm:text-sm/6";

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/30">
      <div className="w-full max-w-md rounded-lg bg-white p-6 shadow-lg">
        <div className="mb-4 flex items-center justify-between">
          <h2 className="text-lg font-semibold">
            {initialData ? "Edit" : "Add"} Book
          </h2>
          <button
            onClick={onClose}
            className="text-gray-500 hover:text-gray-700"
          >
            ✕
          </button>
        </div>

        {formError ? (
          <div className="mb-4 rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700">
            {formError}
          </div>
        ) : null}

        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-900">
              Title
            </label>
            <input
              className={fieldErrors?.title ? errorInput : inputBase}
              value={data.title}
              onChange={(e) =>
                setData((d) => ({ ...d, title: e.target.value }))
              }
              placeholder="Clean Code"
            />
            {fieldErrors?.title ? (
              <p className="mt-1 text-sm text-red-600">{fieldErrors.title}</p>
            ) : null}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-900">
              Author
            </label>
            <input
              className={fieldErrors?.author ? errorInput : inputBase}
              value={data.author}
              onChange={(e) =>
                setData((d) => ({ ...d, author: e.target.value }))
              }
              placeholder="Robert C. Martin"
            />
            {fieldErrors?.author ? (
              <p className="mt-1 text-sm text-red-600">{fieldErrors.author}</p>
            ) : null}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-900">
              Year
            </label>
            <input
              className={fieldErrors?.year ? errorInput : inputBase}
              value={data.year}
              onChange={(e) => setData((d) => ({ ...d, year: e.target.value }))}
              inputMode="numeric"
              placeholder="2008"
            />
            {fieldErrors?.year ? (
              <p className="mt-1 text-sm text-red-600">{fieldErrors.year}</p>
            ) : null}
          </div>
        </div>

        <div className="mt-6 flex justify-end gap-2">
          <button
            onClick={onClose}
            className="rounded-md border border-gray-300 bg-white px-3 py-1.5 text-sm text-gray-700 hover:bg-gray-50"
            disabled={submitting}
          >
            Cancel
          </button>
          <button
            onClick={() => onSubmit(data)}
            className="rounded-md bg-indigo-600 px-3 py-1.5 text-sm text-white hover:bg-indigo-700 disabled:opacity-60"
            disabled={submitting}
          >
            {submitting ? "Saving..." : "Save"}
          </button>
        </div>
      </div>
    </div>
  );
}
