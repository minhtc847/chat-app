package main

import (
	"BE-chat-app/internal/data"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
)

func (app *application) ServeConversationChat(hub *ChatHub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil) // Upgrade HTTP to WebSocket
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer conn.Close()

	conversationIdStr := r.URL.Query().Get("conversation_id")
	userIdStr := r.URL.Query().Get("user_id")

	conversationId, err := uuid.Parse(conversationIdStr)
	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("invalid conversation_id"))
		return
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("invalid user_id"))
		return
	}

	// Add user to the conversation
	hub.AddConnection(conversationId, userId, conn)
	defer hub.RemoveConnection(conversationId, userId)

	for {
		var input struct {
			Content  string `json:"content"`
			Type     string `json:"type"`
			SenderId string `json:"sender_id"`
		}

		if err := conn.ReadJSON(&input); err != nil {
			app.badRequestResponse(w, r, err)
			break
		}
		SenderProfileId, err := uuid.Parse(input.SenderId)
		if err != nil {
			app.badRequestResponse(w, r, err)
			break
		}

		// Save the message to the database
		message := &data.Message{
			ConversationId:  conversationId,
			SenderProfileId: SenderProfileId,
			Content:         input.Content,
			Type:            input.Type,
		}
		err = app.models.Messages.Insert(message)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		// Broadcast the message to all participants in the conversation
		connections := hub.GetConnections(conversationId)
		for _, recipientConn := range connections {
			err = recipientConn.WriteJSON(envelope{"message": message})
			if err != nil {
				log.Printf("Error broadcasting message: %v", err)
			}
		}
	}
}
