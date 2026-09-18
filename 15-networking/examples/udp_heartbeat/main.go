package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	// Listen on dynamic UDP port
	serverConn, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	defer serverConn.Close()

	serverAddr := serverConn.LocalAddr().String()
	fmt.Printf("[UDP Server] Listening on %s\n", serverAddr)

	// Sender client goroutine
	go func() {
		time.Sleep(50 * time.Millisecond)
		clientConn, dialErr := net.Dial("udp", serverAddr)
		if dialErr != nil {
			fmt.Printf("[UDP Client] Dial error: %v\n", dialErr)
			return
		}
		defer clientConn.Close()

		for i := 1; i <= 3; i++ {
			msg := fmt.Sprintf("HEARTBEAT #%d", i)
			_, _ = clientConn.Write([]byte(msg))
			time.Sleep(50 * time.Millisecond)
		}
	}()

	// Read 3 heartbeats
	buf := make([]byte, 1024)
	for i := 1; i <= 3; i++ {
		_ = serverConn.SetReadDeadline(time.Now().Add(1 * time.Second))
		n, peerAddr, readErr := serverConn.ReadFrom(buf)
		if readErr != nil {
			fmt.Printf("[UDP Server] Read error: %v\n", readErr)
			return
		}
		fmt.Printf("[UDP Server] Received %d bytes from %s: %s\n", n, peerAddr, string(buf[:n]))
	}

	fmt.Println("[UDP Server] Finished processing heartbeats successfully.")
}
