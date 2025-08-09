package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	mux := httprouter.New()

	mux.NotFound = http.HandlerFunc(app.notFound)
	mux.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowed)

	mux.HandlerFunc("GET", "/status", app.status)

	mux.HandlerFunc("GET", "/books", app.listBooksHandler)
	mux.HandlerFunc("POST", "/books", app.createBookHandler)
	mux.HandlerFunc("GET", "/books/:id", app.showBookHandler)
	mux.HandlerFunc("PUT", "/books/:id", app.updateBookHandler)
	mux.HandlerFunc("DELETE", "/books/:id", app.deleteBookHandler)

	return app.logAccess(app.recoverPanic(mux))
}
