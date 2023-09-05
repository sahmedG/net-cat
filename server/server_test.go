// server_test.go

package server

import (
	"net"
	"testing"
)
var maxConnections = 10
func TestHandleClient_MaxConnections(t *testing.T) {
	// Create a fake listener for testing
	fakeListener, _ := net.Listen("tcp", "127.0.0.1:0")
	defer fakeListener.Close()

	// Create multiple fake client connections to exceed the maximum limit
	for i := 0; i <= maxConnections; i++ {
		conn, _ := net.Dial("tcp", fakeListener.Addr().String())
		go handleClient(conn)
	}

	// Attempt to create one more connection
	conn, err := net.Dial("tcp", fakeListener.Addr().String())
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Close the connection to clean up
	conn.Close()
}
