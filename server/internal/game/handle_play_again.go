package game

import (
	"encoding/json"
	"log"
	"royaka/internal/model"
	"royaka/internal/network/dto"

	"github.com/gorilla/websocket"
)

func HandlePlayAgain(conn *websocket.Conn, requestID string, data json.RawMessage) {
	var req dto.PlayAgainRequest

	if err := json.Unmarshal(data, &req); err != nil || req.RoomID == "" {
		log.Printf("[WARN][PLAY_AGAIN] Invalid request: %v", err)
		writeToConnection(conn, dto.Fail(dto.MessagePlayAgainResponse, requestID, "invalid_request", invalidRequestMessage))
		return
	}

	roomsMu.RLock()
	room, exists := rooms[req.RoomID]
	roomsMu.RUnlock()
	if !exists {
		log.Printf("[WARN][PLAY_AGAIN] Room %s not found", req.RoomID)
		writeToConnection(conn, dto.Fail(dto.MessagePlayAgainResponse, requestID, "room_not_found", roomRequestMessage))
		return
	}
	player := model.GetPlayerByConn(conn)
	if player == nil || (room.Player1.User.Username != player.User.Username && room.Player2.User.Username != player.User.Username) {
		writeToConnection(conn, dto.Fail(dto.MessagePlayAgainResponse, requestID, "player_not_in_room", "Not a member of this room"))
		return
	}

	roomsMu.Lock()
	delete(rooms, room.ID)
	roomsMu.Unlock()

	log.Printf("[INFO][PLAY_AGAIN] Room %s cleaned up", room.ID)
	writeToConnection(conn, dto.OK(dto.MessagePlayAgainResponse, requestID, "Ready for matchmaking", dto.Empty{}))
}
