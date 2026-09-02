package dto

// MessageType is the closed set of message names supported by the WebSocket
// protocol. Go has no enum keyword, so a named string type plus constants is
// the idiomatic equivalent and keeps the JSON representation unchanged.
type MessageType string

const (
	// Client requests.
	MessageRegister     MessageType = "register"
	MessageLogin        MessageType = "login"
	MessageAuthenticate MessageType = "authenticate"
	MessageLogout       MessageType = "logout"
	MessageGetUser      MessageType = "get_user"
	MessageGetDesk      MessageType = "get_desk"
	MessageFindMatch    MessageType = "find_match"
	MessageGetGame      MessageType = "get_game"
	MessageAttack       MessageType = "attack"
	MessageHeal         MessageType = "heal"
	MessageSkipTurn     MessageType = "skip_turn"
	MessagePlayAgain    MessageType = "play_again"
	MessageLeaveGame    MessageType = "leave_game"
	MessageSelectTroop  MessageType = "select_troop"

	// Direct responses and server-pushed events.
	MessageError                MessageType = "error"
	MessageRegisterResponse     MessageType = "register_response"
	MessageLoginResponse        MessageType = "login_response"
	MessageAuthenticateResponse MessageType = "authenticate_response"
	MessageLogoutResponse       MessageType = "logout_response"
	MessageUserResponse         MessageType = "user_response"
	MessageDeckResponse         MessageType = "deck_response"
	MessageFindMatchResponse    MessageType = "find_match_response"
	MessageMatchFound           MessageType = "match_found"
	MessageMatchTimeout         MessageType = "match_timeout"
	MessageGameResponse         MessageType = "game_response"
	MessageAttackResponse       MessageType = "attack_response"
	MessageHealResponse         MessageType = "heal_response"
	MessageSkipTurnResponse     MessageType = "skip_turn_response"
	MessagePlayAgainResponse    MessageType = "play_again_response"
	MessageLeaveGameResponse    MessageType = "leave_game_response"
	MessageTroopResponse        MessageType = "troop_response"
	MessageManaUpdate           MessageType = "mana_update"
	MessageGameState            MessageType = "game_state"
	MessageGameOverResponse     MessageType = "game_over_response"
)

func (messageType MessageType) String() string { return string(messageType) }

func (messageType MessageType) IsRequest() bool {
	switch messageType {
	case MessageRegister, MessageLogin, MessageAuthenticate, MessageLogout,
		MessageGetUser, MessageGetDesk, MessageFindMatch, MessageGetGame,
		MessageAttack, MessageHeal, MessageSkipTurn, MessagePlayAgain,
		MessageLeaveGame, MessageSelectTroop:
		return true
	default:
		return false
	}
}
