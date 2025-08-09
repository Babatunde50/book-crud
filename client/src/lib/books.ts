import { Book } from "./types";

export async function apiGet(path: string, init?: RequestInit) {
  const base = process.env.API_BASE_URL!;
  const res = await fetch(`${base}${path}`, {
    cache: "no-store",
    ...init,
  });
  if (!res.ok) throw new Error(`GET ${path} failed: ${res.status}`);
  return res.json();
}

export async function getBook(id: string): Promise<Book | null> {
  const base = process.env.API_BASE_URL!;
  const res = await fetch(`${base}/books/${id}`, { cache: "no-store" });

  if (res.status === 404) return null;
  if (!res.ok) throw new Error(`GET /books/${id} failed: ${res.status}`);

  return res.json();
}

export async function getBooks(): Promise<Book[]> {
  return apiGet("/books");
}
