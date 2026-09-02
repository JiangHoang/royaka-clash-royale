package wsconn

import (
	"errors"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 65 * time.Second
	pingPeriod     = 30 * time.Second
	maxMessageSize = 1 << 20
)

var ErrClosed = errors.New("websocket connection is closed")

type Client struct {
	conn      *websocket.Conn
	outbound  chan any
	done      chan struct{}
	closeOnce sync.Once
}

var registry = struct {
	sync.RWMutex
	clients map[*websocket.Conn]*Client
}{clients: make(map[*websocket.Conn]*Client)}

func Register(conn *websocket.Conn) *Client {
	client := &Client{conn: conn, outbound: make(chan any, 256), done: make(chan struct{})}
	registry.Lock()
	registry.clients[conn] = client
	registry.Unlock()
	go client.writePump()
	return client
}

func ConfigureReader(conn *websocket.Conn) {
	conn.SetReadLimit(maxMessageSize)
	_ = conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(pongWait))
	})
}

func Send(conn *websocket.Conn, value any) error {
	registry.RLock()
	client := registry.clients[conn]
	registry.RUnlock()
	if client == nil {
		return ErrClosed
	}
	return client.Send(value)
}

func (c *Client) Send(value any) error {
	select {
	case <-c.done:
		return ErrClosed
	case c.outbound <- value:
		return nil
	}
}

func Close(conn *websocket.Conn) {
	registry.Lock()
	client := registry.clients[conn]
	delete(registry.clients, conn)
	registry.Unlock()
	if client != nil {
		client.close()
	} else {
		_ = conn.Close()
	}
}

func (c *Client) close() {
	c.closeOnce.Do(func() {
		registry.Lock()
		delete(registry.clients, c.conn)
		registry.Unlock()
		close(c.done)
		_ = c.conn.Close()
	})
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case value := <-c.outbound:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteJSON(value); err != nil {
				c.close()
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.close()
				return
			}
		case <-c.done:
			return
		}
	}
}
