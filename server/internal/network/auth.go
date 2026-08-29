package network

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"
	"unicode/utf8"

	"royaka/internal/model"
	"royaka/internal/utils"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

type authenticateRequest struct {
	SessionID string `json:"session_id"`
}
type logoutRequest struct {
	SessionID string `json:"session_id"`
}

func writeAuthFailure(conn *websocket.Conn, responseType, message string) {
	_ = conn.WriteJSON(utils.Response{Type: responseType, Success: false, Message: message})
}

func validUsername(username string) bool {
	length := utf8.RuneCountInString(strings.TrimSpace(username))
	return length >= 1 && length <= 32
}

func handleRegister(conn *websocket.Conn, data json.RawMessage) {
	var req utils.RegisterRequest
	if err := json.Unmarshal(data, &req); err != nil || !validUsername(req.Username) || len(req.Password) < 6 {
		writeAuthFailure(conn, "register_response", "Username must be 1-32 characters and password at least 6 characters")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeAuthFailure(conn, "register_response", "Registration failed")
		return
	}
	user := model.NewUser(strings.TrimSpace(req.Username))
	if err := model.AddUser(user, string(hash)); err != nil {
		if errors.Is(err, model.ErrUserExists) {
			writeAuthFailure(conn, "register_response", "Username is already registered")
		} else {
			log.Printf("[ERROR][AUTH] Registration database error: %v", err)
			writeAuthFailure(conn, "register_response", "Registration failed")
		}
		return
	}
	_ = conn.WriteJSON(utils.Response{Type: "register_response", Success: true, Message: "Registered successfully"})
}

func handleLogin(conn *websocket.Conn, data json.RawMessage) {
	var req utils.LoginRequest
	if err := json.Unmarshal(data, &req); err != nil || strings.TrimSpace(req.Username) == "" || req.Password == "" {
		writeAuthFailure(conn, "login_response", "Invalid login data")
		return
	}
	user, passwordHash, err := model.FindCredentialsByUsername(req.Username)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)) != nil {
		writeAuthFailure(conn, "login_response", "Invalid credentials")
		return
	}
	session, err := CreateSession(context.Background(), user.AuthID)
	if err != nil {
		log.Printf("[ERROR][AUTH] Create session failed: %v", err)
		writeAuthFailure(conn, "login_response", "Could not create session")
		return
	}
	user.LastLogin = time.Now()
	if err := model.SaveUser(&user); err != nil {
		log.Printf("[ERROR][AUTH] Update login time failed: %v", err)
	}
	bindIdentity(conn, connectionIdentity{AuthID: user.AuthID, Username: user.Username})
	_ = conn.WriteJSON(utils.Response{Type: "login_response", Success: true, Message: "Login successful", Data: map[string]any{
		"session_id": session.SessionID, "expires_at": session.ExpiresAt.Unix(),
	}})
}

func handleAuthenticate(conn *websocket.Conn, data json.RawMessage) {
	var req authenticateRequest
	if err := json.Unmarshal(data, &req); err != nil || req.SessionID == "" {
		writeAuthFailure(conn, "authenticate_response", "Session ID is required")
		return
	}
	_, user, err := FindSessionByID(req.SessionID)
	if err != nil || !user.IsActive {
		writeAuthFailure(conn, "authenticate_response", "Session is invalid or expired")
		return
	}
	bindIdentity(conn, connectionIdentity{AuthID: user.AuthID, Username: user.Username})
	_ = conn.WriteJSON(utils.Response{Type: "authenticate_response", Success: true, Message: "Authenticated"})
}

func handleLogout(conn *websocket.Conn, data json.RawMessage) {
	var req logoutRequest
	if err := json.Unmarshal(data, &req); err != nil || req.SessionID == "" {
		writeAuthFailure(conn, "logout_response", "Session ID is required")
		return
	}
	if err := DeleteSession(context.Background(), req.SessionID); err != nil {
		writeAuthFailure(conn, "logout_response", "Logout failed")
		return
	}
	removeIdentity(conn)
	_ = conn.WriteJSON(utils.Response{Type: "logout_response", Success: true, Message: "Logged out"})
}

func handleGetUser(conn *websocket.Conn, data json.RawMessage) {
	identity, ok := getIdentity(conn)
	if !ok {
		var req utils.UserRequest
		if json.Unmarshal(data, &req) == nil && req.SessionID != "" {
			_, user, err := FindSessionByID(req.SessionID)
			if err == nil {
				identity, ok = connectionIdentity{AuthID: user.AuthID, Username: user.Username}, true
				bindIdentity(conn, identity)
			}
		}
	}
	if !ok {
		writeAuthFailure(conn, "user_response", "Authentication required")
		return
	}
	user, err := model.FindUserByAuthID(identity.AuthID)
	if err != nil {
		writeAuthFailure(conn, "user_response", "User not found")
		return
	}
	_ = conn.WriteJSON(utils.Response{Type: "user_response", Success: true, Data: map[string]any{"user": user, "maxExp": model.GetMaxExp(user.Level)}})
}
