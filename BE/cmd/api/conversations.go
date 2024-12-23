package main

import (
	"BE-chat-app/internal/data"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

func (app *application) getConversationsHandler(w http.ResponseWriter, r *http.Request) {
	conservationStr := app.readStringParam(r, "id")
	if conservationStr == "" {
		app.notFoundResponse(w, r)
	}
	conservationId := uuid.MustParse(conservationStr)
	conversation, err := app.models.Conversations.Get(conservationId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"conversation": conversation}, nil)
	if err != nil {
	}
}

func (app *application) createConversationsHandler(w http.ResponseWriter, r *http.Request) {
	receiverIdStr := app.readString(r.URL.Query(), "receiver_id", "")
	receiverId, err := uuid.Parse(receiverIdStr)
	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("invalid receiver_id"))
		return
	}
	senderId := app.contextGetUser(r).ID
	newConservation := &data.Conversations{
		ProfileOneId: senderId,
		ProfileTwoId: receiverId,
	}
	err = app.models.Conversations.Insert(newConservation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	newConservation, err = app.models.Conversations.GetByProfiles(senderId, receiverId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"conversation": newConservation}, nil)
	if err != nil {
	}
}
func (app *application) deleteConversationsHandler(w http.ResponseWriter, r *http.Request) {

	conservationId := uuid.MustParse(app.readStringParam(r, "id"))

	conservation, err := app.models.Conversations.Get(conservationId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if conservation.ProfileOneId != app.contextGetUser(r).ID && conservation.ProfileTwoId != app.contextGetUser(r).ID {
		app.notFoundResponse(w, r)
		return
	}
	err = app.models.Conversations.Delete(conservationId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "delete successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
