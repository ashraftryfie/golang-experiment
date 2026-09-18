package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"
)

func handleClient(conn net.Conn) {
	defer conn.Close()
	remoteAddr := conn.RemoteAddr().String()
	fmt.Printf("[Server] Client connected: %s\n", remoteAddr)

	scanner := bufio.NewScanner(conn)
	for {
		// Set sliding read deadline of 30 seconds
		_ = conn.SetReadDeadline(time.Now().Add(30 * time.Second))

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Printf("[Server] Error reading from %s: %v\n", remoteAddr, err)
			} else {
				fmt.Printf("[Server] Client %s disconnected gracefully\n", remoteAddr)
			}
			return
		}

		line := scanner.Text()
		fmt.Printf("[Server] Received from %s: %s\n", remoteAddr, line)

		// Echo back with newline
		_, err := io.WriteString(conn, fmt.Sprintf("ECHO: %s\n", line))
		if err != nil {
			fmt.Printf("[Server] Error writing to %s: %v\n", remoteAddr, err)
			return
		}
	}
}

func main() {
	// Listen on dynamic localhost TCP port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	serverAddr := listener.Addr().String()
	fmt.Printf("[Server] Echo server listening on %s\n", serverAddr)

	// Launch test client in goroutine to demonstrate communication
	go func() {
		time.Sleep(50 * time.Millisecond)
		conn, dialErr := net.Dial("tcp", serverAddr)
		if dialErr != nil {
			fmt.Printf("[Client] Dial error: %v\n", dialErr)
			return
		}
		defer conn.Close()

		_, _ = io.WriteString(conn, "Hello Go Sockets!\n")
		resp, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Printf("[Client] Server response: %s", resp)
	}()

	// Accept single connection for demonstration
	conn, err := listener.Accept()
	if err != nil {
		fmt.Printf("[Server] Accept error: %v\n", err)
		return
	}
	handleClient(conn)
}
