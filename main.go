// main.go

package main

import (
	"fmt"
	"netcat/client"
	"netcat/server"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Usage: ./main server/client")
		return
	}

	mode := args[1]

	switch mode {
	case "server":
		server.StartServer()
	case "client":
		client.StartClient(os.Args[2])
	default:
		fmt.Println("Invalid mode. Use 'server' or 'client'.")
	}
}
