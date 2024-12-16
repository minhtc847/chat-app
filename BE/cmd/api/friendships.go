package main

import (
	"BE-chat-app/internal/data"
	"github.com/google/uuid"
	"net/http"
)

func (app *application) getFriendsHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse("18e3afc2-88ca-4af8-8b5e-3768efc5202f")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	friends, err := app.models.Friends.GetAllFriends(userID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	response := map[string]interface{}{
		"friends": friends,
	}

	err = app.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) sendInviteHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse("e69d784c-71f2-453b-9331-bbb2871dbbae")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	profileID, err := uuid.Parse("18e3afc2-88ca-4af8-8b5e-3768efc5202f")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.models.Friends.SendInvite(userID, profileID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	response := map[string]interface{}{
		"message": "Invite to make friend successfully",
	}

	err = app.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) confirmInviteHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse("18e3afc2-88ca-4af8-8b5e-3768efc5202f")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	profileID, err := uuid.Parse("e69d784c-71f2-453b-9331-bbb2871dbbae")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	invite, err := app.models.Friends.GetInvite(profileID, userID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	status := "Accepted"

	var friendship data.Friendship
	friendship.ID = invite.ID
	friendship.Created_at = invite.Created_at

	err = app.models.Friends.ConfirmInvite(friendship.ID, status)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	response := map[string]interface{}{
		"message": "Invite friend is accepted or decline",
	}

	err = app.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
