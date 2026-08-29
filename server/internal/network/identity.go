package network

import (
	"sync"

	"github.com/gorilla/websocket"
)

type connectionIdentity struct {
	AuthID   string
	Username string
	Legacy   bool
}

var identities = struct {
	sync.RWMutex
	values map[*websocket.Conn]connectionIdentity
}{values: make(map[*websocket.Conn]connectionIdentity)}

func bindIdentity(conn *websocket.Conn, identity connectionIdentity) {
	identities.Lock()
	identities.values[conn] = identity
	identities.Unlock()
}

func getIdentity(conn *websocket.Conn) (connectionIdentity, bool) {
	identities.RLock()
	identity, ok := identities.values[conn]
	identities.RUnlock()
	return identity, ok
}

func removeIdentity(conn *websocket.Conn) {
	identities.Lock()
	delete(identities.values, conn)
	identities.Unlock()
}
