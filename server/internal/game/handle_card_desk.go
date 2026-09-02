package game

import (
	"encoding/json"
	"log"
	"royaka/internal/model"
	"royaka/internal/network/dto"

	"github.com/gorilla/websocket"
)

func HandleGetDesk(conn *websocket.Conn, requestID string, data json.RawMessage) {
	troops, err := model.LoadTroop()
	if err != nil {
		log.Println("loadTroop error:", err)
		writeToConnection(conn, dto.Fail(dto.MessageDeckResponse, requestID, "deck_unavailable", "Failed to load troops"))
		return
	}

	writeToConnection(conn, dto.OK(dto.MessageDeckResponse, requestID, "Troop data loaded", dto.ToTroopValues(troops)))
}
