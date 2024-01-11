// main.go

package main

import (
	"TCPChat/server"
	"fmt"
	"os"
)

func main() {
	args := os.Args

	if len(args) == 2 {
		server.StartServer(os.Args[1])
	} else if len(args) == 1 {
		server.StartServer("8989")
	} else {
		fmt.Println("[USAGE]: ./TCPChat [$port]")
	}
}
