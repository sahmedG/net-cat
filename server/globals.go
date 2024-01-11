package server

import "sync"

/* connections lock and count */
var (
	MaxConn       = 10        // Max connections
	ConnCount     = 0        // Current connections count
	ConnCountLock sync.Mutex // Provides locking for max connections
)

var (
	clients     = make(map[*Client]bool)
	clientsLock sync.Mutex // Locks the client from accesing variables
)

var (
	chatRooms     = make(map[string][]ChatMsg)
	chatRoomsLock sync.Mutex
)
