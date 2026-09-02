package dto

type Empty struct{}

type Session struct {
	SessionID string `json:"session_id"`
	ExpiresAt int64  `json:"expires_at"`
}

type UserData struct {
	User   User `json:"user"`
	MaxEXP int  `json:"maxExp"`
}

type MatchFound struct {
	RoomID   string `json:"room_id"`
	Opponent string `json:"opponent"`
}

type GameData struct {
	User     Player         `json:"user"`
	Opponent Player         `json:"opponent"`
	Player1  string         `json:"player1,omitempty"`
	Map      []BattleEntity `json:"map,omitempty"`
	Time     int64          `json:"time,omitempty"`
	Turn     string         `json:"turn,omitempty"`
}

type AttackResult struct {
	Attacker    Player `json:"attacker"`
	Defender    Player `json:"defender"`
	Troop       string `json:"troop"`
	Target      string `json:"target"`
	Damage      int    `json:"damage"`
	IsCrit      bool   `json:"isCrit"`
	IsDestroyed bool   `json:"isDestroyed"`
	Turn        string `json:"turn"`
}

type HealResult struct {
	Player      Player `json:"player"`
	Opponent    Player `json:"opponent"`
	Troop       string `json:"troop"`
	HealedTower Tower  `json:"healedTower"`
	HealAmount  int    `json:"healAmount"`
	Turn        string `json:"turn"`
}

type SkipTurnResult struct {
	Turn    string `json:"turn"`
	Player1 Player `json:"player1"`
	Player2 Player `json:"player2"`
}

type TroopResult struct {
	Player Player `json:"player"`
}
type ManaUpdate struct {
	Player Player `json:"player"`
}
type GameOver struct {
	Winner Player `json:"winner"`
}

type GameState struct {
	BattleMap     []BattleEntity `json:"battleMap"`
	TimeLeft      int64          `json:"timeLeft"`
	Player1Guard1 float64        `json:"player1Guard1"`
	Player1Guard2 float64        `json:"player1Guard2"`
	Player2Guard1 float64        `json:"player2Guard1"`
	Player2Guard2 float64        `json:"player2Guard2"`
}
