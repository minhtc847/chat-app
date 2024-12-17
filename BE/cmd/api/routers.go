package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	hub := NewChatHub()

	router.GET("/ws/conversation", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		app.ServeConversationChat(hub, w, r)
	})
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodPost, "/users", app.registerUserHandler)

	router.HandlerFunc(http.MethodGet, "/users/activate", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/users/login", app.createAuthenticationJWTTokenHandler)
	router.HandlerFunc(http.MethodGet, "/users/logout", app.requireAuthenticatedUser(app.logout))
	router.HandlerFunc(http.MethodGet, "/conversation/:id/messages", app.requireAuthenticatedUser(app.getMessagesByConversation))
	router.HandlerFunc(http.MethodGet, "/conversation", app.requireAuthenticatedUser(app.getConversationsHandler))
	return app.enableCORS(app.authenticate(router))

}
