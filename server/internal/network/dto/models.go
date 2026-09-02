package dto

import (
	"royaka/internal/model"
	"time"
)

type User struct {
	ID          string    `json:"id"`
	Username    string    `json:"username"`
	CreatedAt   time.Time `json:"createdAt"`
	LastLogin   time.Time `json:"lastLogin"`
	IsActive    bool      `json:"isActive"`
	EXP         int       `json:"exp"`
	Level       int       `json:"level"`
	GamesPlayed int       `json:"gamesPlayed"`
	GamesWon    int       `json:"gamesWon"`
	Avatar      string    `json:"avatar"`
	Gold        int       `json:"gold"`
}

type Troop struct {
	Name          string  `json:"name"`
	MaxHP         float64 `json:"max_hp"`
	HP            float64 `json:"hp"`
	DMG           float64 `json:"dmg"`
	ATK           float64 `json:"atk"`
	DEF           float64 `json:"def"`
	Mana          int     `json:"mana"`
	Crit          int     `json:"crit"`
	EXP           int     `json:"exp"`
	Speed         float64 `json:"speed"`
	Range         float64 `json:"range"`
	Type          string  `json:"type"`
	Image         string  `json:"image"`
	Description   string  `json:"description"`
	AttackSpeed   float64 `json:"attack_speed"`
	AggroPriority string  `json:"aggro_priority"`
	Rarity        string  `json:"rarity"`
}

type Tower struct {
	Type        string  `json:"type"`
	MaxHP       float64 `json:"max_hp"`
	HP          float64 `json:"hp"`
	ATK         float64 `json:"atk"`
	DEF         float64 `json:"def"`
	Crit        float64 `json:"crit"`
	EXP         int     `json:"exp"`
	Range       float64 `json:"range"`
	AttackSpeed float64 `json:"attack_speed"`
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Area struct {
	TopLeft     Position `json:"top_left"`
	BottomRight Position `json:"bottom_right"`
}

type Player struct {
	Mana           int              `json:"mana"`
	Towers         map[string]Tower `json:"towers"`
	TowerInstances []TowerInstance  `json:"tower_instances,omitempty"`
	Troops         []Troop          `json:"troops"`
	TroopInstances []TroopInstance  `json:"troop_instances,omitempty"`
	Active         bool             `json:"active"`
	User           User             `json:"user"`
	Turn           int              `json:"turn"`
	Gold           int              `json:"gold"`
}

type TroopInstance struct {
	ID         string    `json:"id"`
	Template   Troop     `json:"template"`
	TypeEntity string    `json:"type_entity"`
	Owner      string    `json:"owner"`
	Position   Position  `json:"position"`
	IsDead     bool      `json:"is_dead"`
	LastAttack time.Time `json:"last_attack"`
}

type TowerInstance struct {
	ID          string    `json:"id"`
	Template    Tower     `json:"template"`
	TypeEntity  string    `json:"type_entity"`
	Owner       string    `json:"owner"`
	Area        Area      `json:"area"`
	IsDestroyed bool      `json:"is_destroyed"`
	LastAttack  time.Time `json:"last_attack"`
}

type BattleEntity interface{ battleEntityDTO() }

func (TroopInstance) battleEntityDTO() {}
func (TowerInstance) battleEntityDTO() {}

func ToUser(value *model.User) User {
	if value == nil {
		return User{}
	}
	return User{ID: value.ID, Username: value.Username, CreatedAt: value.CreatedAt, LastLogin: value.LastLogin, IsActive: value.IsActive, EXP: value.EXP, Level: value.Level, GamesPlayed: value.GamesPlayed, GamesWon: value.GamesWon, Avatar: value.Avatar, Gold: value.Gold}
}

func ToTroop(value *model.Troop) Troop {
	if value == nil {
		return Troop{}
	}
	return Troop{Name: value.Name, MaxHP: value.MaxHP, HP: value.HP, DMG: value.DMG, ATK: value.ATK, DEF: value.DEF, Mana: value.MANA, Crit: value.CRIT, EXP: value.EXP, Speed: value.Speed, Range: value.Range, Type: value.Type, Image: value.Image, Description: value.Description, AttackSpeed: value.AttackSpeed, AggroPriority: value.AggroPriority, Rarity: value.Rarity}
}

func ToTroops(values []*model.Troop) []Troop {
	result := make([]Troop, 0, len(values))
	for _, value := range values {
		result = append(result, ToTroop(value))
	}
	return result
}

func ToTroopValues(values []model.Troop) []Troop {
	result := make([]Troop, 0, len(values))
	for index := range values {
		result = append(result, ToTroop(&values[index]))
	}
	return result
}

func ToTower(value *model.Tower) Tower {
	if value == nil {
		return Tower{}
	}
	return Tower{Type: value.Type, MaxHP: value.MaxHP, HP: value.HP, ATK: value.ATK, DEF: value.DEF, Crit: value.CRIT, EXP: value.EXP, Range: value.Range, AttackSpeed: value.AttackSpeed}
}

func toPosition(value model.Position) Position { return Position{X: value.X, Y: value.Y} }

func ToTroopInstance(value *model.TroopInstance) TroopInstance {
	if value == nil {
		return TroopInstance{}
	}
	return TroopInstance{ID: value.ID, Template: ToTroop(value.Template), TypeEntity: value.TypeEntity, Owner: value.Owner, Position: toPosition(value.Position), IsDead: value.IsDead, LastAttack: value.LastAttackTime}
}

func ToTowerInstance(value *model.TowerInstance) TowerInstance {
	if value == nil {
		return TowerInstance{}
	}
	return TowerInstance{ID: value.ID, Template: ToTower(value.Template), TypeEntity: value.TypeEntity, Owner: value.Owner, Area: Area{TopLeft: toPosition(value.Area.TopLeft), BottomRight: toPosition(value.Area.BottomRight)}, IsDestroyed: value.IsDestroyed, LastAttack: value.LastAttackTime}
}

func ToPlayer(value *model.Player) Player {
	if value == nil {
		return Player{}
	}
	towers := make(map[string]Tower, len(value.Towers))
	for key, tower := range value.Towers {
		towers[key] = ToTower(tower)
	}
	towerInstances := make([]TowerInstance, 0, len(value.TowerInstances))
	for _, item := range value.TowerInstances {
		towerInstances = append(towerInstances, ToTowerInstance(item))
	}
	troopInstances := make([]TroopInstance, 0, len(value.TroopInstances))
	for _, item := range value.TroopInstances {
		troopInstances = append(troopInstances, ToTroopInstance(item))
	}
	return Player{Mana: value.Mana, Towers: towers, TowerInstances: towerInstances, Troops: ToTroops(value.Troops), TroopInstances: troopInstances, Active: value.Active, User: ToUser(value.User), Turn: value.Turn, Gold: value.Gold}
}
