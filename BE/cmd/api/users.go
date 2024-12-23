package main

import (
	"BE-chat-app/internal/data"
	"BE-chat-app/internal/validator"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Gender   bool   `json:"gender"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	user := &data.Profile{
		Name:      input.Name,
		Email:     input.Email,
		Gender:    input.Gender,
		Activated: false,
	}
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	v := validator.New()
	if data.ValidateProfile(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Profiles.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Generate activation token
	token := GenerateToken()
	redisKey := fmt.Sprintf("activation:%s", token)
	redisExpiration := 24 * 3 * time.Hour // Token expiration time

	// Store token in Redis
	err = app.redis.storeActivationToken(redisKey, user.ID.String(), redisExpiration)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.background(func() {
		data := map[string]any{
			"activationToken": token,
			"userID":          user.ID,
		}

		err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			// Importantly, if there is an error sending the email then we use the
			// app.logger.PrintError() helper to manage it, instead of the
			// app.serverErrorResponse() helper like before.
			app.logger.PrintError(err, nil)
		}
	})
	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {

	qs := r.URL.Query()
	token := app.readString(qs, "token", "")

	if app.redis.ActivateToken["activation:"+token] == "" {

		app.serverErrorResponse(w, r, fmt.Errorf("invalid activation token"))

	}
	user, err := app.models.Profiles.Get(uuid.MustParse(app.redis.getActivationToken("activation:" + token)))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	user.Activated = true

	err = app.models.Profiles.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// If everything went successfully, then we delete all activation tokens for the
	// user.
	err = app.redis.removeActivateToken("activation:" + token)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Send the updated user details to the client in a JSON response.
	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// create login handler
//func (app *application) loginUserHandler(w http.ResponseWriter, r *http.Request) {
//	var input struct {
//		Email    string `json:"email"`
//		Password string `json:"password"`
//	}
//	err := app.readJSON(w, r, &input)
//	if err != nil {
//		app.badRequestResponse(w, r, err)
//		return
//	}
//	v := validator.New()
//	if data.ValidateEmail(v, input.Email); !v.Valid() {
//		app.failedValidationResponse(w, r, v.Errors)
//		return
//	}
//	user, err := app.models.Profiles.GetByEmail(input.Email)
//	if err != nil {
//		switch {
//		case errors.Is(err, data.ErrRecordNotFound):
//			app.invalidCredentialsResponse(w, r)
//		default:
//			app.serverErrorResponse(w, r, err)
//		}
//		return
//	}
//	match, err := user.Password.Matches(input.Password)
//	if err != nil {
//		app.serverErrorResponse(w, r, err)
//		return
//	}
//	if !match {
//		app.invalidCredentialsResponse(w, r)
//		return
//	}
//	if !user.Activated {
//		app.inactiveAccountResponse(w, r)
//		return
//	}
//	token, err := app.createAuthenticationJWTTokenHandler(user.ID, 3*24*time.Hour, data.ScopeAuthentication)
//	if err != nil {
//		app.serverErrorResponse(w, r, err)
//		return
//	}
//	err = app.writeJSON(w, http.StatusOK, envelope{
//		"authenticationToken": token,
//		"user":                user,
//	}, nil)
//	if err != nil {
//		app.serverErrorResponse(w, r, err)
//	}
//}
