package main

import (
	"time"

	"github.com/Babatunde50/byfood-assessment/server/business/book"
	"github.com/google/uuid"
)

// ---------- DTOs returned to clients ----------

// Book represents information about a book record.
type BookResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Year        int       `json:"year"`
	DateCreated time.Time `json:"date_created"`
	DateUpdated time.Time `json:"date_updated"`
	Version     int       `json:"-"`
}

// NewBook contains information needed to create a new book.
type NewBookRequest struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}

// UpdateBook contains information needed to update a book.
type UpdateBookRequest struct {
	Title  *string `json:"title,omitempty"`
	Author *string `json:"author,omitempty"`
	Year   *int    `json:"year,omitempty"`
}

// toBookResponse converts book.Book to the dbBook type for database storage.
func toBookResponse(bk book.Book) BookResponse {
	return BookResponse{
		ID:          bk.ID,
		Title:       bk.Title,
		Author:      bk.Author,
		Year:        bk.Year,
		DateCreated: bk.DateCreated,
		DateUpdated: bk.DateUpdated,
		Version:     bk.Version,
	}
}

func toBooksResponse(books []book.Book) []BookResponse {
	bookResponses := make([]BookResponse, len(books))
	for i, book := range books {
		bookResponses[i] = toBookResponse(book)
	}
	return bookResponses
}

type URLRequest struct {
	URL       string `json:"url"`
	Operation string `json:"operation"`
}

type URLResponse struct {
	ProcessedURL string `json:"processed_url"`
}
