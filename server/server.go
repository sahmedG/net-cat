package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"time"
)

func StartServer(port string) {
	// fmt.Println(GetLocalIP())

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
		if ConnCount >= MaxConn {
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
		no_get: true,
	}

	logo := []string{
		"Welcome to TCP-Chat!",
		"         _nnnn_",
		"        dGGGGMMb",
		"       @p~qp~~qMb",
		"       M|@||@) M|",
		"       @,----.JM|",
		"      JS^\\__/  qKL",
		"     dZP        qKRb",
		"    dZP          qKKb",
		"   fZP            SMMb",
		"   HZM            MMMM",
		"   FqM            MMMM",
		" __| \".        |\\dS\"qML",
		" |    `.       | `' \\Zq",
		"_)      \\.___.,|     .'",
		"\\____   )MMMMMP|   .'",
		"     `-'       `--'",
	}

	for _, line := range logo {
		client.Writer.WriteString(line + "\n")
	}
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
		Broadcast(ChatMsg{
			Sender:  "System",
			Content: fmt.Sprintf("User %s left the chat", client.Name),
			RoomID:  client.RoomID,
			Time:    time.Now(),
		})
	}()

loop:
	client.Writer.WriteString("\nEnter your name: ")
	client.Writer.Flush()
	name, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Println(err)
		return
	}
	client.Name = strings.TrimSpace(name)
	if client.Name == "" {
		client.Writer.WriteString("\nCan't enter empty name")
		goto loop
	}

	/* unblock user to get messeges */
	clientsLock.Lock()
	client.no_get = false
	clientsLock.Unlock()

	go func() {
		Broadcast(ChatMsg{
			Sender:  "System",
			Content: fmt.Sprintf("User %s joined the chat", client.Name),
			RoomID:  "",
			Time:    time.Now(),
		})
	}()
	go func() {
		messages := chatRooms[""]
		for _, message := range messages {
			if message.Sender == "System" {
				continue
			} else {
				client.Writer.WriteString(formatMessage(message) + "\n")
			}
		}

		for i := 0; i < len(messages); i++ {
			if messages[i].Sender == "System" {
				continue
			}
			if i == len(messages)-1 {
				client.Writer.WriteString("------chat history-----------\n")
			}
		}
	}()
	//go pingClient(client, conn)
	go func() {
		/*  */
		scanner := bufio.NewScanner(conn)

		/* Scanner seems to keep running until client disconnects */
		for scanner.Scan() {
			//go pingClient(client, conn)
			/* get the message from the client */
			message := scanner.Text()
			/* If the message is not empty */
			if message != "" {
				/* let him join a room when he type /join */
				if strings.HasPrefix(message, "/join") {
					/* wrong input */
					args := strings.Fields(message)
					if len(args) != 2 {
						client.Writer.WriteString("[System] Invalid usage. Use: /join <room_id>\n")
						client.Writer.Flush()
						continue
					}
					/* Join room with ID */
					roomID := args[1]
					if rune(roomID[0]) >= 49 && rune(roomID[0]) <= 57 {
						client.JoinRoom(roomID)
					} else {
						client.Writer.WriteString("[System] Invalid usage. Use: /join <room_id> ie 1 - 9\n")
						client.Writer.Flush()
						continue
					}

					/* let the client leave the room when he type /leave */
				} else if message == "/leave" {
					client.LeaveRoom()
				} else if strings.HasPrefix(message, "/help") {
					client.Writer.WriteString("\nProgram usage: /join [room number], /leave, /rn [new user name], /exit\n")
					client.Writer.Flush()
					continue
				} else if strings.HasPrefix(message, "/rn") {
					args := strings.Fields(message)
					if !Renaming_Arg_check(args) {
						client.Writer.WriteString("Usage: /rn [Name]\n")
						client.Writer.Flush()
						continue
					}
					new_name := args[1]
					client.ChangeName(new_name)
				} else if strings.HasPrefix(message, "/exit") {
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
					break
				} else {
					Broadcast(ChatMsg{
						Sender:  client.Name,
						Content: message,
						RoomID:  client.RoomID,
						Time:    time.Now(),
					})
				}
			}
		}
		client.ClientExit(conn)
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		s := <-c
		fmt.Println(s)
		os.Exit(1)
	}()
	clientsLock.Lock()
	for _, c := range clients {
		fmt.Printf("c: %v\n", c)
	}
	clientsLock.Unlock()
	select {}
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

func formatMessage(message ChatMsg) string {
	return fmt.Sprintf("[%s][%s]: %s", message.Time.Format("2006-01-02 15:04:05"), message.Sender, message.Content)
}
