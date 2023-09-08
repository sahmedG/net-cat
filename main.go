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
		if len(os.Args) == 3{
			port := os.Args[2]
			server.StartServer(port)
		} else {
			port := "8989"
			server.StartServer(port)
		}

	case "client":
		if len(os.Args) == 3 {
			client.StartClient(os.Args[2])
		} else {
			fmt.Println("Missing the host name and port!")
		}
	default:
		fmt.Println("Invalid mode. Use 'server' or 'client'.")
	}
}
