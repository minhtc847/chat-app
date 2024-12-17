package main

import (
	"BE-chat-app/internal/data"
	"BE-chat-app/internal/validator"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/pascaldekloe/jwt"
	"net/http"
	"time"
)

// login
func (app *application) createAuthenticationJWTTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	v := validator.New()
	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlaintext(v, input.Password)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	profile, err := app.models.Profiles.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	match, err := profile.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}
	if !profile.Activated {
		app.inactiveAccountResponse(w, r)
		return
	}
	// Create a JWT claims struct containing the user ID as the subject, with an issued
	// time of now and validity window of the next 24 hours. We also set the issuer and
	// audience to a unique identifier for our application.
	var claims jwt.Claims
	claims.Subject = profile.ID.String()
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	claims.Issuer = "chatappbe.minhtc47.net"
	claims.Audiences = []string{"chatappfe.minhtc47.net"}
	// Sign the JWT claims using the HMAC-SHA256 algorithm and the secret key from the
	// application config. This returns a []byte slice containing the JWT as a base64
	// encoded string.
	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.jwt.secret))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Convert the []byte slice to a string and return it in a JSON response.
	err = app.writeJSON(w, http.StatusCreated, envelope{"authentication_token": string(jwtBytes)}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func GenerateToken() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err) // Handle error appropriately in production
	}
	return hex.EncodeToString(b)
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {

	token := app.contextGetToken(r)
	userId := app.contextGetUser(r).ID.String()
	err := app.redis.storeLogoutToken(*token, userId, 24*time.Hour)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "logout successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
