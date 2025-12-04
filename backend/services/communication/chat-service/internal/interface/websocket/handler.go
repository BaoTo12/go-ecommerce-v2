package websocket

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/titan-commerce/backend/chat-service/internal/application"
	"github.com/titan-commerce/backend/chat-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, validate origin
	},
}

type WebSocketMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type SendMessagePayload struct {
	ConversationID string `json:"conversation_id"`
	Content        string `json:"content"`
	MessageType    string `json:"message_type"`
}

type ChatWebSocketHandler struct {
	service *application.ChatService
	logger  *logger.Logger
	clients map[*websocket.Conn]string
	mu      sync.RWMutex
}

func NewChatWebSocketHandler(service *application.ChatService, logger *logger.Logger) *ChatWebSocketHandler {
	return &ChatWebSocketHandler{
		service: service,
		logger:  logger,
		clients: make(map[*websocket.Conn]string),
	}
}

func (h *ChatWebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error(err, "failed to upgrade connection")
		return
	}

	connID := uuid.New().String()
	msgChan := make(chan *domain.Message, 100)

	h.mu.Lock()
	h.clients[conn] = userID
	h.mu.Unlock()

	// Register connection
	h.service.GetConnectionManager().AddConnection(userID, connID, msgChan)

	defer func() {
		h.service.GetConnectionManager().RemoveConnection(userID, connID)
		h.mu.Lock()
		delete(h.clients, conn)
		h.mu.Unlock()
		conn.Close()
	}()

	// Goroutine to send messages
	go func() {
		for msg := range msgChan {
			data, _ := json.Marshal(map[string]interface{}{
				"type":    "new_message",
				"message": msg,
			})
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}
		}
	}()

	// Read messages
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				h.logger.Error(err, "websocket error")
			}
			break
		}

		var wsMsg WebSocketMessage
		if err := json.Unmarshal(data, &wsMsg); err != nil {
			continue
		}

		switch wsMsg.Type {
		case "send_message":
			var payload SendMessagePayload
			if err := json.Unmarshal(wsMsg.Payload, &payload); err != nil {
				continue
			}

			msgType := domain.MessageTypeText
			if payload.MessageType == "image" {
				msgType = domain.MessageTypeImage
			}

			msg, err := h.service.SendMessage(r.Context(), payload.ConversationID, userID, payload.Content, msgType)
			if err != nil {
				h.logger.Error(err, "failed to send message")
				continue
			}

			// Send confirmation to sender
			response, _ := json.Marshal(map[string]interface{}{
				"type":    "message_sent",
				"message": msg,
			})
			conn.WriteMessage(websocket.TextMessage, response)

		case "typing":
			// Broadcast typing indicator to other participants
			h.logger.Debugf("User %s is typing", userID)

		case "mark_read":
			var payload struct {
				ConversationID string `json:"conversation_id"`
			}
			json.Unmarshal(wsMsg.Payload, &payload)
			h.service.MarkAsRead(r.Context(), payload.ConversationID, userID)
		}
	}
}
