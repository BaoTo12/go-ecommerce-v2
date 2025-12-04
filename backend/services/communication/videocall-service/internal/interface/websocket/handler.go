package websocket

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/titan-commerce/backend/videocall-service/internal/application"
	"github.com/titan-commerce/backend/videocall-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type SignalingHandler struct {
	service *application.VideoCallService
	logger  *logger.Logger
}

func NewSignalingHandler(service *application.VideoCallService, logger *logger.Logger) *SignalingHandler {
	return &SignalingHandler{
		service: service,
		logger:  logger,
	}
}

func (h *SignalingHandler) HandleSignaling(w http.ResponseWriter, r *http.Request) {
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

	// Register for signaling messages
	sigChan := h.service.GetSignalingHub().Register(userID)

	defer func() {
		h.service.GetSignalingHub().Unregister(userID)
		conn.Close()
	}()

	// Goroutine to send signaling messages
	go func() {
		for msg := range sigChan {
			data, _ := json.Marshal(msg)
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}
		}
	}()

	// Read signaling messages
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				h.logger.Error(err, "websocket error")
			}
			break
		}

		var msg domain.SignalingMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}

		msg.FromUser = userID

		switch msg.Type {
		case "call":
			// Initiate call
			callType := domain.CallTypeVideo
			if t, ok := msg.Payload["call_type"].(string); ok && t == "audio" {
				callType = domain.CallTypeAudio
			}
			h.service.InitiateCall(r.Context(), userID, msg.ToUser, callType)

		case "accept":
			h.service.AcceptCall(r.Context(), msg.CallID, userID)

		case "reject":
			h.service.RejectCall(r.Context(), msg.CallID, userID)

		case "end":
			h.service.EndCall(r.Context(), msg.CallID, userID)

		case "offer", "answer", "ice-candidate":
			// Relay WebRTC signaling
			h.service.RelaySignaling(r.Context(), &msg)
		}
	}
}
