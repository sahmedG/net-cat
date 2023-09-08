package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	MaxConn       = 10
	ConnCount     = 0
	ConnCountLock sync.Mutex
)

type Client struct {
	Name   string
	Conn   net.Conn
	Writer *bufio.Writer
	RoomID string
}

type ChatMsg struct {
	Sender  string
	Content string
	RoomID  string
	Time    time.Time
}

var (
	clients     = make(map[*Client]bool)
	clientsLock sync.Mutex
)

var (
	chatRooms     = make(map[string][]ChatMsg)
	chatRoomsLock sync.Mutex
)

func StartServer(port string) {
	fmt.Println(GetLocalIP())

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Printf("Server listening on port %s\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		ConnCountLock.Lock()
		if ConnCount >= MaxConn{
			log.Println("Reached maximum number of connections.")
			conn.Write([]byte("Can't accept anymore connections at the moment!"))
			conn.Close()
			ConnCountLock.Unlock()
			continue
		}
		ConnCount++
		ConnCountLock.Unlock()

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {

	defer conn.Close()
	client := &Client{
		Conn:   conn,
		Writer: bufio.NewWriter(conn),
	}

	logo := `
	Welcome to TCP-Chat!
	_nnnn_
	dGGGGMMb
   @p~qp~~qMb
   M|@||@) M|
   @,----.JM|
  JS^\__/  qKL
 dZP        qKRb
dZP          qKKb
fZP            SMMb
HZM            MMMM
FqM            MMMM
__| ".        |\dS"qML
|    ".       | ' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
 	 '-'       '--`
	client.Writer.WriteString(logo)
	client.Writer.Flush()

	clientsLock.Lock()
	clients[client] = true
	clientsLock.Unlock()

	defer func() {
		ConnCountLock.Lock()
		ConnCount--
		ConnCountLock.Unlock()
		clientsLock.Lock()
		delete(clients, client)
		clientsLock.Unlock()
		broadcast(ChatMsg{
			Sender:  "System",
			Content: fmt.Sprintf("User %s left the chat", client.Name),
			RoomID:  client.RoomID,
			Time:    time.Now(),
		})
	}()

	client.Writer.WriteString("\nEnter your name: ")
	client.Writer.Flush()
	name, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Println(err)
		return
	}
	client.Name = strings.TrimSpace(name)

	broadcast(ChatMsg{
		Sender:  "System",
		Content: fmt.Sprintf("User %s joined the chat", client.Name),
		RoomID:  "",
		Time:    time.Now(),
	})


	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			message := scanner.Text()
			if message != "" {
				if strings.HasPrefix(message, "/join") {
					args := strings.Fields(message)
					if len(args) != 2 {
						client.Writer.WriteString("[System] Invalid usage. Use: /join <room_id>\n")
						client.Writer.Flush()
						continue
					}
					roomID := args[1]
					joinRoom(client, roomID)
				} else if message == "/leave" {
					leaveRoom(client)
				} else {
					broadcast(ChatMsg{
						Sender:  client.Name,
						Content: message,
						RoomID:  client.RoomID,
						Time:    time.Now(),
					})
				}
			}
		}
	}()

	clientsLock.Lock()
	for _, c := range clients {
		fmt.Printf("c: %v\n", c)
	}
	clientsLock.Unlock()
	select {}
}

func broadcast(message ChatMsg) {
	clientsLock.Lock()
	defer clientsLock.Unlock()

	logFile, err := os.OpenFile("chat_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening log file:", err)
		return
	}
	defer logFile.Close()

	_, err = logFile.WriteString(message.Content + "\n")
	if err != nil {
		log.Println("Error writing to log file:", err)
		return
	}

	chatRooms[message.RoomID] = append(chatRooms[message.RoomID], message)

	for client := range clients {
		if client.RoomID == message.RoomID {
			client.Writer.WriteString(formatMessage(message) + "\n")
			client.Writer.Flush()
		}
	}
}

func GetLocalIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddress := conn.LocalAddr().(*net.UDPAddr)

	return localAddress.IP
}

func joinRoom(client *Client, roomID string) {
	leaveRoom(client)
	client.RoomID = roomID
	broadcast(ChatMsg{
		Sender:  "System",
		Content: client.Name + " joined the room.",
		RoomID:  roomID,
		Time:    time.Now(),
	})
	chatRoomsLock.Lock()
	messages := chatRooms[roomID]
	chatRoomsLock.Unlock()
	for _, message := range messages {
		client.Writer.WriteString(formatMessage(message) + "\n")
		client.Writer.Flush()
	}
}

func formatMessage(message ChatMsg) string {
	return fmt.Sprintf("[%s][%s]: %s", message.Time.Format("2006-01-02 15:04:05"), message.Sender, message.Content)
}

func leaveRoom(client *Client) {
	if client.RoomID != "" {
		broadcast(ChatMsg{
			Sender:  "System",
			Content: client.Name + " left the room.",
			RoomID:  client.RoomID,
			Time:    time.Now(),
		})
		client.RoomID = ""
	}
}
