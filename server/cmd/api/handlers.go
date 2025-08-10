package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Babatunde50/book-crud/server/business/book"
	"github.com/Babatunde50/book-crud/server/business/urlprocessor"
	"github.com/Babatunde50/book-crud/server/internal/request"
	"github.com/Babatunde50/book-crud/server/internal/response"
	"github.com/Babatunde50/book-crud/server/internal/validator"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

func (app *application) status(w http.ResponseWriter, r *http.Request) {
	_ = response.JSON(w, http.StatusOK, map[string]string{"Status": "OK"})
}

func validateBookRequest(br NewBookRequest) validator.Validator {
	var v validator.Validator

	// Title validations
	v.CheckField(br.Title != "", "title", "title is required")
	v.CheckField(len(br.Title) <= 100, "title", "title must not exceed 100 characters")

	// Author validations
	v.CheckField(br.Author != "", "author", "author is required")
	v.CheckField(len(br.Author) <= 50, "author", "author must not exceed 50 characters")

	// Year validations
	currentYear := time.Now().Year()
	v.CheckField(br.Year >= 1, "year", "year must be a positive number")
	v.CheckField(br.Year <= currentYear, "year", fmt.Sprintf("year cannot be in the future (max %d)", currentYear))

	return v
}

// @Summary Create a new book
// @Description Adds a book to the database
// @Accept json
// @Tags         books
// @Produce json
// @Param book body NewBookRequest true "Book data"
// @Success 201 {object} book.Book
// @Failure 400 {object} map[string]string
// @Failure 422 {object} map[string]interface{}
// @Router /books [post]
func (app *application) createBookHandler(w http.ResponseWriter, r *http.Request) {
	var input NewBookRequest

	err := request.DecodeJSON(w, r, &input)

	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	v := validateBookRequest(input)
	if v.HasErrors() {
		app.failedValidation(w, r, v)
		return
	}

	newBook := book.NewBook{
		Title:  input.Title,
		Author: input.Author,
		Year:   input.Year,
	}

	bk, err := app.bookCore.Create(r.Context(), newBook)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = response.JSON(w, http.StatusCreated, toBookResponse(bk))

	if err != nil {
		app.serverError(w, r, err)
		return
	}
}

// @Summary      List all books
// @Tags         books
// @Produce      json
// @Success      200 {array} BookResponse
// @Failure      500 {object} map[string]string
// @Router       /books [get]
func (app *application) listBooksHandler(w http.ResponseWriter, r *http.Request) {
	books, err := app.bookCore.QueryAll(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = response.JSON(w, http.StatusOK, toBooksResponse(books))
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}

// @Summary      Get a book by ID
// @Tags         books
// @Produce      json
// @Param        id path string true "Book ID (UUID)"
// @Success      200 {object} BookResponse
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /books/{id} [get]
func (app *application) showBookHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	bk, err := app.bookCore.QueryByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, book.ErrNotFound):
			app.notFound(w, r)
		default:
			app.serverError(w, r, err)
		}
		return
	}

	err = response.JSON(w, http.StatusOK, toBookResponse(bk))
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}

func validateUpdateBookRequest(br UpdateBookRequest) validator.Validator {
	var v validator.Validator

	if br.Title != nil {
		v.CheckField(*br.Title != "", "title", "must be provided")
		v.CheckField(len(*br.Title) <= 100, "title", "must not be more than 100 characters long")
	}

	if br.Author != nil {
		v.CheckField(*br.Author != "", "author", "must be provided")
		v.CheckField(len(*br.Author) <= 50, "author", "must not be more than 50 characters long")
	}

	if br.Year != nil {
		v.CheckField(*br.Year >= 1, "year", "must be a valid year")
		v.CheckField(*br.Year <= time.Now().Year(), "year", "must not be in the future")
	}

	return v
}

// @Summary      Update a book by ID
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id   path string              true "Book ID (UUID)"
// @Param        book body UpdateBookRequest   true "Partial fields to update"
// @Success      200 {object} BookResponse
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      422 {object} validator.Validator
// @Failure      500 {object} map[string]string
// @Router       /books/{id} [put]
func (app *application) updateBookHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	var input UpdateBookRequest
	err = request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	v := validateUpdateBookRequest(input)
	if v.HasErrors() {
		app.failedValidation(w, r, v)
		return
	}

	// Step 1: Fetch the existing book
	existing, err := app.bookCore.QueryByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, book.ErrNotFound):
			app.notFound(w, r)
		default:
			app.serverError(w, r, err)
		}
		return
	}

	// Step 2: Prepare the update
	updates := book.UpdateBook{}

	if input.Title != nil {
		updates.Title = input.Title
	}
	if input.Author != nil {
		updates.Author = input.Author
	}
	if input.Year != nil {
		updates.Year = input.Year
	}

	// Step 3: Apply update
	updated, err := app.bookCore.Update(r.Context(), existing, updates)
	if err != nil {
		switch {
		case errors.Is(err, book.ErrNotFound):
			app.notFound(w, r)
		default:
			app.serverError(w, r, err)
		}
		return
	}

	err = response.JSON(w, http.StatusOK, toBookResponse(updated))
	if err != nil {
		app.serverError(w, r, err)
	}
}

// @Summary      Delete a book by ID
// @Tags         books
// @Param        id path string true "Book ID (UUID)"
// @Success      204
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /books/{id} [delete]
func (app *application) deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	err = app.bookCore.Delete(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, book.ErrNotFound):
			app.notFound(w, r)
		default:
			app.serverError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func validateURLRequest(input URLRequest) validator.Validator {
	var v validator.Validator

	// URL must not be empty
	v.CheckField(input.URL != "", "url", "must be provided")

	// Validate URL structure
	if _, err := url.ParseRequestURI(input.URL); err != nil {
		v.AddFieldError("url", "must be a valid URL")
	}

	// Normalize and validate operation
	op := strings.ToLower(input.Operation)
	validOps := map[string]bool{
		"canonical":   true,
		"redirection": true,
		"all":         true,
	}
	v.CheckField(validOps[op], "operation", "must be one of 'canonical', 'redirection', or 'all'")

	return v
}

// @Summary      Process a URL
// @Description  Canonicalize or redirect-normalize a URL
// @Tags         url
// @Accept       json
// @Produce      json
// @Param        payload body URLRequest true "URL payload"
// @Success      200 {object} URLResponse
// @Failure      400 {object} map[string]string
// @Failure      422 {object} validator.Validator
// @Failure      500 {object} map[string]string
// @Router       /url/process [post]
func (app *application) processURLHandler(w http.ResponseWriter, r *http.Request) {
	var input URLRequest

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	v := validateURLRequest(input)
	if v.HasErrors() {
		app.failedValidation(w, r, v)
		return
	}

	result, err := app.urlProcessorCore.Process(input.URL, input.Operation)
	if err != nil {
		switch {
		case errors.Is(err, urlprocessor.ErrInvalidOperation):
			app.badRequest(w, r, err)
		default:
			app.serverError(w, r, err)
		}
		return
	}

	resp := URLResponse{ProcessedURL: result}
	err = response.JSON(w, http.StatusOK, resp)
	if err != nil {
		app.serverError(w, r, err)
	}
}
