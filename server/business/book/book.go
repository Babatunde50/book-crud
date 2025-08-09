// business/book/book.go
package book

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound      = errors.New("book not found")
	ErrTitleConflict = errors.New("book title already exists")
)

// Storer defines the behavior the book package expects from the data store layer.
type Storer interface {
	Create(ctx context.Context, book Book) error
	Update(ctx context.Context, book Book) error
	Delete(ctx context.Context, bookID uuid.UUID) error
	QueryByID(ctx context.Context, bookID uuid.UUID) (Book, error)
	QueryAll(ctx context.Context) ([]Book, error)
}

// Core manages the set of APIs for book access.
type Core struct {
	storer Storer
}

// NewCore constructs a core for book API access.
func NewCore(storer Storer) *Core {
	return &Core{
		storer: storer,
	}
}

// Create adds a new book to the system.
func (c *Core) Create(ctx context.Context, nb NewBook) (Book, error) {
	now := time.Now()

	book := Book{
		ID:          uuid.New(),
		Title:       nb.Title,
		Author:      nb.Author,
		Year:        nb.Year,
		DateCreated: now,
		DateUpdated: now,
		Version:     1,
	}

	if err := c.storer.Create(ctx, book); err != nil {
		return Book{}, fmt.Errorf("create: %w", err)
	}

	return book, nil
}

// Update modifies information about a book.
func (c *Core) Update(ctx context.Context, book Book, ub UpdateBook) (Book, error) {
	if ub.Title != nil {
		book.Title = *ub.Title
	}

	if ub.Author != nil {
		book.Author = *ub.Author
	}

	if ub.Year != nil {
		book.Year = *ub.Year
	}

	book.DateUpdated = time.Now()
	book.Version++

	if err := c.storer.Update(ctx, book); err != nil {
		return Book{}, fmt.Errorf("update: %w", err)
	}

	return book, nil
}

// Delete removes a book from the system.
func (c *Core) Delete(ctx context.Context, bookID uuid.UUID) error {
	if err := c.storer.Delete(ctx, bookID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

// QueryByID finds a book by its ID.
func (c *Core) QueryByID(ctx context.Context, bookID uuid.UUID) (Book, error) {
	book, err := c.storer.QueryByID(ctx, bookID)
	if err != nil {
		return Book{}, fmt.Errorf("query: id[%s]: %w", bookID, err)
	}
	return book, nil
}

// QueryAll returns all books in the system.
func (c *Core) QueryAll(ctx context.Context) ([]Book, error) {
	books, err := c.storer.QueryAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("query all: %w", err)
	}
	return books, nil
}
