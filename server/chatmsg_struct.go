package server

import (
	"time"
)

/* Holds message informations */
type ChatMsg struct {
	Sender  string
	Content string
	RoomID  string
	Time    time.Time
}
