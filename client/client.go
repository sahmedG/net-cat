package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func StartClient(hostname string) {

	conn, err := net.Dial("tcp", hostname)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	conn.Write([]byte(name + "\n"))

	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
			if scanner.Text() == "Can't accept anymore connections at the moment!" {
				conn.Close()
				break
			}
		}
	}()

	for {
		message, _ := reader.ReadString('\n')
		if strings.HasPrefix(message,"/help"){
			conn.Write([]byte("you can use the following commands within the program: /join -> followed by room number to join a room /leave -> to leave the current room /exit -> to exit the program"))
		}
		
		if strings.HasPrefix(message, "/exit\n") {
			conn.Write([]byte("User "+name+" left the chat."))
			break
		}

		if strings.HasPrefix(message, "/join") || message == "/leave\n" {
			conn.Write([]byte(message))
		} else {
			conn.Write([]byte(message))
		}
	}
}
