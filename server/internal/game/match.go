package game

import (
	"encoding/json"
	"fmt"
	"log"
	"royaka/internal/model"
	"royaka/internal/network/dto"
	"royaka/internal/utils"
	"time"

	"github.com/gorilla/websocket"
)

func HandleFindMatch(conn *websocket.Conn, requestID string, data json.RawMessage) {
	var req dto.FindMatchRequest

	// Parse & validate request data
	if err := json.Unmarshal(data, &req); err != nil || req.Username == "" || req.Mode == "" {
		log.Printf("[WARN][MATCH] invalid request: %v", err)
		writeToConnection(conn, dto.Fail(dto.MessageFindMatchResponse, requestID, "invalid_request", invalidRequestMessage))
		return
	}
	if req.Mode != "simple" && req.Mode != "enhanced" {
		writeToConnection(conn, dto.Fail(dto.MessageFindMatchResponse, requestID, "invalid_mode", "Invalid game mode"))
		return
	}

	username := req.Username
	log.Printf("[INFO][MATCH] matchmaking request: user=%s, mode=%s", username, req.Mode)

	// Check if user already in matchmaking queue
	pendingMu.Lock()
	if pendingPlayers[username] {
		pendingMu.Unlock()
		log.Printf("[WARN][MATCH] user %s already in queue", username)
		writeToConnection(conn, dto.Fail(dto.MessageFindMatchResponse, requestID, "already_in_queue", "Already in queue"))
		return
	}
	pendingPlayers[username] = true
	pendingMu.Unlock()

	// Save client connection for future communication
	clientConn := &ClientConnection{Conn: conn}
	clientsMu.Lock()
	clients[username] = clientConn
	clientsMu.Unlock()

	// Create Player instance and register
	user, err := model.FindUserByUsername(req.Username)
	if err != nil {
		clientConn.SafeWrite(dto.Fail(dto.MessageFindMatchResponse, requestID, "profile_unavailable", "Failed to load player profile"))
		CleanupUser(username)
		return
	}
	player, err := model.NewPlayer(&user, req.Mode)
	if err != nil {
		log.Printf("[ERROR][MATCH] Failed to load game definitions: %v", err)
		clientConn.SafeWrite(dto.Fail(dto.MessageFindMatchResponse, requestID, "game_data_unavailable", "Game data is unavailable"))
		CleanupUser(username)
		return
	}
	model.RegisterConnection(conn, player)

	// Confirm queue entry
	clientConn.SafeWrite(dto.OK(dto.MessageFindMatchResponse, requestID, "Added to match queue. Waiting for opponent...", dto.Empty{}))

	// Push player to matchmaking queue
	go queuePlayer(player, clientConn, req.Mode, username)

	// Start the matchmaker loop once
	matchmakerOnce.Do(func() {
		startMatchmaker()
		log.Println("[MATCH] matchmaker started")
	})
}

func queuePlayer(player *model.Player, clientConn *ClientConnection, mode, username string) {
	queue, ok := matchQueues[mode]
	fmt.Printf("[INFO][MATCH] player queued for mode %s (queue: %v)\n", mode, matchQueues[mode])

	if !ok {
		log.Printf("[WARN][MATCH] invalid mode %s for user %s", mode, username)
		clientConn.SafeWrite(dto.Fail(dto.MessageFindMatchResponse, "", "invalid_mode", "Invalid game mode"))
		return
	}

	select {
	case queue <- player:
	case <-time.After(30 * time.Second):
		log.Printf("[WARN][MATCH] enqueue timeout for user %s", username)
		CleanupUser(username)
		return
	}

	timer := time.NewTimer(30 * time.Second)
	defer timer.Stop()

	select {
	case <-player.Matched:
		log.Printf("[INFO][MATCH] user %s matched", username)
	case <-timer.C:
		log.Printf("[WARN][MATCH] matchmaking timeout for user %s", username)
		RemovePlayerFromQueue(player)
		CleanupUser(username)
		clientConn.SafeWrite(dto.Event[dto.Empty]{Type: dto.MessageMatchTimeout, Success: false, Message: "Matchmaking timed out. No opponents found.", Error: &dto.Error{Code: "match_timeout", Message: "Matchmaking timed out. No opponents found."}})
	}
}

func RemovePlayerFromQueue(targetPlayer *model.Player) {
	for mode, queue := range matchQueues {
		var tempPlayers []*model.Player

		// Drain queue and filter out target player
		queueLength := len(queue)
		for i := 0; i < queueLength; i++ {
			select {
			case player := <-queue:
				if player.User.Username != targetPlayer.User.Username {
					tempPlayers = append(tempPlayers, player)
				} else {
					log.Printf("[INFO][QUEUE] removed player %s from %s queue", player.User.Username, mode)
				}
			default:
				break
			}
		}

		// Put back remaining players
		for _, player := range tempPlayers {
			select {
			case queue <- player:
			default:
				log.Printf("[WARN][QUEUE] failed to restore player %s to queue", player.User.Username)
			}
		}
	}
}

func CleanupUser(username string) {
	pendingMu.Lock()
	delete(pendingPlayers, username)
	pendingMu.Unlock()

	clientsMu.Lock()
	delete(clients, username)
	clientsMu.Unlock()

	log.Printf("[INFO][CLEANUP] removed user %s from queues", username)
}

func startMatchmaker() {
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			for mode, queue := range matchQueues {
				if len(queue) < 2 {
					continue
				}
				player1, player2 := <-queue, <-queue
				log.Printf("[INFO][MATCH] pairing players in mode %s: %s vs %s", mode, player1.User.Username, player2.User.Username)
				if validatePlayers(player1, player2, mode) {
					handleMatch(player1, player2, mode)
				}
			}
		}
	}()
}

func validatePlayers(p1, p2 *model.Player, mode string) bool {
	clientsMu.RLock()
	_, ok1 := clients[p1.User.Username]
	_, ok2 := clients[p2.User.Username]
	clientsMu.RUnlock()

	if !ok1 || !ok2 || p1.User.Username == p2.User.Username {
		if ok1 {
			matchQueues[mode] <- p1
		}
		if ok2 && p1.User.Username != p2.User.Username {
			matchQueues[mode] <- p2
		}
		log.Printf("[WARN][MATCH] validation failed for %s vs %s", p1.User.Username, p2.User.Username)
		return false
	}

	return true
}

func handleMatch(p1, p2 *model.Player, mode string) {
	clientsMu.RLock()
	conn1 := clients[p1.User.Username]
	conn2 := clients[p2.User.Username]
	clientsMu.RUnlock()

	// Signal matched
	p1.Matched <- true
	p2.Matched <- true

	// Create and register room
	roomID := utils.GenerateRoomID()
	room := NewRoom(roomID, p1, p2, mode)

	roomsMu.Lock()
	rooms[roomID] = room
	roomsMu.Unlock()
	RegisterRoom(roomID, room)

	log.Printf("[INFO][ROOM] created room %s with players %s and %s", roomID, p1.User.Username, p2.User.Username)

	// Notify players
	notifyMatchFound(conn1, p2.User.Username, roomID)
	notifyMatchFound(conn2, p1.User.Username, roomID)

	// Remove from pending
	pendingMu.Lock()
	delete(pendingPlayers, p1.User.Username)
	delete(pendingPlayers, p2.User.Username)
	pendingMu.Unlock()
}

func notifyMatchFound(conn *ClientConnection, opponent, roomID string) {
	conn.SafeWrite(dto.Push(dto.MessageMatchFound, "Match found!", dto.MatchFound{RoomID: roomID, Opponent: opponent}))
}
