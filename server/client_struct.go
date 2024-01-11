package server

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

/* holds clients informations */
type Client struct {
	Name   string
	Conn   net.Conn
	Writer *bufio.Writer
	RoomID string
	no_get bool // if true, it will block user from recieving messeges.
}

/* This function will allow client to join a group chat (room), will update the RoomID variable */
func (client *Client) JoinRoom(roomID string) {
	client.LeaveRoom()
	client.RoomID = roomID
	Broadcast(ChatMsg{
		Sender:  "System",
		Content: client.Name + " joined the room.",
		RoomID:  roomID,
		Time:    time.Now(),
	})

	chatRoomsLock.Lock()
	messages := chatRooms[roomID]
	chatRoomsLock.Unlock()
	fmt.Println(len(messages))
	for i := 0; i < len(messages); i++ {
		if messages[i].Sender == "System" {
			continue
		}
		client.Writer.WriteString(formatMessage(messages[i]) + "\n")
	}
	if len(messages) > 1 {
		client.Writer.WriteString("------chat history-----------\n")
	}
	client.Writer.Flush()
}

/* Resets RoomID of the client */
func (client *Client) LeaveRoom() {

	if client.RoomID != "" {
		Broadcast(ChatMsg{
			Sender:  "System",
			Content: client.Name + " left the room.",
			RoomID:  client.RoomID,
			Time:    time.Now(),
		})
		client.RoomID = ""
	}
}

func (client *Client) ChangeName(new_name string) {

	clientsLock.Lock()
	old_name := client.Name
	client.Name = new_name
	clientsLock.Unlock()

	Broadcast(ChatMsg{
		Sender:  "System",
		Content: old_name + " changed his name to " + client.Name,
		RoomID:  client.RoomID,
		Time:    time.Now(),
	})
	client.RoomID = ""
}

func (client *Client) ClientExit(conn net.Conn) {
	ConnCountLock.Lock()
	ConnCount--
	ConnCountLock.Unlock()
	clientsLock.Lock()
	delete(clients, client)
	clientsLock.Unlock()
	Broadcast(ChatMsg{
		Sender:  "System",
		Content: fmt.Sprintf("User %s left the chat", client.Name),
		RoomID:  "",
		Time:    time.Now(),
	})
	conn.Close()
}
