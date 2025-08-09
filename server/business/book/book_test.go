package book_test

import (
	"context"
	"errors"
	"runtime/debug"
	"testing"
	"time"

	"github.com/Babatunde50/byfood-assessment/server/business/book"
	"github.com/Babatunde50/byfood-assessment/server/business/book/bookdb"
	"github.com/Babatunde50/byfood-assessment/server/internal/dbtest"
	"github.com/google/uuid"
)

func Test_Book_CRUD(t *testing.T) {
	c, err := dbtest.StartDB(t)
	if err != nil {
		t.Fatalf("starting test DB: %v", err)
	}
	defer dbtest.StopDB(c)

	test := dbtest.NewTest(t, c)
	defer func() {
		if r := recover(); r != nil {
			t.Log(r)
			t.Error(string(debug.Stack()))
		}
		test.Teardown()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	store := bookdb.New(test.DB)
	core := book.NewCore(store)

	// ---------------------------------------------------------------------

	t.Log("Given the need to CRUD a book record")

	newBook := book.NewBook{
		Title:  "Test Driven Development",
		Author: "Kent Beck",
		Year:   2003,
	}

	t.Log("\tWhen creating a new book")
	createdBook, err := core.Create(ctx, newBook)
	if err != nil {
		t.Fatalf("\t\tShould be able to create book: %s", err)
	}

	if createdBook.Title != newBook.Title {
		t.Errorf("\t\tShould match title: got %q, want %q", createdBook.Title, newBook.Title)
	}

	// ---------------------------------------------------------------------

	t.Log("\tWhen querying the book by ID")
	queriedBook, err := core.QueryByID(ctx, createdBook.ID)
	if err != nil {
		t.Fatalf("\t\tShould be able to retrieve book: %s", err)
	}

	if queriedBook.ID != createdBook.ID {
		t.Errorf("\t\tQueried ID doesn't match created ID: got %v, want %v", queriedBook.ID, createdBook.ID)
	}

	// ---------------------------------------------------------------------

	t.Log("\tWhen updating the book")
	updatedTitle := "Refactoring"
	updatedAuthor := "Martin Fowler"
	updatedYear := 1999

	update := book.UpdateBook{
		Title:  &updatedTitle,
		Author: &updatedAuthor,
		Year:   &updatedYear,
	}

	updatedBook, err := core.Update(ctx, createdBook, update)
	if err != nil {
		t.Fatalf("\t\tShould be able to update book: %s", err)
	}

	if updatedBook.Title != updatedTitle {
		t.Errorf("\t\tUpdated title mismatch: got %q, want %q", updatedBook.Title, updatedTitle)
	}

	// ---------------------------------------------------------------------

	t.Log("\tWhen deleting the book")
	err = core.Delete(ctx, updatedBook.ID)
	if err != nil {
		t.Fatalf("\t\tShould be able to delete book: %s", err)
	}

	_, err = core.QueryByID(ctx, updatedBook.ID)
	if err == nil {
		t.Errorf("\t\tShould not be able to retrieve deleted book")
	}

	if err != nil && !errors.Is(err, book.ErrNotFound) {
		t.Fatalf("\t\tExpected ErrNotFound after delete, got: %s", err)
	}

	// ---------------------------------------------------------------------

	t.Log("\tWhen querying non-existent book")
	_, err = core.QueryByID(ctx, uuid.New())
	if err == nil {
		t.Errorf("\t\tExpected error querying non-existent book")
	}
	if err != nil && !errors.Is(err, book.ErrNotFound) {
		t.Errorf("\t\tExpected ErrNotFound, got %v", err)
	}
}
