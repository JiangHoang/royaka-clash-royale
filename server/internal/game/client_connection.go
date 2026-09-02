package game

import (
	"royaka/internal/network/wsconn"

	"github.com/gorilla/websocket"
)

type ClientConnection struct {
	Conn     *websocket.Conn
	Username string
}

func (c *ClientConnection) SafeWrite(data interface{}) error {
	return wsconn.Send(c.Conn, data)
}
