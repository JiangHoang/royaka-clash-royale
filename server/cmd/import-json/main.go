package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"royaka/internal/database"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type legacyUser struct {
	ID          string    `json:"id"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
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
type legacySession struct {
	SessionID     string `json:"session_id"`
	Username      string `json:"username"`
	Authenticated bool   `json:"authenticated"`
}
type troop struct {
	Name          string  `json:"name"`
	MaxHP         float64 `json:"max_hp"`
	HP            float64 `json:"hp"`
	DMG           float64 `json:"dmg"`
	ATK           float64 `json:"atk"`
	DEF           float64 `json:"def"`
	MANA          int     `json:"mana"`
	CRIT          int     `json:"crit"`
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
type tower struct {
	Type        string  `json:"type"`
	MaxHP       float64 `json:"max_hp"`
	HP          float64 `json:"hp"`
	ATK         float64 `json:"atk"`
	DEF         float64 `json:"def"`
	CRIT        float64 `json:"crit"`
	EXP         int     `json:"exp"`
	Range       float64 `json:"range"`
	AttackSpeed float64 `json:"attack_speed"`
}

func readJSON(path string, target any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func main() {
	dataDir := flag.String("data-dir", "assets/data", "legacy JSON directory")
	flag.Parse()
	databaseURL := os.Getenv("DATABASE_URL")
	cutoff, err := time.Parse(time.RFC3339, os.Getenv("LEGACY_SESSION_CUTOFF"))
	if databaseURL == "" || err != nil {
		log.Fatal("DATABASE_URL and an RFC3339 LEGACY_SESSION_CUTOFF are required")
	}
	if cutoff.Before(time.Now()) || cutoff.After(time.Now().Add(7*24*time.Hour+5*time.Minute)) {
		log.Fatal("LEGACY_SESSION_CUTOFF must be within the next seven days")
	}

	var users []legacyUser
	var sessions []legacySession
	var troops []troop
	var towers []tower
	for name, target := range map[string]any{"users.json": &users, "sessions.json": &sessions, "troops.json": &troops, "towers.json": &towers} {
		if err := readJSON(filepath.Join(*dataDir, name), target); err != nil {
			log.Fatalf("Read %s: %v", name, err)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	if err := database.Connect(ctx, databaseURL); err != nil {
		log.Fatalf("Connect database: %v", err)
	}
	defer database.Close()

	userIDs := make(map[string]string, len(users))
	for _, user := range users {
		var authID string
		err := database.Pool().QueryRow(ctx, `select auth_id from public.profiles where lower(btrim(username))=lower(btrim($1))`, user.Username).Scan(&authID)
		if err == pgx.ErrNoRows {
			authID = uuid.NewString()
		} else if err != nil {
			log.Fatalf("Lookup user %q: %v", user.Username, err)
		}
		_, err = database.Pool().Exec(ctx, `insert into public.profiles
			(auth_id,legacy_id,username,password_hash,created_at,last_login,is_active,exp,level,games_played,games_won,avatar,gold)
			values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
			on conflict(auth_id) do update set legacy_id=excluded.legacy_id,username=excluded.username,password_hash=excluded.password_hash,
			created_at=excluded.created_at,last_login=excluded.last_login,is_active=excluded.is_active,exp=excluded.exp,
			level=excluded.level,games_played=excluded.games_played,games_won=excluded.games_won,avatar=excluded.avatar,gold=excluded.gold`,
			authID, user.ID, user.Username, user.Password, user.CreatedAt, user.LastLogin, user.IsActive, user.EXP, user.Level, user.GamesPlayed, user.GamesWon, user.Avatar, user.Gold)
		if err != nil {
			log.Fatalf("Upsert user %q: %v", user.Username, err)
		}
		userIDs[strings.ToLower(strings.TrimSpace(user.Username))] = authID
	}
	if _, err := database.Pool().Exec(ctx, `alter table public.profiles alter column password_hash set not null`); err != nil {
		log.Fatalf("Enforce profile password hashes: %v", err)
	}
	for _, item := range troops {
		_, err := database.Pool().Exec(ctx, `insert into public.troops(name,max_hp,hp,dmg,atk,def,mana,crit,exp,speed,range,type,image,description,attack_speed,aggro_priority,rarity)
		values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17) on conflict(name) do update set max_hp=excluded.max_hp,hp=excluded.hp,dmg=excluded.dmg,atk=excluded.atk,def=excluded.def,mana=excluded.mana,crit=excluded.crit,exp=excluded.exp,speed=excluded.speed,range=excluded.range,type=excluded.type,image=excluded.image,description=excluded.description,attack_speed=excluded.attack_speed,aggro_priority=excluded.aggro_priority,rarity=excluded.rarity`, item.Name, item.MaxHP, item.HP, item.DMG, item.ATK, item.DEF, item.MANA, item.CRIT, item.EXP, item.Speed, item.Range, item.Type, item.Image, item.Description, item.AttackSpeed, item.AggroPriority, item.Rarity)
		if err != nil {
			log.Fatalf("Upsert troop %q: %v", item.Name, err)
		}
	}
	for _, item := range towers {
		_, err := database.Pool().Exec(ctx, `insert into public.towers(type,max_hp,hp,atk,def,crit,exp,range,attack_speed) values($1,$2,$3,$4,$5,$6,$7,$8,$9) on conflict(type) do update set max_hp=excluded.max_hp,hp=excluded.hp,atk=excluded.atk,def=excluded.def,crit=excluded.crit,exp=excluded.exp,range=excluded.range,attack_speed=excluded.attack_speed`, item.Type, item.MaxHP, item.HP, item.ATK, item.DEF, item.CRIT, item.EXP, item.Range, item.AttackSpeed)
		if err != nil {
			log.Fatalf("Upsert tower %q: %v", item.Type, err)
		}
	}
	for _, item := range sessions {
		profileID := userIDs[strings.ToLower(strings.TrimSpace(item.Username))]
		_, err := database.Pool().Exec(ctx, `insert into public.legacy_sessions(session_id,profile_id,authenticated,expires_at) values($1,$2,$3,$4) on conflict(session_id) do update set profile_id=excluded.profile_id,authenticated=excluded.authenticated,expires_at=excluded.expires_at`, item.SessionID, profileID, item.Authenticated, cutoff)
		if err != nil {
			log.Fatalf("Upsert session %q: %v", item.SessionID, err)
		}
	}
	log.Printf("Import complete: users=%d sessions=%d troops=%d towers=%d", len(users), len(sessions), len(troops), len(towers))
}
