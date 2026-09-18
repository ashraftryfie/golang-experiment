package chatbroker

import (
	"bufio"
	"io"
	"net"
	"strings"
	"testing"
	"time"
)

func TestChatServer_JoinAndBroadcast(t *testing.T) {
	server := NewChatServer()
	if err := server.Start("127.0.0.1:0"); err != nil {
		t.Fatalf("failed to start server: %v", err)
	}
	defer server.Stop()

	serverAddr := server.Addr().String()

	// Client 1: Alice
	conn1, err := net.Dial("tcp", serverAddr)
	if err != nil {
		t.Fatalf("Alice failed to connect: %v", err)
	}
	defer conn1.Close()
	reader1 := bufio.NewReader(conn1)

	// Send nickname
	_, _ = io.WriteString(conn1, "alice\n")
	status1, _ := reader1.ReadString('\n')
	if strings.TrimSpace(status1) != "OK" {
		t.Fatalf("Alice expected OK, got %s", status1)
	}

	// Client 2: Bob
	conn2, err := net.Dial("tcp", serverAddr)
	if err != nil {
		t.Fatalf("Bob failed to connect: %v", err)
	}
	defer conn2.Close()
	reader2 := bufio.NewReader(conn2)

	// Send nickname
	_, _ = io.WriteString(conn2, "bob\n")
	status2, _ := reader2.ReadString('\n')
	if strings.TrimSpace(status2) != "OK" {
		t.Fatalf("Bob expected OK, got %s", status2)
	}

	// Alice should receive Bob's JOIN announcement
	joinMsg, err := reader1.ReadString('\n')
	if err != nil {
		t.Fatalf("Alice failed reading join message: %v", err)
	}
	if strings.TrimSpace(joinMsg) != "JOIN: bob" {
		t.Errorf("Alice expected 'JOIN: bob', got '%s'", strings.TrimSpace(joinMsg))
	}

	// Bob sends a message
	_, _ = io.WriteString(conn2, "Hello Alice!\n")

	// Alice should receive Bob's message
	chatMsg, err := reader1.ReadString('\n')
	if err != nil {
		t.Fatalf("Alice failed reading chat message: %v", err)
	}
	expectedMsg := "MSG: bob: Hello Alice!"
	if strings.TrimSpace(chatMsg) != expectedMsg {
		t.Errorf("Alice expected '%s', got '%s'", expectedMsg, strings.TrimSpace(chatMsg))
	}
}

func TestChatServer_DuplicateNicknameRejected(t *testing.T) {
	server := NewChatServer()
	if err := server.Start("127.0.0.1:0"); err != nil {
		t.Fatalf("failed to start server: %v", err)
	}
	defer server.Stop()

	serverAddr := server.Addr().String()

	// Connect first Alice
	conn1, err := net.Dial("tcp", serverAddr)
	if err != nil {
		t.Fatalf("Alice 1 failed to connect: %v", err)
	}
	defer conn1.Close()
	_, _ = io.WriteString(conn1, "alice\n")
	resp1, _ := bufio.NewReader(conn1).ReadString('\n')
	if strings.TrimSpace(resp1) != "OK" {
		t.Fatalf("Alice 1 expected OK, got %s", resp1)
	}

	// Connect second Alice
	conn2, err := net.Dial("tcp", serverAddr)
	if err != nil {
		t.Fatalf("Alice 2 failed to connect: %v", err)
	}
	defer conn2.Close()
	_, _ = io.WriteString(conn2, "alice\n")
	resp2, _ := bufio.NewReader(conn2).ReadString('\n')
	if !strings.HasPrefix(resp2, "ERR: nickname already taken") {
		t.Fatalf("Alice 2 expected duplicate error, got %s", resp2)
	}
}

func TestChatServer_QuitAndLeaveBroadcast(t *testing.T) {
	server := NewChatServer()
	if err := server.Start("127.0.0.1:0"); err != nil {
		t.Fatalf("failed to start server: %v", err)
	}
	defer server.Stop()

	serverAddr := server.Addr().String()

	conn1, _ := net.Dial("tcp", serverAddr)
	defer conn1.Close()
	_, _ = io.WriteString(conn1, "watcher\n")
	reader1 := bufio.NewReader(conn1)
	_, _ = reader1.ReadString('\n') // read OK

	conn2, _ := net.Dial("tcp", serverAddr)
	defer conn2.Close()
	_, _ = io.WriteString(conn2, "quitter\n")
	reader2 := bufio.NewReader(conn2)
	_, _ = reader2.ReadString('\n') // read OK

	// watcher reads quitter JOIN
	_, _ = reader1.ReadString('\n')

	// quitter quits
	_, _ = io.WriteString(conn2, "/quit\n")
	bye, _ := reader2.ReadString('\n')
	if strings.TrimSpace(bye) != "BYE" {
		t.Fatalf("quitter expected BYE, got %s", bye)
	}

	// watcher should receive LEAVE
	leaveMsg, err := reader1.ReadString('\n')
	if err != nil {
		t.Fatalf("watcher failed reading leave: %v", err)
	}
	if strings.TrimSpace(leaveMsg) != "LEAVE: quitter" {
		t.Fatalf("watcher expected 'LEAVE: quitter', got '%s'", strings.TrimSpace(leaveMsg))
	}
}

func TestChatServer_GracefulShutdown(t *testing.T) {
	server := NewChatServer()
	if err := server.Start("127.0.0.1:0"); err != nil {
		t.Fatalf("failed to start server: %v", err)
	}

	serverAddr := server.Addr().String()
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		t.Fatalf("dial failed: %v", err)
	}
	defer conn.Close()

	_, _ = io.WriteString(conn, "user1\n")
	_, _ = bufio.NewReader(conn).ReadString('\n') // OK

	// Stop server
	if err := server.Stop(); err != nil {
		t.Fatalf("Stop() error: %v", err)
	}

	// Any subsequent dial must fail
	_, dialErr := net.DialTimeout("tcp", serverAddr, 50*time.Millisecond)
	if dialErr == nil {
		t.Fatalf("expected dial to fail on stopped server")
	}
}
