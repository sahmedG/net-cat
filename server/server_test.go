// server_test.go

package server

import (
	"net"
	"testing"
)
var maxConnections = 10
func TestHandleClient_MaxConnections(t *testing.T) {
	fakeListener, _ := net.Listen("tcp", "127.0.0.1:0")
	defer fakeListener.Close()

	for i := 0; i <= maxConnections; i++ {
		conn, _ := net.Dial("tcp", fakeListener.Addr().String())
		go handleClient(conn)
	}

	conn, err := net.Dial("tcp", fakeListener.Addr().String())
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	conn.Close()
}
