package game

import (
	"encoding/json"
	"log"
	"royaka/internal/network/dto"

	"github.com/gorilla/websocket"
)

func HandleSkipTurn(conn *websocket.Conn, requestID string, data json.RawMessage) {
	var req dto.GameRequest
	if err := json.Unmarshal(data, &req); err != nil || req.RoomID == "" || req.Username == "" {
		log.Printf("[WARN][SKIP_TURN] Invalid request from conn %v: %v | Data: %s", conn.RemoteAddr(), err, string(data))
		writeToConnection(conn, dto.Fail(dto.MessageSkipTurnResponse, requestID, "invalid_request", "Invalid skip turn request"))
		return
	}

	roomsMu.RLock()
	room, exists := rooms[req.RoomID]
	roomsMu.RUnlock()

	if !exists {
		log.Printf("[WARN][SKIP_TURN] Room not found: %s by user %s", req.RoomID, req.Username)
		writeToConnection(conn, dto.Fail(dto.MessageSkipTurnResponse, requestID, "room_not_found", "Room not found"))
		return
	}
	if !room.Game.Started || room.Game.WinnerDeclared {
		writeToConnection(conn, dto.Fail(dto.MessageSkipTurnResponse, requestID, "game_finished", "Game has already ended"))
		return
	}

	current := room.Game.CurrentPlayer()
	if current.User.Username != req.Username {
		log.Printf("[WARN][SKIP_TURN] Not %s's turn in room %s", req.Username, req.RoomID)
		writeToConnection(conn, dto.Fail(dto.MessageSkipTurnResponse, requestID, "not_your_turn", "It's not your turn!"))
		return
	}

	room.Game.SkipTurn(current)

	log.Printf("[DEBUG][SKIP_TURN] Turn switched to: %s", room.Game.Turn)

	payload := dto.OK(dto.MessageSkipTurnResponse, requestID, "Turn skipped", dto.SkipTurnResult{Turn: room.Game.Turn, Player1: dto.ToPlayer(room.Game.Player1), Player2: dto.ToPlayer(room.Game.Player2)})

	sendToClient(room.Player1.User.Username, payload)
	sendToClient(room.Player2.User.Username, payload)
}
