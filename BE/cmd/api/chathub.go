package main

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections for simplicity; restrict in production
		return true
	},
}

type ChatHub struct {
	connections map[uuid.UUID]map[uuid.UUID]*websocket.Conn // conversationId -> map[userId -> websocket.Conn]
	mu          sync.Mutex                                  // Mutex for thread safety
}

func NewChatHub() *ChatHub {
	return &ChatHub{
		connections: make(map[uuid.UUID]map[uuid.UUID]*websocket.Conn),
	}
}

// Add a connection to a conversation
func (hub *ChatHub) AddConnection(conversationId, userId uuid.UUID, conn *websocket.Conn) {
	hub.mu.Lock()
	if hub.connections[conversationId] == nil {
		hub.connections[conversationId] = make(map[uuid.UUID]*websocket.Conn)
	}
	hub.connections[conversationId][userId] = conn
	hub.mu.Unlock()
}

// Remove a connection from a conversation
func (hub *ChatHub) RemoveConnection(conversationId, userId uuid.UUID) {
	hub.mu.Lock()
	if _, exists := hub.connections[conversationId]; exists {
		delete(hub.connections[conversationId], userId)
		if len(hub.connections[conversationId]) == 0 {
			delete(hub.connections, conversationId) // Clean up empty conversations
		}
	}
	hub.mu.Unlock()
}

// Get all connections in a conversation
func (hub *ChatHub) GetConnections(conversationId uuid.UUID) map[uuid.UUID]*websocket.Conn {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	return hub.connections[conversationId]
}
