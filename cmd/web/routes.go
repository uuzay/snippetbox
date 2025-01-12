package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	middleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	sessionMiddleware := alice.New(app.session.Enable, noSurf)

	mux := pat.New()
	mux.Get("/", sessionMiddleware.ThenFunc(app.home))
	mux.Get("/snippet/create", sessionMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", sessionMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippet))

	// ordering matters for pat
	// if /snippet/create was after this pattern, it would still get matched to this one
	mux.Get("/snippet/:id", sessionMiddleware.ThenFunc(app.showSnippet))

	mux.Get("/user/signup", sessionMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", sessionMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", sessionMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", sessionMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", sessionMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", neuter(fileServer)))

	return middleware.Then(mux)
}
