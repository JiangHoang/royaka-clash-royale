package game

import (
	"encoding/json"
	"log"
	"royaka/internal/model"
	"royaka/internal/network/dto"

	"github.com/gorilla/websocket"
)

func HandleAttack(conn *websocket.Conn, requestID string, data json.RawMessage) {
	var req dto.AttackRequest

	// Parse & validate request data
	if err := json.Unmarshal(data, &req); err != nil || req.RoomID == "" || req.Username == "" || req.Troop == "" || req.Target == "" {
		log.Printf("[WARN][ATTACK] invalid request: %v", err)
		writeToConnection(conn, dto.Fail(dto.MessageAttackResponse, requestID, "invalid_request", invalidRequestMessage))
		return
	}

	// Fetch the room from memory
	roomsMu.RLock()
	room, exists := rooms[req.RoomID]
	roomsMu.RUnlock()
	if !exists {
		log.Printf("[WARN][ATTACK] Room %s not found for user %s", req.RoomID, req.Username)
		writeToConnection(conn, dto.Fail(dto.MessageAttackResponse, requestID, "room_not_found", roomRequestMessage))
		return
	}

	// Identify the attacker
	var attacker, defender *model.Player
	if room.Player1.User.Username == req.Username {
		attacker = room.Player1
		defender = room.Player2
	} else if room.Player2.User.Username == req.Username {
		attacker = room.Player2
		defender = room.Player1
	} else {
		log.Printf("[WARN][ATTACK] User %s not in room %s", req.Username, req.RoomID)
		writeToConnection(conn, dto.Fail(dto.MessageAttackResponse, requestID, "player_not_in_room", "You are not part of this match"))
		return
	}

	if room.Game.CurrentPlayer().User.Username != attacker.User.Username {
		writeToConnection(conn, dto.Fail(dto.MessageAttackResponse, requestID, "not_your_turn", "It's not your turn!"))
		return
	}

	// Find the troop being used for attack
	var troop *model.Troop
	for i := range attacker.Troops {
		if attacker.Troops[i].Name == req.Troop {
			troop = attacker.Troops[i]
			break
		}
	}
	if troop == nil {
		log.Printf("[WARN][ATTACK] Troop %s not found for user %s", req.Troop, req.Username)
		writeToConnection(conn, dto.Fail(dto.MessageAttackResponse, requestID, "invalid_troop", "Invalid troop used for attack"))
		return
	}

	// Process the attack via game logic
	log.Printf("[INFO][ATTACK] %s attacking with %s targeting %s in room %s", attacker.User.Username, troop.Name, req.Target, req.RoomID)
	damage, isCrit, message := room.Game.PlayTurnSimple(attacker, troop, req.Target)
	isDestroyed := defender.Towers[req.Target].HP <= 0

	success := damage > 0 || isDestroyed

	payload := dto.Response[dto.AttackResult]{Type: dto.MessageAttackResponse, RequestID: requestID, Success: success, Message: message, Data: dto.AttackResult{Attacker: dto.ToPlayer(attacker), Defender: dto.ToPlayer(defender), Troop: troop.Name, Target: req.Target, Damage: int(damage), IsCrit: isCrit, IsDestroyed: isDestroyed, Turn: room.Game.Turn}}

	sendToClient(room.Player1.User.Username, payload)
	sendToClient(room.Player2.User.Username, payload)

	if defender.Towers["king"].HP <= 0 {
		winner, result := room.Game.CheckWinner()
		if result == "" {
			return
		}
		gameOverPayload := dto.Push(dto.MessageGameOverResponse, result, dto.GameOver{Winner: dto.ToPlayer(winner)})
		sendToClient(room.Player1.User.Username, gameOverPayload)
		sendToClient(room.Player2.User.Username, gameOverPayload)
	}
}
