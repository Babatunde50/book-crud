package book

import (
	"time"

	"github.com/google/uuid"
)

// Book represents information about a book record.
type Book struct {
	ID          uuid.UUID
	Title       string
	Author      string
	Year        int
	DateCreated time.Time
	DateUpdated time.Time
	Version     int
}

// NewBook holds data required to create a new book.
type NewBook struct {
	Title  string
	Author string
	Year   int
}

// UpdateBook holds data required to update an existing book.
type UpdateBook struct {
	Title  *string
	Author *string
	Year   *int
}
