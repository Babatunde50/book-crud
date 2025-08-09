package bookdb

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Babatunde50/byfood-assessment/server/business/book"
	"github.com/Babatunde50/byfood-assessment/server/internal/database"
	"github.com/google/uuid"
)

type Store struct {
	db *database.DB
}

// New creates a new bookdb store that satisfies the book.Store interface.
func New(db *database.DB) *Store {
	return &Store{db: db}
}

// Create inserts a new book into the database.
func (s *Store) Create(ctx context.Context, bk book.Book) error {
	const query = `
		INSERT INTO books (
			id, title, author, year, date_created, date_updated, version
		)
		VALUES (
			:id, :title, :author, :year, :date_created, :date_updated, :version
		)`

	dbBook := toDBBook(bk)

	if _, err := s.db.NamedExecContext(ctx, query, dbBook); err != nil {
		return err
	}

	return nil
}

// Update modifies an existing book record.
func (s *Store) Update(ctx context.Context, bk book.Book) error {
	const query = `
		UPDATE books SET
			title = :title,
			author = :author,
			year = :year,
			date_updated = :date_updated,
			version = :version
		WHERE id = :id AND version = :version - 1`

	dbBook := toDBBook(bk)

	result, err := s.db.NamedExecContext(ctx, query, dbBook)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return book.ErrNotFound
	}

	return nil
}

// Delete removes a book by its ID.
func (s *Store) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM books WHERE id = $1`

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return book.ErrNotFound
	}

	return nil
}

// QueryByID retrieves a book by its ID.
func (s *Store) QueryByID(ctx context.Context, id uuid.UUID) (book.Book, error) {
	const query = `SELECT * FROM books WHERE id = $1`

	var dbBook dbBook
	if err := s.db.GetContext(ctx, &dbBook, query, id); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return book.Book{}, book.ErrNotFound
		}

		return book.Book{}, err
	}

	return toCoreBook(dbBook), nil
}

// Query retrieves all books.
func (s *Store) QueryAll(ctx context.Context) ([]book.Book, error) {
	const query = `SELECT * FROM books ORDER BY date_created DESC`

	var dbBooks []dbBook
	if err := s.db.SelectContext(ctx, &dbBooks, query); err != nil {
		return nil, err
	}

	books := make([]book.Book, len(dbBooks))
	for i, dbBook := range dbBooks {
		books[i] = toCoreBook(dbBook)
	}

	return books, nil
}
