package main

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	dynamicMiddleware := alice.New(app.session.Enable, noSurf)
	//mux := pat.New()
	r := mux.NewRouter()

	r.Handle("/", dynamicMiddleware.ThenFunc(app.home))

	r.Handle("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippet)).Methods(http.MethodPost)
	r.Handle("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippetForm)).Methods(http.MethodGet)
	r.Handle("/snippet/{id}", dynamicMiddleware.ThenFunc(app.showSnippet)).Methods(http.MethodGet)
	r.Handle("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser)).Methods(http.MethodPost)
	r.Handle("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm)).Methods(http.MethodGet)
	r.Handle("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm)).Methods(http.MethodPost)
	r.Handle("/user/login", dynamicMiddleware.ThenFunc(app.loginUser)).Methods(http.MethodGet)
	r.Handle("/user/logout", dynamicMiddleware.ThenFunc(app.logoutUser)).Methods(http.MethodPost)

	//mux.Get("/snippet/:id", http.HandlerFunc(app.showSnippet))
	//mux.Get("/", http.HandlerFunc(app.home))
	//mux.Post("/snippet/:id", http.HandlerFunc(app.createSnippet))

	fileServer := http.FileServer(http.Dir("./ui/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	http.Handle("/", r)

	return standardMiddleware.Then(nil)
}
