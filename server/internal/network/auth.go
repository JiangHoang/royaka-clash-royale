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
	"royaka/internal/network/dto"
	"royaka/internal/network/wsconn"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

func writeAuthFailure(conn *websocket.Conn, responseType dto.MessageType, requestID, code, message string) {
	_ = wsconn.Send(conn, dto.Fail(responseType, requestID, code, message))
}

func validUsername(username string) bool {
	length := utf8.RuneCountInString(strings.TrimSpace(username))
	return length >= 1 && length <= 32
}

func handleRegister(conn *websocket.Conn, requestID string, data json.RawMessage) {
	var req dto.RegisterRequest
	if err := json.Unmarshal(data, &req); err != nil || !validUsername(req.Username) || len(req.Password) < 6 {
		writeAuthFailure(conn, dto.MessageRegisterResponse, requestID, "invalid_registration", "Username must be 1-32 characters and password at least 6 characters")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeAuthFailure(conn, dto.MessageRegisterResponse, requestID, "registration_failed", "Registration failed")
		return
	}
	user := model.NewUser(strings.TrimSpace(req.Username))
	if err := model.AddUser(user, string(hash)); err != nil {
		if errors.Is(err, model.ErrUserExists) {
			writeAuthFailure(conn, dto.MessageRegisterResponse, requestID, "username_exists", "Username is already registered")
		} else {
			log.Printf("[ERROR][AUTH] Registration database error: %v", err)
			writeAuthFailure(conn, dto.MessageRegisterResponse, requestID, "registration_failed", "Registration failed")
		}
		return
	}
	_ = wsconn.Send(conn, dto.OK(dto.MessageRegisterResponse, requestID, "Registered successfully", dto.Empty{}))
}

func handleLogin(conn *websocket.Conn, requestID string, data json.RawMessage) {
	var req dto.LoginRequest
	if err := json.Unmarshal(data, &req); err != nil || strings.TrimSpace(req.Username) == "" || req.Password == "" {
		writeAuthFailure(conn, dto.MessageLoginResponse, requestID, "invalid_login", "Invalid login data")
		return
	}
	user, passwordHash, err := model.FindCredentialsByUsername(req.Username)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)) != nil {
		writeAuthFailure(conn, dto.MessageLoginResponse, requestID, "invalid_credentials", "Invalid credentials")
		return
	}
	session, err := CreateSession(context.Background(), user.AuthID)
	if err != nil {
		log.Printf("[ERROR][AUTH] Create session failed: %v", err)
		writeAuthFailure(conn, dto.MessageLoginResponse, requestID, "session_creation_failed", "Could not create session")
		return
	}
	user.LastLogin = time.Now()
	if err := model.SaveUser(&user); err != nil {
		log.Printf("[ERROR][AUTH] Update login time failed: %v", err)
	}
	bindIdentity(conn, connectionIdentity{AuthID: user.AuthID, Username: user.Username})
	_ = wsconn.Send(conn, dto.OK(dto.MessageLoginResponse, requestID, "Login successful", dto.Session{SessionID: session.SessionID, ExpiresAt: session.ExpiresAt.Unix()}))
}

func handleAuthenticate(conn *websocket.Conn, requestID string, data json.RawMessage) {
	var req dto.SessionRequest
	if err := json.Unmarshal(data, &req); err != nil || req.SessionID == "" {
		writeAuthFailure(conn, dto.MessageAuthenticateResponse, requestID, "session_required", "Session ID is required")
		return
	}
	_, user, err := FindSessionByID(req.SessionID)
	if err != nil || !user.IsActive {
		writeAuthFailure(conn, dto.MessageAuthenticateResponse, requestID, "invalid_session", "Session is invalid or expired")
		return
	}
	bindIdentity(conn, connectionIdentity{AuthID: user.AuthID, Username: user.Username})
	_ = wsconn.Send(conn, dto.OK(dto.MessageAuthenticateResponse, requestID, "Authenticated", dto.Empty{}))
}

func handleLogout(conn *websocket.Conn, requestID string, data json.RawMessage) {
	var req dto.SessionRequest
	if err := json.Unmarshal(data, &req); err != nil || req.SessionID == "" {
		writeAuthFailure(conn, dto.MessageLogoutResponse, requestID, "session_required", "Session ID is required")
		return
	}
	if err := DeleteSession(context.Background(), req.SessionID); err != nil {
		writeAuthFailure(conn, dto.MessageLogoutResponse, requestID, "logout_failed", "Logout failed")
		return
	}
	removeIdentity(conn)
	_ = wsconn.Send(conn, dto.OK(dto.MessageLogoutResponse, requestID, "Logged out", dto.Empty{}))
}

func handleGetUser(conn *websocket.Conn, requestID string, data json.RawMessage) {
	identity, ok := getIdentity(conn)
	if !ok {
		var req dto.SessionRequest
		if json.Unmarshal(data, &req) == nil && req.SessionID != "" {
			_, user, err := FindSessionByID(req.SessionID)
			if err == nil {
				identity, ok = connectionIdentity{AuthID: user.AuthID, Username: user.Username}, true
				bindIdentity(conn, identity)
			}
		}
	}
	if !ok {
		writeAuthFailure(conn, dto.MessageUserResponse, requestID, "authentication_required", "Authentication required")
		return
	}
	user, err := model.FindUserByAuthID(identity.AuthID)
	if err != nil {
		writeAuthFailure(conn, dto.MessageUserResponse, requestID, "user_not_found", "User not found")
		return
	}
	_ = wsconn.Send(conn, dto.OK(dto.MessageUserResponse, requestID, "", dto.UserData{User: dto.ToUser(&user), MaxEXP: model.GetMaxExp(user.Level)}))
}
