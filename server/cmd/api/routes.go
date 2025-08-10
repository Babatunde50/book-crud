package main

import (
	"expvar"
	"net/http"

	_ "github.com/Babatunde50/book-crud/server/cmd/api/docs"
	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger"
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

	mux.HandlerFunc("POST", "/url/process", app.processURLHandler)

	mux.Handler(http.MethodGet, "/swagger/*any", httpSwagger.WrapHandler)

	mux.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	return app.logAccess(app.recoverPanic(mux))
}
