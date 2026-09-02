package game

import (
	"log"
	"royaka/internal/model"
	"royaka/internal/network/dto"
	"royaka/internal/network/wsconn"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	clients   = make(map[string]*ClientConnection)
	clientsMu sync.RWMutex

	pendingPlayers = make(map[string]bool)
	pendingMu      sync.RWMutex

	rooms   = make(map[string]*Room)
	roomsMu sync.RWMutex

	matchQueues = map[string]chan *model.Player{
		"simple":   make(chan *model.Player, 100),
		"enhanced": make(chan *model.Player, 100),
	}
	matchmakerOnce sync.Once

	invalidRequestMessage = "Invalid request"
	roomRequestMessage    = "Room not found"
	manaRequestMessage    = "Not enough mana!"
)

func writeToConnection(conn *websocket.Conn, payload any) {
	if err := wsconn.Send(conn, payload); err != nil {
		log.Printf("[ERROR][SEND] Failed to send response: %v", err)
	}
}

func sendToClient(username string, payload any) {
	clientsMu.RLock()
	client, exists := clients[username]
	clientsMu.RUnlock()

	if !exists || client == nil || client.Conn == nil {
		log.Printf("[WARN][SEND] Client %s not found or connection is nil", username)
		return
	}

	if err := client.SafeWrite(payload); err != nil {
		log.Printf("[ERROR][SEND] Failed to send to %s: %v", username, err)
	}
}

func toBattleEntities(entities []BattleEntity) []dto.BattleEntity {
	result := make([]dto.BattleEntity, 0, len(entities))
	for _, entity := range entities {
		switch value := entity.(type) {
		case *model.TroopInstance:
			result = append(result, dto.ToTroopInstance(value))
		case *model.TowerInstance:
			result = append(result, dto.ToTowerInstance(value))
		}
	}
	return result
}
