package main

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (app *application) getMessagesByConversation(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	conservationId, err := uuid.Parse(app.readString(qs, "conversation_id", ""))
	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("invalid conversation_id"))
		return
	}
	timestampStr := app.readString(qs, "timestamp", "")
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
