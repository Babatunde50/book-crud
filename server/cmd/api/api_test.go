package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Babatunde50/book-crud/server/business/book"
	"github.com/Babatunde50/book-crud/server/business/book/bookdb"
	"github.com/Babatunde50/book-crud/server/business/urlprocessor"
	"github.com/Babatunde50/book-crud/server/internal/database"
	"github.com/Babatunde50/book-crud/server/internal/docker"
	"github.com/google/uuid"
)

type testApp struct {
	handler  http.Handler
	teardown func()
}

func setupTestApp(t *testing.T) *testApp {
	t.Helper()

	image := "postgres:16"
	port := "5432"
	args := []string{"-e", "POSTGRES_PASSWORD=postgres"}

	c, err := docker.StartContainer(image, port, args...)

	if err != nil {
		t.Fatalf("could not start postgres container: %v", err)
	}

	dsn := fmt.Sprintf("postgres:postgres@%s/postgres?sslmode=disable", c.Host)
	var db *database.DB

	for i := 0; i < 20; i++ {
		db, err = database.New(dsn, true)
		if err == nil {
			break
		}
		t.Logf("waiting for db to be ready: %v", err)
		time.Sleep(time.Second)
	}

	if err != nil {
		t.Fatalf("could not connect to database: %v", err)
	}

	bookStore := bookdb.New(db)
	bookCore := book.NewCore(bookStore)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	app := &application{
		bookCore:         bookCore,
		urlProcessorCore: urlprocessor.New(),
		logger:           logger,
		db:               db,
	}

	h := app.routes()

	return &testApp{
		handler: h,
		teardown: func() {
			db.Close()
			_ = docker.StopContainer(c.ID)
		},
	}
}

func Test_CreateBookHandler(t *testing.T) {
	t.Parallel()
	test := setupTestApp(t)
	defer test.teardown()

	tests := []struct {
		name           string
		payload        string
		expectedStatus int
		assert         func(t *testing.T, body string)
	}{
		{
			name:           "valid book",
			payload:        `{"title":"Clean Code","author":"Robert C. Martin","year":2008}`,
			expectedStatus: http.StatusCreated,
			assert: func(t *testing.T, body string) {
				if !strings.Contains(body, "Clean Code") {
					t.Errorf("response body does not contain book title")
				}
			},
		},
		{
			name:           "missing title",
			payload:        `{"author":"Robert","year":2000}`,
			expectedStatus: http.StatusUnprocessableEntity,
			assert: func(t *testing.T, body string) {
				if !strings.Contains(body, "title") {
					t.Errorf("expected field error for title")
				}
			},
		},
		{
			name:           "invalid year",
			payload:        `{"title":"Go","author":"Ken","year":2050}`,
			expectedStatus: http.StatusUnprocessableEntity,
			assert: func(t *testing.T, body string) {
				if !strings.Contains(body, "year") {
					t.Errorf("expected field error for year")
				}
			},
		},
		{
			name:           "malformed json",
			payload:        `{title:bad}`,
			expectedStatus: http.StatusBadRequest,
			assert: func(t *testing.T, body string) {
				if !strings.Contains(body, "badly-formed") {
					t.Errorf("expected badly-formed json error")
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader([]byte(tc.payload)))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			test.handler.ServeHTTP(w, r)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tc.expectedStatus {
				t.Errorf("got status %d, want %d", res.StatusCode, tc.expectedStatus)
			}

			body, _ := io.ReadAll(res.Body)
			tc.assert(t, string(body))
		})
	}
}

func Test_ListBooksHandler(t *testing.T) {
	t.Parallel()
	test := setupTestApp(t)
	defer test.teardown()

	newBook := book.NewBook{
		Title:  "The Pragmatic Programmer",
		Author: "Andy Hunt",
		Year:   1999,
	}

	_, err := testCreateBook(test, newBook)

	if err != nil {
		t.Fatalf("failed to create book: %v", err)
	}

	tests := []struct {
		name           string
		setup          func()
		expectedStatus int
		assert         func(t *testing.T, body string)
	}{
		{
			name:           "list existing books",
			expectedStatus: http.StatusOK,
			assert: func(t *testing.T, body string) {
				if !strings.Contains(body, "The Pragmatic Programmer") {
					t.Errorf("expected book title not found in response: %s", body)
				}
				if !strings.Contains(body, "Andy Hunt") {
					t.Errorf("expected book author not found in response")
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup()
			}

			req := httptest.NewRequest(http.MethodGet, "/books", nil)
			res := httptest.NewRecorder()

			test.handler.ServeHTTP(res, req)

			if res.Result().StatusCode != tc.expectedStatus {
				t.Errorf("got status %d, want %d", res.Result().StatusCode, tc.expectedStatus)
			}

			body, _ := io.ReadAll(res.Body)
			tc.assert(t, string(body))
		})
	}
}

func Test_ShowBookHandler(t *testing.T) {
	t.Parallel()
	test := setupTestApp(t)
	defer test.teardown()

	newBook := book.NewBook{
		Title:  "Design Patterns",
		Author: "Erich Gamma",
		Year:   1994,
	}

	created, err := testCreateBook(test, newBook)

	if err != nil {
		t.Fatalf("failed to create book: %v", err)
	}

	tests := []struct {
		name           string
		bookID         string
		expectedStatus int
		assert         func(t *testing.T, body string)
	}{
		{
			name:           "valid ID",
			bookID:         created.ID.String(),
			expectedStatus: http.StatusOK,
			assert: func(t *testing.T, body string) {
				if !strings.Contains(body, "Design Patterns") {
					t.Errorf("expected book title not found")
				}
				if !strings.Contains(body, "Erich Gamma") {
					t.Errorf("expected book author not found")
				}
			},
		},
		{
			name:           "non-existent ID",
			bookID:         "123e4567-e89b-12d3-a456-426614174000", // valid UUID format but not in DB
			expectedStatus: http.StatusNotFound,
			assert: func(t *testing.T, body string) {
				if !strings.Contains(body, "could not be found") {
					t.Errorf("expected 'could not be found' in response, got: %s", body)
				}
			},
		},
		{
			name:           "invalid UUID",
			bookID:         "not-a-uuid",
			expectedStatus: http.StatusBadRequest,
			assert: func(t *testing.T, body string) {
				if !strings.Contains(body, "UUID") && !strings.Contains(body, "invalid") {
					t.Errorf("expected UUID-related error, got: %s", body)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/books/%s", tc.bookID)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			res := httptest.NewRecorder()

			test.handler.ServeHTTP(res, req)

			if res.Result().StatusCode != tc.expectedStatus {
				t.Errorf("got status %d, want %d", res.Result().StatusCode, tc.expectedStatus)
			}

			body, _ := io.ReadAll(res.Body)
			tc.assert(t, string(body))
		})
	}
}

func Test_UpdateBookHandler(t *testing.T) {
	t.Parallel()
	test := setupTestApp(t)
	defer test.teardown()

	// Seed a book to update
	newBook := book.NewBook{
		Title:  "Refactor",
		Author: "Martin Fowler",
		Year:   2000,
	}
	created, err := testCreateBook(test, newBook)
	if err != nil {
		t.Fatalf("failed to create book: %v", err)
	}

	tests := []struct {
		name           string
		id             string
		payload        string
		expectedStatus int
		assert         func(t *testing.T, body string)
	}{
		{
			name:           "valid full update",
			id:             created.ID.String(),
			payload:        `{"title":"Refactoring","author":"Fowler","year":2001}`,
			expectedStatus: http.StatusOK,
			assert: func(t *testing.T, body string) {
				if !strings.Contains(body, "Refactoring") {
					t.Errorf("expected updated title in response")
				}
			},
		},
		{
			name:           "partial update - author only",
			id:             created.ID.String(),
			payload:        `{"author":"Fowler Jr."}`,
			expectedStatus: http.StatusOK,
			assert: func(t *testing.T, body string) {
				if !strings.Contains(body, "Fowler Jr.") {
					t.Errorf("expected updated author in response")
				}
			},
		},
		{
			name:           "invalid UUID",
			id:             "not-a-uuid",
			payload:        `{"title":"X"}`,
			expectedStatus: http.StatusBadRequest,
			assert: func(t *testing.T, body string) {
				if !strings.Contains(body, "UUID") {
					t.Errorf("expected UUID error, got: %s", body)
				}
			},
		},
		{
			name:           "non-existent UUID",
			id:             uuid.New().String(),
			payload:        `{"title":"Y"}`,
			expectedStatus: http.StatusNotFound,
			assert: func(t *testing.T, body string) {
				if !strings.Contains(body, "could not be found") {
					t.Errorf("expected not found error, got: %s", body)
				}
			},
		},
		{
			name:           "invalid year",
			id:             created.ID.String(),
			payload:        `{"year":3000}`,
			expectedStatus: http.StatusUnprocessableEntity,
			assert: func(t *testing.T, body string) {
				if !strings.Contains(body, "year") {
					t.Errorf("expected validation error for year")
				}
			},
		},
		{
			name:           "malformed JSON",
			id:             created.ID.String(),
			payload:        `{bad}`,
			expectedStatus: http.StatusBadRequest,
			assert: func(t *testing.T, body string) {
				if !strings.Contains(body, "badly-formed") {
					t.Errorf("expected badly-formed JSON error, got: %s", body)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/books/%s", tc.id)
			r := httptest.NewRequest(http.MethodPut, url, bytes.NewReader([]byte(tc.payload)))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			test.handler.ServeHTTP(w, r)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tc.expectedStatus {
				t.Errorf("got status %d, want %d", res.StatusCode, tc.expectedStatus)
			}

			body, _ := io.ReadAll(res.Body)
			tc.assert(t, string(body))
		})
	}
}

func Test_DeleteBookHandler(t *testing.T) {
	t.Parallel()
	test := setupTestApp(t)
	defer test.teardown()

	// Seed one book
	created, err := testCreateBook(test, book.NewBook{
		Title:  "Design Patterns",
		Author: "Gang of Four",
		Year:   1994,
	})
	if err != nil {
		t.Fatalf("failed to create book: %v", err)
	}

	tests := []struct {
		name           string
		id             string
		expectedStatus int
		assert         func(t *testing.T, body string)
	}{
		{
			name:           "valid delete",
			id:             created.ID.String(),
			expectedStatus: http.StatusNoContent,
			assert: func(t *testing.T, body string) {
				if body != "" {
					t.Errorf("expected empty body on delete, got: %q", body)
				}
			},
		},
		{
			name:           "non-existent UUID",
			id:             uuid.New().String(),
			expectedStatus: http.StatusNotFound,
			assert: func(t *testing.T, body string) {
				if !strings.Contains(body, "could not be found") {
					t.Errorf("expected not found error, got: %s", body)
				}
			},
		},
		{
			name:           "invalid UUID",
			id:             "bad-id",
			expectedStatus: http.StatusBadRequest,
			assert: func(t *testing.T, body string) {
				if !strings.Contains(body, "UUID") {
					t.Errorf("expected bad UUID error, got: %s", body)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/books/%s", tc.id)
			r := httptest.NewRequest(http.MethodDelete, url, nil)
			w := httptest.NewRecorder()

			test.handler.ServeHTTP(w, r)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tc.expectedStatus {
				t.Errorf("got status %d, want %d", res.StatusCode, tc.expectedStatus)
			}

			body, _ := io.ReadAll(res.Body)
			tc.assert(t, string(body))
		})
	}
}

func Test_ProcessURLHandler(t *testing.T) {
	t.Parallel()
	test := setupTestApp(t)
	defer test.teardown()

	tests := []struct {
		name           string
		payload        string
		expectedStatus int
		assert         func(t *testing.T, body string)
	}{
		{
			name:           "valid - canonical",
			payload:        `{"url":"https://example.com/path","operation":"canonical"}`,
			expectedStatus: http.StatusOK,
			assert: func(t *testing.T, body string) {
				if !strings.Contains(body, "example.com") {
					t.Errorf("expected processed url in body")
				}
			},
		},
		{
			name:           "valid - redirection",
			payload:        `{"url":"https://EXAMPLE.com/FOOD-Experience?query=abc","operation":"redirection"}`,
			expectedStatus: http.StatusOK,
			assert: func(t *testing.T, body string) {
				expected := `https://www.byfood.com/food-experience`
				if !strings.Contains(strings.ToLower(body), expected) {
					t.Errorf("expected %q in response body, got: %s", expected, body)
				}
			},
		},
		{
			name:           "invalid url",
			payload:        `{"url":"bad-url","operation":"canonical"}`,
			expectedStatus: http.StatusUnprocessableEntity,
			assert: func(t *testing.T, body string) {
				if !strings.Contains(body, "url") {
					t.Errorf("expected url validation error")
				}
			},
		},
		{
			name:           "invalid operation",
			payload:        `{"url":"https://example.com","operation":"compress"}`,
			expectedStatus: http.StatusUnprocessableEntity,
			assert: func(t *testing.T, body string) {
				if !strings.Contains(body, "operation") {
					t.Errorf("expected operation validation error")
				}
			},
		},
		{
			name:           "malformed json",
			payload:        `{url:"bad"}`,
			expectedStatus: http.StatusBadRequest,
			assert: func(t *testing.T, body string) {
				if !strings.Contains(body, "badly-formed") {
					t.Errorf("expected JSON error")
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/url/process", bytes.NewReader([]byte(tc.payload)))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			test.handler.ServeHTTP(w, r)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tc.expectedStatus {
				t.Errorf("got %d, want %d", res.StatusCode, tc.expectedStatus)
			}

			body, _ := io.ReadAll(res.Body)
			tc.assert(t, string(body))
		})
	}
}

func testCreateBook(test *testApp, bk book.NewBook) (*book.Book, error) {
	payload := fmt.Sprintf(`{"title":"%s","author":"%s","year":%d}`, bk.Title, bk.Author, bk.Year)
	r := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader([]byte(payload)))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	test.handler.ServeHTTP(w, r)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("expected status 201 Created, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	var created book.Book
	_ = json.Unmarshal(body, &created)

	return &created, nil
}
