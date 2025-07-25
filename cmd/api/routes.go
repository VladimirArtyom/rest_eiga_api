package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	var router *httprouter.Router = httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router = app.movieRoutes(router)
	router = app.userRoutes(router)
	router = app.userTokens(router)

	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}

func (app *application) movieRoutes(router *httprouter.Router) *httprouter.Router {

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.requireActivatedUser(app.createMovieHandler))
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.requireActivatedUser(app.listMovieHandler))
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.requireActivatedUser(app.showMovieHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.requireActivatedUser(app.updateMovieHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.requireActivatedUser(app.deleteMovieHandler))

	return router
}

func (app *application) userRoutes(router *httprouter.Router) *httprouter.Router {
	
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	return router
	
}

func (app *application) userTokens(router *httprouter.Router) *httprouter.Router {
	router.HandlerFunc(http.MethodPost,"/v1/tokens/authentication", app.createAuthenticationTokenHandler)
	return router
}
