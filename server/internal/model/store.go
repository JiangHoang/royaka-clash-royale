package model

import (
	"context"
	"errors"
	"strings"

	"royaka/internal/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)

const userColumns = `auth_id, legacy_id, username, created_at, last_login,
 is_active, exp, level, games_played, games_won, avatar, gold`

func scanUser(row pgx.Row) (User, error) {
	var user User
	err := row.Scan(
		&user.AuthID, &user.ID, &user.Username, &user.CreatedAt, &user.LastLogin,
		&user.IsActive, &user.EXP, &user.Level, &user.GamesPlayed, &user.GamesWon,
		&user.Avatar, &user.Gold,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrUserNotFound
	}
	return user, err
}

func FindUserByUsername(username string) (User, error) {
	return scanUser(database.Pool().QueryRow(context.Background(), `
		select `+userColumns+` from public.profiles where lower(username) = lower($1)
	`, strings.TrimSpace(username)))
}

func FindCredentialsByUsername(username string) (User, string, error) {
	var user User
	var passwordHash string
	err := database.Pool().QueryRow(context.Background(), `
		select `+userColumns+`, password_hash from public.profiles
		where lower(btrim(username)) = lower(btrim($1))
	`, username).Scan(&user.AuthID, &user.ID, &user.Username, &user.CreatedAt, &user.LastLogin,
		&user.IsActive, &user.EXP, &user.Level, &user.GamesPlayed, &user.GamesWon,
		&user.Avatar, &user.Gold, &passwordHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, "", ErrUserNotFound
	}
	return user, passwordHash, err
}

func AddUser(user *User, passwordHash string) error {
	err := database.Pool().QueryRow(context.Background(), `
		insert into public.profiles
		(legacy_id, username, password_hash, created_at, last_login, is_active,
		 exp, level, games_played, games_won, avatar, gold)
		values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		returning auth_id
	`, user.ID, strings.TrimSpace(user.Username), passwordHash, user.CreatedAt, user.LastLogin,
		user.IsActive, user.EXP, user.Level, user.GamesPlayed, user.GamesWon, user.Avatar, user.Gold).Scan(&user.AuthID)
	var postgresError *pgconn.PgError
	if errors.As(err, &postgresError) && postgresError.Code == "23505" {
		return ErrUserExists
	}
	return err
}

func FindUserByAuthID(authID string) (User, error) {
	return scanUser(database.Pool().QueryRow(context.Background(), `
		select `+userColumns+` from public.profiles where auth_id = $1
	`, authID))
}

// SaveUser persists mutable profile fields without rewriting the password hash.
func SaveUser(user *User) error {
	command, err := database.Pool().Exec(context.Background(), `
		update public.profiles set
			last_login = $2, is_active = $3, exp = $4, level = $5,
			games_played = $6, games_won = $7, avatar = $8, gold = $9
		where auth_id = $1
	`, user.AuthID, user.LastLogin, user.IsActive, user.EXP, user.Level,
		user.GamesPlayed, user.GamesWon, user.Avatar, user.Gold)
	if err != nil {
		return err
	}
	if command.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}

// ApplyGameResult locks both profiles and applies deltas in one transaction so
// simultaneous game completions cannot lose wins, EXP, or gold.
func ApplyGameResult(winner, loser *User, isDraw bool, winnerGold, loserGold int) error {
	ctx := context.Background()
	tx, err := database.Pool().Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, `
		select auth_id from public.profiles where auth_id::text = any($1::text[]) order by auth_id for update
	`, []string{winner.AuthID, loser.AuthID})
	if err != nil {
		return err
	}
	locked := 0
	for rows.Next() {
		var ignored string
		if err := rows.Scan(&ignored); err != nil {
			rows.Close()
			return err
		}
		locked++
	}
	rows.Close()
	if locked != 2 {
		return ErrUserNotFound
	}

	currentWinner, err := scanUser(tx.QueryRow(ctx, `select `+userColumns+` from public.profiles where auth_id=$1`, winner.AuthID))
	if err != nil {
		return err
	}
	currentLoser, err := scanUser(tx.QueryRow(ctx, `select `+userColumns+` from public.profiles where auth_id=$1`, loser.AuthID))
	if err != nil {
		return err
	}
	if isDraw {
		currentWinner.AddExp(10)
		currentLoser.AddExp(10)
	} else {
		currentWinner.GamesWon++
		currentWinner.AddExp(30)
	}
	currentWinner.GamesPlayed++
	currentLoser.GamesPlayed++
	currentWinner.Gold += winnerGold
	currentLoser.Gold += loserGold

	for _, user := range []*User{&currentWinner, &currentLoser} {
		if _, err := tx.Exec(ctx, `update public.profiles set exp=$2, level=$3,
			games_played=$4, games_won=$5, gold=$6 where auth_id=$1`,
			user.AuthID, user.EXP, user.Level, user.GamesPlayed, user.GamesWon, user.Gold); err != nil {
			return err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	*winner, *loser = currentWinner, currentLoser
	return nil
}
