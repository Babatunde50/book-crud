package bookdb

import (
	"time"

	"github.com/Babatunde50/book-crud/server/business/book"
	"github.com/google/uuid"
)

// dbBook represents how a book is stored in the database.
type dbBook struct {
	ID          uuid.UUID `db:"id"`
	Title       string    `db:"title"`
	Author      string    `db:"author"`
	Year        int       `db:"year"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
	Version     int       `db:"version"`
}

// toCoreBook converts a dbBook to the core book.Book type.
func toCoreBook(db dbBook) book.Book {
	return book.Book{
		ID:          db.ID,
		Title:       db.Title,
		Author:      db.Author,
		Year:        db.Year,
		DateCreated: db.DateCreated,
		DateUpdated: db.DateUpdated,
		Version:     db.Version,
	}
}

// toDBBook converts book.Book to the dbBook type for database storage.
func toDBBook(bk book.Book) dbBook {
	return dbBook{
		ID:          bk.ID,
		Title:       bk.Title,
		Author:      bk.Author,
		Year:        bk.Year,
		DateCreated: bk.DateCreated,
		DateUpdated: bk.DateUpdated,
		Version:     bk.Version,
	}
}
