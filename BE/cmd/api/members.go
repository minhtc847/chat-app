package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
)

func (app *application) addMembersToGroupHandler(w http.ResponseWriter, r *http.Request) {
	channelID, err := app.readUUIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var input struct {
		UserIDs []uuid.UUID `json:"user_ids"`
	}
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	exist, err := app.models.Channel.ExistsChannel(channelID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !exist {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Member.AddMembers(channelID, input.UserIDs)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	response := map[string]interface{}{
		"message": "Members added successfully",
	}

	app.writeJSON(w, http.StatusOK, response, nil)
}
