package main

import (
	"expvar"
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

	router = app.metricRoutes(router)
	return app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router))))
}

func (app *application) movieRoutes(router *httprouter.Router) *httprouter.Router {

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.requirePermission("movies:write", app.createMovieHandler))
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.requirePermission("movies:read", app.listMovieHandler))
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.requirePermission("movies:read", app.showMovieHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.requirePermission("movies:write", app.updateMovieHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.requirePermission("movies:write", app.deleteMovieHandler))

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

func (app *application) metricRoutes(router *httprouter.Router) *httprouter.Router {
		
	router.Handler(http.MethodGet, "/v1/metrics", expvar.Handler())
	return router
}
