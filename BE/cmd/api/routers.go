package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodPost, "/users", app.registerUserHandler)

	router.HandlerFunc(http.MethodGet, "/friends", app.getFriendsHandler)
	router.HandlerFunc(http.MethodPost, "/friends", app.sendInviteHandler)
	router.HandlerFunc(http.MethodPut, "/friends", app.confirmInviteHandler)

	return app.enableCORS(app.authenticate(router))

}
