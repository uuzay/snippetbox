package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	middleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	sessionMiddleware := alice.New(app.session.Enable)

	mux := pat.New()
	mux.Get("/", sessionMiddleware.ThenFunc(app.home))
	mux.Get("/snippet/create", sessionMiddleware.ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", sessionMiddleware.ThenFunc(app.createSnippet))

	// ordering matters for pat
	// if /snippet/create was after this pattern, it would still get matched to this one
	mux.Get("/snippet/:id", sessionMiddleware.ThenFunc(app.showSnippet))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", neuter(fileServer)))

	return middleware.Then(mux)
}
