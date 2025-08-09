import { apiGet } from "./http";
import { Book } from "./types";

export async function getBooks(): Promise<Book[]> {
  return apiGet("/books");
}

export async function getBook(id: string): Promise<Book> {
  return apiGet(`/books/${id}`);
}
