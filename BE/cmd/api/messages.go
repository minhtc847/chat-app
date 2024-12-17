package main

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (app *application) getMessagesByConversation(w http.ResponseWriter, r *http.Request) {

	conservationId := uuid.MustParse(app.readStringParam(r, "conversation_id"))
	timestampStr := app.readString(r.URL.Query(), "timestamp", "")
	timestamp, err := time.Parse(time.RFC3339, timestampStr)
	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("invalid timestamp"))
		return
	}

	messages, err := app.models.Messages.GetMessagesAfterTime(conservationId, timestamp)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"messages": messages, "count": len(messages), "timestamp": timestamp}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateMessageHandler(w http.ResponseWriter, r *http.Request) {

	messageId := uuid.MustParse(app.readStringParam(r, "id"))
	message, err := app.models.Messages.Get(messageId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if message.Deleted {
		app.notFoundResponse(w, r)
		return
	}
	if message.SenderProfileId != app.contextGetUser(r).ID {
		app.badRequestResponse(w, r, fmt.Errorf("you are not the sender of this message"))
	}
	var input struct {
		Content string `json:"content"`
		Deleted bool   `json:"deleted"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	message.Content = input.Content
	message.Deleted = input.Deleted
	err = app.models.Messages.Update(message)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "update successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteMessageHandler(w http.ResponseWriter, r *http.Request) {

	messageId := uuid.MustParse(app.readStringParam(r, "id"))
	message, err := app.models.Messages.Get(messageId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if message.Deleted {
		app.notFoundResponse(w, r)
		return
	}
	if message.SenderProfileId != app.contextGetUser(r).ID {
		app.badRequestResponse(w, r, fmt.Errorf("you are not the sender of this message"))
		return
	}
	message.Deleted = true
	err = app.models.Messages.Update(message)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "delete successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
