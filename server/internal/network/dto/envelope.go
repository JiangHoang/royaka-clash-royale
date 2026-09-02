package dto

import "encoding/json"

// Envelope is the common inbound WebSocket message shape. RequestID is
// optional so older clients remain compatible during the migration.
type Envelope struct {
	Type      MessageType     `json:"type"`
	RequestID string          `json:"request_id,omitempty"`
	Data      json.RawMessage `json:"data,omitempty"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Response[T any] struct {
	Type      MessageType `json:"type"`
	RequestID string      `json:"request_id,omitempty"`
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Error     *Error      `json:"error"`
	Data      T           `json:"data,omitempty"`
}

// Event uses the same wire shape as a response, but intentionally has no
// request ID because it is initiated by the server.
type Event[T any] struct {
	Type    MessageType `json:"type"`
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Error   *Error      `json:"error"`
	Data    T           `json:"data,omitempty"`
}

func OK[T any](messageType MessageType, requestID, message string, data T) Response[T] {
	return Response[T]{Type: messageType, RequestID: requestID, Success: true, Message: message, Data: data}
}

func Fail(messageType MessageType, requestID, code, message string) Response[any] {
	return Response[any]{
		Type:      messageType,
		RequestID: requestID,
		Success:   false,
		Message:   message,
		Error:     &Error{Code: code, Message: message},
	}
}

func Push[T any](messageType MessageType, message string, data T) Event[T] {
	return Event[T]{Type: messageType, Success: true, Message: message, Data: data}
}
