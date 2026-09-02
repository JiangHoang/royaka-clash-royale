package game

import (
	"encoding/json"
	"log"
	"royaka/internal/model"
	"royaka/internal/network/dto"

	"github.com/gorilla/websocket"
)

func HandleGetGame(conn *websocket.Conn, requestID string, data json.RawMessage) {
	var req dto.GameRequest

	// Parse & validate request
	if err := json.Unmarshal(data, &req); err != nil || req.RoomID == "" || req.Username == "" {
		log.Printf("[WARN][GAME] invalid request: %v", err)
		writeToConnection(conn, dto.Fail(dto.MessageGameResponse, requestID, "invalid_request", invalidRequestMessage))
		return
	}

	// Get room safely
	roomsMu.RLock()
	room, exists := rooms[req.RoomID]
	roomsMu.RUnlock()
	if !exists {
		log.Printf("[WARN][GAME] room %s not found for user %s", req.RoomID, req.Username)
		writeToConnection(conn, dto.Fail(dto.MessageGameResponse, requestID, "room_not_found", roomRequestMessage))
		return
	}

	// Identify current player and opponent
	var currentUser, opponent *model.Player
	if room.Player1.User.Username == req.Username {
		currentUser, opponent = room.Player1, room.Player2
	} else if room.Player2.User.Username == req.Username {
		currentUser, opponent = room.Player2, room.Player1
	} else {
		log.Printf("[WARN][GAME] user %s not in room %s", req.Username, req.RoomID)
		writeToConnection(conn, dto.Fail(dto.MessageGameResponse, requestID, "player_not_in_room", "Player not in room"))
		return
	}

	dataPayload := dto.GameData{User: dto.ToPlayer(currentUser), Opponent: dto.ToPlayer(opponent)}

	if room.Game.Enhanced {
		dataPayload.Player1 = room.Player1.User.Username
		dataPayload.Map = toBattleEntities(room.Game.BattleSystem.GetEntityList())
		dataPayload.Time = room.Game.MaxTime.Milliseconds()
	} else {
		dataPayload.Turn = room.Game.Turn
	}

	writeToConnection(conn, dto.OK(dto.MessageGameResponse, requestID, "Game info loaded", dataPayload))

	log.Printf("[INFO][GAME] sent game state to %s in room %s", req.Username, req.RoomID)
}
