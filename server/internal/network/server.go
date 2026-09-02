// internal/network/server.go

package network

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"royaka/internal/network/dto"
	"royaka/internal/network/wsconn"
	"strings"

	"royaka/internal/game"
	"royaka/internal/model"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return isOriginAllowed(r.Header.Get("Origin"), os.Getenv("ALLOWED_ORIGINS"))
	},
}

func isOriginAllowed(origin, configuredOrigins string) bool {
	origin = strings.TrimRight(strings.TrimSpace(origin), "/")
	if origin == "" {
		return false
	}

	if strings.TrimSpace(configuredOrigins) == "" {
		configuredOrigins = "http://localhost:5173,http://127.0.0.1:5173"
	}
	for _, allowed := range strings.Split(configuredOrigins, ",") {
		if strings.TrimRight(strings.TrimSpace(allowed), "/") == origin {
			return true
		}
	}
	return false
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[ERROR][WS] Upgrade failed: %v", err)
		http.Error(w, "WebSocket upgrade failed", http.StatusInternalServerError)
		return
	}
	wsconn.Register(conn)
	wsconn.ConfigureReader(conn)

	// Recover panic inside the goroutine safely
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ERROR][WS] Panic recovered: %v", r)
		}
		if player := model.GetPlayerByConn(conn); player != nil {
			game.RemovePlayerFromQueue(player)
			game.CleanupUser(player.User.Username)
		}
		game.HandleDisconnect(conn)
		model.RemoveConnection(conn)
		removeIdentity(conn)
		wsconn.Close(conn)
		log.Println("[WS] Connection closed")
	}()

	log.Println("[WS] WebSocket connection established")

	for {

		if !readAndProcessMessage(conn) {
			break
		}
	}
}

func readAndProcessMessage(conn *websocket.Conn) bool {
	_, msg, err := conn.ReadMessage()
	if err != nil {
		logWebSocketError(err)
		return false
	}

	var pdu dto.Envelope
	if err := json.Unmarshal(msg, &pdu); err != nil {
		log.Printf("[WARN][WS] Invalid JSON: %v", err)
		sendError(conn, "", "invalid_message", "Invalid message format")
		return true
	}
	if !pdu.Type.IsRequest() {
		log.Printf("[WARN][WS] Unsupported request type: %s", pdu.Type)
		sendError(conn, pdu.RequestID, "unknown_message_type", "Unknown message type")
		return true
	}

	log.Printf("[INFO][WS] Message type: %s", pdu.Type)
	processMessage(conn, pdu)
	return true
}

func processMessage(conn *websocket.Conn, pdu dto.Envelope) {
	switch pdu.Type {
	case dto.MessageRegister:
		handleRegister(conn, pdu.RequestID, pdu.Data)
		return
	case dto.MessageLogin:
		handleLogin(conn, pdu.RequestID, pdu.Data)
		return
	case dto.MessageAuthenticate:
		handleAuthenticate(conn, pdu.RequestID, pdu.Data)
		return
	case dto.MessageLogout:
		handleLogout(conn, pdu.RequestID, pdu.Data)
		return
	case dto.MessageGetUser:
		handleGetUser(conn, pdu.RequestID, pdu.Data)
		return
	}

	identity, authenticated := getIdentity(conn)
	if !authenticated {
		sendError(conn, pdu.RequestID, "authentication_required", "Authentication required")
		return
	}
	if len(pdu.Data) > 0 {
		var envelope struct {
			Username string `json:"username"`
		}
		if json.Unmarshal(pdu.Data, &envelope) == nil && envelope.Username != "" &&
			strings.TrimSpace(envelope.Username) != identity.Username {
			log.Printf("[WARN][AUTH] Rejected username spoof on %s", pdu.Type)
			sendError(conn, pdu.RequestID, "username_mismatch", "Username does not match authenticated user")
			return
		}
	}

	switch pdu.Type {
	case dto.MessageGetDesk:
		game.HandleGetDesk(conn, pdu.RequestID, pdu.Data)
	case dto.MessageFindMatch:
		game.HandleFindMatch(conn, pdu.RequestID, pdu.Data)
	case dto.MessageGetGame:
		game.HandleGetGame(conn, pdu.RequestID, pdu.Data)
	case dto.MessageAttack:
		game.HandleAttack(conn, pdu.RequestID, pdu.Data)
	case dto.MessageHeal:
		game.HandleHeal(conn, pdu.RequestID, pdu.Data)
	case dto.MessageSkipTurn:
		game.HandleSkipTurn(conn, pdu.RequestID, pdu.Data)
	case dto.MessagePlayAgain:
		game.HandlePlayAgain(conn, pdu.RequestID, pdu.Data)
	case dto.MessageLeaveGame:
		game.HandleLeaveGame(conn, pdu.RequestID, pdu.Data)
	case dto.MessageSelectTroop:
		game.HandleSelectTroop(conn, pdu.RequestID, pdu.Data)
	default:
		log.Printf("[WARN][WS] Unknown message type: %s", pdu.Type)
		sendError(conn, pdu.RequestID, "unknown_message_type", "Unknown message type")
	}
}

func sendError(conn *websocket.Conn, requestID, code, message string) {
	err := wsconn.Send(conn, dto.Fail(dto.MessageError, requestID, code, message))
	if err != nil {
		log.Printf("[ERROR][WS] Failed to send error response: %v", err)
	}
}

func logWebSocketError(err error) {
	if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
		log.Printf("[ERROR][WS] Unexpected closure: %v", err)
	} else {
		log.Printf("[WARN][WS] Client disconnected: %v", err)
	}
}
