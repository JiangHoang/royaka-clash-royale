package game

import (
	"log"
	"royaka/internal/network/dto"
)

func (g *Game) HandleTurnTimeout() {
	if !g.Started || g.WinnerDeclared {
		log.Printf("[TURN] ignored timeout because game has ended")
		return
	}

	log.Printf("[TURN] %s skipped due to timeout", g.Turn)

	nextPlayer := g.CurrentPlayer()
	g.SkipTurn(nextPlayer)

	log.Printf("[DEBUG][SKIP_TURN] Turn switched to: %s", g.Turn)

	payload := dto.Push(dto.MessageSkipTurnResponse, "Turn skipped", dto.SkipTurnResult{Turn: g.Turn, Player1: dto.ToPlayer(g.Player1), Player2: dto.ToPlayer(g.Player2)})

	sendToClient(g.Player1.User.Username, payload)
	sendToClient(g.Player2.User.Username, payload)
}
