export type Book = {
  id: string;
  title: string;
  author: string;
  year: number;
  date_created: string;
  date_updated: string;
  version: number;
};

export type NewBook = {
  title: string;
  author: string;
  year: number;
};

export type UpdateBook = Partial<NewBook>;
