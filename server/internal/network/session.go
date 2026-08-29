package network

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"royaka/internal/database"
	"royaka/internal/model"

	"github.com/jackc/pgx/v5"
)

type Session struct {
	SessionID     string    `json:"session_id"`
	Username      string    `json:"username"`
	Authenticated bool      `json:"authenticated"`
	ExpiresAt     time.Time `json:"expires_at"`
}

var legacySessionCutoff time.Time

func ConfigureLegacySessionCutoff(cutoff time.Time) {
	legacySessionCutoff = cutoff
}

func FindSessionByID(sessionID string) (Session, model.User, error) {
	current := database.Pool().QueryRow(context.Background(), `
		select s.session_id, p.username, true, s.expires_at,
		       p.auth_id, p.legacy_id, p.username, p.created_at, p.last_login,
		       p.is_active, p.exp, p.level, p.games_played, p.games_won, p.avatar, p.gold
		from public.sessions s join public.profiles p on p.auth_id = s.profile_id
		where s.session_id = $1 and s.expires_at > now()
	`, sessionID)
	if session, user, err := scanSession(current); err == nil {
		return session, user, nil
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return Session{}, model.User{}, err
	}
	legacy := database.Pool().QueryRow(context.Background(), `
		select s.session_id, p.username, s.authenticated, s.expires_at,
		       p.auth_id, p.legacy_id, p.username, p.created_at, p.last_login,
		       p.is_active, p.exp, p.level, p.games_played, p.games_won, p.avatar, p.gold
		from public.legacy_sessions s
		join public.profiles p on p.auth_id = s.profile_id
		where s.session_id = $1 and s.authenticated and s.expires_at > now() and $2 > now()
	`, sessionID, legacySessionCutoff)
	return scanSession(legacy)
}

func scanSession(row pgx.Row) (Session, model.User, error) {
	var session Session
	var user model.User
	err := row.Scan(&session.SessionID, &session.Username, &session.Authenticated, &session.ExpiresAt,
		&user.AuthID, &user.ID, &user.Username, &user.CreatedAt, &user.LastLogin,
		&user.IsActive, &user.EXP, &user.Level, &user.GamesPlayed, &user.GamesWon,
		&user.Avatar, &user.Gold)
	if errors.Is(err, pgx.ErrNoRows) {
		return Session{}, model.User{}, pgx.ErrNoRows
	}
	return session, user, err
}

func CreateSession(ctx context.Context, authID string) (Session, error) {
	random := make([]byte, 32)
	if _, err := rand.Read(random); err != nil {
		return Session{}, err
	}
	session := Session{
		SessionID: base64.RawURLEncoding.EncodeToString(random), Authenticated: true,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}
	_, err := database.Pool().Exec(ctx, `
		insert into public.sessions(session_id, profile_id, expires_at) values($1,$2,$3)
	`, session.SessionID, authID, session.ExpiresAt)
	return session, err
}

func DeleteSession(ctx context.Context, sessionID string) error {
	_, err := database.Pool().Exec(ctx, `delete from public.sessions where session_id=$1`, sessionID)
	return err
}

func CleanupExpiredSessions(ctx context.Context) error {
	_, err := database.Pool().Exec(ctx, `
		delete from public.legacy_sessions where expires_at <= now() or $1 <= now()
	`, legacySessionCutoff)
	if err == nil {
		_, err = database.Pool().Exec(ctx, `delete from public.sessions where expires_at <= now()`)
	}
	return err
}
