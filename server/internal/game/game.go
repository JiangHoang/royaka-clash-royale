package game

import (
	"fmt"
	"log"
	"math/rand"
	"royaka/internal/model"
	"royaka/internal/utils"
	"sync"
	"time"
)

type Game struct {
	Player1         *model.Player
	Player2         *model.Player
	Turn            string
	Started         bool
	Enhanced        bool
	StartTime       time.Time
	MaxTime         time.Duration
	LastTick        time.Time
	TickerStopChan  chan struct{}
	BattleSystem    *BattleSystem
	WinnerDeclared  bool
	TurnTimerCancel func()
}

// ===================== Game Initialization =====================

func NewGame(p1, p2 *model.Player, mode string) *Game {
	if mode != "simple" && mode != "enhanced" {
		log.Fatal("Invalid game mode")
	}

	startingPlayer := p1.User.Username
	if rand.Intn(2) == 0 {
		startingPlayer = p2.User.Username
	}

	p1.LastManaRegen = time.Now()
	p2.LastManaRegen = time.Now()

	p1.TowerInstances = model.CreateTowerInstances(p1.Towers, p1.User.Username, true)
	p2.TowerInstances = model.CreateTowerInstances(p2.Towers, p2.User.Username, false)

	initialEntities := make(map[string][]BattleEntity)
	for _, ti := range p1.TowerInstances {
		posKey := ti.GetPosition().String()
		initialEntities[posKey] = append(initialEntities[posKey], ti)
	}
	for _, ti := range p2.TowerInstances {
		posKey := ti.GetPosition().String()
		initialEntities[posKey] = append(initialEntities[posKey], ti)
	}

	battleSystem := &BattleSystem{
		BattleMap:      initialEntities,
		MapMutex:       sync.RWMutex{},
		TickerStopChan: make(chan struct{}),
		TickRate:       100 * time.Millisecond,
	}

	game := &Game{
		Player1:        p1,
		Player2:        p2,
		Turn:           startingPlayer,
		Started:        true,
		Enhanced:       (mode == "enhanced"),
		BattleSystem:   battleSystem,
		TickerStopChan: battleSystem.TickerStopChan,
		WinnerDeclared: false,
	}

	if game.Enhanced {
		game.StartTime = time.Now()
		game.MaxTime = 3 * time.Minute
		time.AfterFunc(3*time.Second, func() {
			game.StartTime = time.Now()
			go game.startTicker()
		})
	} else if !game.Enhanced {
		game.StartTurnTimer()
	}

	return game
}

// ===================== Turn Management =====================

func (g *Game) CurrentPlayer() *model.Player {
	if g.Enhanced {
		return nil
	}
	if g.Turn == g.Player1.User.Username {
		return g.Player1
	}
	return g.Player2
}

func (g *Game) SwitchTurn() {
	if g.Turn == g.Player1.User.Username {
		g.Turn = g.Player2.User.Username
	} else {
		g.Turn = g.Player1.User.Username
	}

	g.LastTick = time.Now()

	// khởi động timer cho lượt mới
	g.StartTurnTimer()

	nextPlayer := g.CurrentPlayer()
	if nextPlayer.Turn > 0 {
		nextPlayer.Mana += 3
		if nextPlayer.Mana > 10 {
			nextPlayer.Mana = 10
		}
	}
}

func (g *Game) SkipTurn(player *model.Player) {
	player.Turn++
	g.SwitchTurn()
}

func (g *Game) StartTurnTimer() {
	// Hủy timer cũ nếu còn
	if g.TurnTimerCancel != nil {
		g.TurnTimerCancel()
	}

	timer := time.NewTimer(30 * time.Second)
	cancelChan := make(chan struct{})

	g.TurnTimerCancel = func() {
		timer.Stop()
		close(cancelChan)
	}

	go func(turn string) {

		time.Sleep(1 * time.Second)

		select {
		case <-timer.C:
			// Chỉ xử lý timeout nếu vẫn là lượt đó
			if g.Turn == turn {
				log.Printf("[TURN] player %s timed out", turn)
				g.HandleTurnTimeout()
			}
		case <-cancelChan:
			// Lượt kết thúc hợp lệ
		}
	}(g.Turn)
}

// ===================== Game Tick & Loop =====================

func (g *Game) startTicker() {
	manaTicker := time.NewTicker(200 * time.Millisecond)
	tickTicker := time.NewTicker(g.BattleSystem.TickRate)
	cleanupTicker := time.NewTicker(5 * time.Second)

	defer func() {
		tickTicker.Stop()
		cleanupTicker.Stop()
	}()

	for {
		select {
		case <-tickTicker.C:
			g.UpdateBattleMap()
			g.BroadcastGameState()
		case <-manaTicker.C:
			g.UpdateMana()
		case <-cleanupTicker.C:
			g.BattleSystem.CleanupDeadEntities()
		case <-g.BattleSystem.TickerStopChan:
			return
		}
	}
}

func (g *Game) StopGameLoop() {
	if g.Started {
		g.Started = false
		close(g.TickerStopChan)
	}
}

func (g *Game) UpdateMana() {
	now := time.Now()

	for _, player := range []*model.Player{g.Player1, g.Player2} {
		if player.Mana < 10 && now.Sub(player.LastManaRegen) >= 2*time.Second {
			player.Mana++
			player.LastManaRegen = now

			sendToClient(player.User.Username, utils.Response{
				Type:    "mana_update",
				Success: true,
				Message: fmt.Sprintf("Mana: %d", player.Mana),
				Data: map[string]interface{}{
					"player": player,
				},
			})
		}
	}
}

// ===================== Game Outcome =====================

func (g *Game) CheckWinner() (*model.Player, string) {
	if g.WinnerDeclared {
		return nil, ""
	}

	// Kiểm tra King Tower bị phá
	if g.Player1.Towers["king"].HP <= 0.0 {
		g.WinnerDeclared = true
		g.StopGameLoop()

		p1Gold, p2Gold := 0, 0
		if g.Enhanced {
			p1Gold, p2Gold = g.Player1.Gold, g.Player2.Gold
		}
		AwardEXP(g.Player2.User, g.Player1.User, false, p2Gold, p1Gold)
		fmt.Printf("Winner: %s\n", g.Player2.User.Username)
		return g.Player2, g.Player2.User.Username + " wins!"
	}

	if g.Player2.Towers["king"].HP <= 0.0 {
		g.WinnerDeclared = true
		g.StopGameLoop()

		p1Gold, p2Gold := 0, 0
		if g.Enhanced {
			p1Gold, p2Gold = g.Player1.Gold, g.Player2.Gold
		}
		AwardEXP(g.Player1.User, g.Player2.User, false, p1Gold, p2Gold)
		fmt.Printf("Winner: %s\n", g.Player1.User.Username)
		return g.Player1, g.Player1.User.Username + " wins!"
	}

	// Hết giờ trong enhanced mode => xử lý tính điểm
	if g.Enhanced && time.Since(g.StartTime) > g.MaxTime {
		p1Score := g.Player1.DestroyedCount()
		p2Score := g.Player2.DestroyedCount()

		g.WinnerDeclared = true
		g.StopGameLoop()

		// Cộng gold
		if p1Score < p2Score {
			AwardEXP(g.Player1.User, g.Player2.User, false, g.Player1.Gold, g.Player2.Gold)
			fmt.Printf("Winner by score: %s\n", g.Player1.User.Username)
			return g.Player1, g.Player1.User.Username + " wins by score!"
		}

		if p2Score < p1Score {
			AwardEXP(g.Player2.User, g.Player1.User, false, g.Player2.Gold, g.Player1.Gold)
			fmt.Printf("Winner by score: %s\n", g.Player2.User.Username)
			return g.Player2, g.Player2.User.Username + " wins by score!"
		}

		// Hòa điểm
		AwardEXP(g.Player1.User, g.Player2.User, true, g.Player1.Gold, g.Player2.Gold)
		fmt.Println("Game ended in a draw by score")
		return nil, "It's a draw!"
	}

	// Nếu chưa có ai thắng
	return nil, ""
}

func (g *Game) SetWinner(winner *model.Player) {
	if winner == g.Player1 {
		g.WinnerDeclared = true
		g.StopGameLoop()
		AwardEXP(g.Player1.User, g.Player2.User, false, 0, 0)
	} else if winner == g.Player2 {
		g.WinnerDeclared = true
		g.StopGameLoop()
		AwardEXP(g.Player2.User, g.Player1.User, false, 0, 0)
	}
}

func AwardEXP(winner, loser *model.User, isDraw bool, winnerGold, loserGold int) {
	if err := model.ApplyGameResult(winner, loser, isDraw, winnerGold, loserGold); err != nil {
		log.Printf("[ERROR][GAME] Failed to atomically save game result: %v", err)
	}
}

// ===================== Game State Broadcasting =====================

func (g *Game) BroadcastGameState() {
	elapsed := time.Since(g.StartTime)
	timeLeft := g.MaxTime - elapsed
	if timeLeft < 0 {
		timeLeft = 0
	}

	for _, player := range []*model.Player{g.Player1, g.Player2} {
		sendToClient(player.User.Username, utils.Response{
			Type:    "game_state",
			Success: true,
			Message: "Game updated",
			Data: map[string]interface{}{
				"battleMap":     g.BattleSystem.GetEntityList(),
				"timeLeft":      timeLeft.Milliseconds(),
				"player1Guard1": g.Player1.Towers["guard1"].HP,
				"player1Guard2": g.Player1.Towers["guard2"].HP,
				"player2Guard1": g.Player2.Towers["guard1"].HP,
				"player2Guard2": g.Player2.Towers["guard2"].HP,
			},
		})
	}

	if timeLeft == 0 && !g.WinnerDeclared {
		g.checkWinCondition()
	}
}

// ===================== Utility =====================

func (g *Game) Opponent(p *model.Player) *model.Player {
	if g.Player1.User.Username == p.User.Username {
		return g.Player2
	}
	return g.Player1
}

func (g *Game) getPlayerID(isPlayer1 bool) string {
	if isPlayer1 {
		return g.Player1.User.Username
	}
	return g.Player2.User.Username
}

func (g *Game) isInSafeZone(currentY, safeZoneY float64, isPlayer1 bool) bool {
	return (isPlayer1 && currentY < safeZoneY) || (!isPlayer1 && currentY > safeZoneY)
}

func (g *Game) isInRiver(currentY float64) bool {
	return currentY >= RIVER_TOP && currentY <= RIVER_BOTTOM
}
