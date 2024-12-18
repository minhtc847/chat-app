package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) createChannelHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name == "" || len(input.Name) > 15 {
		app.badRequestResponse(w, r, err)
		return
	}

	userID := app.contextGetUser(r).ID

	err = app.models.Channel.CreateChannel(input.Name, userID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	response := map[string]interface{}{
		"message": "Created channel success",
	}

	err = app.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
