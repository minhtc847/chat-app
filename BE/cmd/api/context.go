package main

import (
	"BE-chat-app/internal/data"
	"context"
	"net/http"
)

type contextKey string

const userContextKey = contextKey("user")

func (app *application) contextSetUser(r *http.Request, user *data.Profile) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) *data.Profile {
	user, ok := r.Context().Value(userContextKey).(*data.Profile)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}
