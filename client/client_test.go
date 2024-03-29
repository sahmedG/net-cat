// client_test.go

package client

import (
	"fmt"
	"net"
	"testing"
)

func TestClientConnection(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to start listener: %v", err)
	}
	defer listener.Close()

	serverAddr := listener.Addr().String()

	go func() {
		conn, _ := listener.Accept()
		defer conn.Close()

		name := "TestUser\n"
		conn.Write([]byte(name))

		fmt.Fprintln(conn, "Hello, server!")
	}()

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	response := make([]byte, len("Welcome to TCP-Chat!\n"))
	n, err := conn.Read(response)
	if err != nil {
		t.Fatalf("Failed to read from server: %v", err)
	}

	expectedResponse := "Welcome to TCP-Chat!\n"
	if string(response[:n]) != expectedResponse {
		t.Errorf("Expected response: %s, but got: %s", expectedResponse, string(response[:n]))
	}
}
