package game

import (
	"encoding/json"
	"log"
	"royaka/internal/model"
	"royaka/internal/network/dto"

	"github.com/gorilla/websocket"
)

func HandleLeaveGame(conn *websocket.Conn, requestID string, data json.RawMessage) {
	var req dto.GameRequest

	if err := json.Unmarshal(data, &req); err != nil || req.RoomID == "" || req.Username == "" {
		writeToConnection(conn, dto.Fail(dto.MessageLeaveGameResponse, requestID, "invalid_request", "Invalid request"))
		return
	}

	roomsMu.RLock()
	room, found := rooms[req.RoomID]
	roomsMu.RUnlock()
	if !found {
		writeToConnection(conn, dto.Fail(dto.MessageLeaveGameResponse, requestID, "room_not_found", "Room not found"))
		return
	}

	room.mu.Lock()
	defer room.mu.Unlock()

	player1 := room.Game.Player1
	player2 := room.Game.Player2

	var winner *model.Player

	if player1 != nil && player1.User.Username == req.Username {
		winner = player2
	} else if player2 != nil && player2.User.Username == req.Username {
		winner = player1
	}

	if winner != nil {
		room.Game.SetWinner(winner)

		payload := dto.Push(dto.MessageGameOverResponse, "", dto.GameOver{Winner: dto.ToPlayer(winner)})

		sendToClient(winner.User.Username, payload)
	}

	if room.Game.TurnTimerCancel != nil {
		room.Game.TurnTimerCancel()
	}

	writeToConnection(conn, dto.OK(dto.MessageLeaveGameResponse, requestID, "Left room and winner set if applicable", dto.Empty{}))
}

func HandleDisconnect(conn *websocket.Conn) {
	player := model.GetPlayerByConn(conn)
	if player == nil {
		return
	}

	username := player.User.Username
	roomID := GetRoomIDByUsername(username)
	if roomID == "" {
		return
	}

	log.Printf("[INFO] %s disconnected, handling leave...", username)

	req := dto.GameRequest{
		RoomID:   roomID,
		Username: username,
	}
	raw, _ := json.Marshal(req)
	HandleLeaveGame(conn, "", raw)
}
