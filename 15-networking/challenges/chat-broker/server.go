package chatbroker

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

var (
	ErrNicknameTaken = errors.New("nickname already taken")
	ErrServerClosed  = errors.New("server is closed")
)

type client struct {
	nick string
	conn net.Conn
	send chan string
}

type ChatServer struct {
	listener net.Listener
	clients  map[string]*client
	mu       sync.RWMutex

	broadcastCh chan string
	registerCh  chan *client
	leaveCh     chan *client
	shutdownCh  chan struct{}

	wg sync.WaitGroup
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		clients:     make(map[string]*client),
		broadcastCh: make(chan string, 128),
		registerCh:  make(chan *client),
		leaveCh:     make(chan *client),
		shutdownCh:  make(chan struct{}),
	}
}

func (s *ChatServer) Start(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on tcp %s: %w", addr, err)
	}
	s.listener = l

	s.wg.Add(2)
	go s.brokerLoop()
	go s.acceptLoop()

	return nil
}

func (s *ChatServer) Addr() net.Addr {
	if s.listener == nil {
		return nil
	}
	return s.listener.Addr()
}

func (s *ChatServer) acceptLoop() {
	defer s.wg.Done()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.shutdownCh:
				return
			default:
				// If accept error occurs, continue or return if fatal
				return
			}
		}

		s.wg.Add(1)
		go s.handleConn(conn)
	}
}

func (s *ChatServer) handleConn(conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Step 1: Handshake for nickname with 10s deadline
	_ = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	rawNick, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	nick := strings.TrimSpace(rawNick)
	if nick == "" {
		_, _ = io.WriteString(conn, "ERR: empty nickname\n")
		return
	}

	c := &client{
		nick: nick,
		conn: conn,
		send: make(chan string, 64),
	}

	// Register client safely
	s.mu.Lock()
	if _, exists := s.clients[nick]; exists {
		s.mu.Unlock()
		_, _ = io.WriteString(conn, "ERR: nickname already taken\n")
		return
	}
	s.clients[nick] = c
	s.mu.Unlock()

	// Clear initial handshake deadline
	_ = conn.SetReadDeadline(time.Time{})

	// Confirm registration to client
	_, _ = io.WriteString(conn, "OK\n")

	// Broadcast join to other clients
	s.broadcast(fmt.Sprintf("JOIN: %s\n", nick), nick)

	// Start writer goroutine for this client
	clientDone := make(chan struct{})
	go func() {
		defer close(clientDone)
		for msg := range c.send {
			_, writeErr := io.WriteString(conn, msg)
			if writeErr != nil {
				return
			}
		}
	}()

	defer func() {
		// Cleanup client on departure
		s.mu.Lock()
		delete(s.clients, nick)
		s.mu.Unlock()

		close(c.send)
		<-clientDone

		s.broadcast(fmt.Sprintf("LEAVE: %s\n", nick), nick)
	}()

	// Step 2: Read loop
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		msg := strings.TrimRight(line, "\r\n")
		if msg == "/quit" {
			_, _ = io.WriteString(conn, "BYE\n")
			return
		}

		// Broadcast message to everyone else
		s.broadcast(fmt.Sprintf("MSG: %s: %s\n", nick, msg), nick)
	}
}

// broadcast sends a message to all connected clients except optional excludeNick.
func (s *ChatServer) broadcast(msg string, excludeNick string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for nick, c := range s.clients {
		if excludeNick != "" && nick == excludeNick {
			continue
		}
		select {
		case c.send <- msg:
		default:
			// If client buffer is full, drop or skip to prevent blocking broker
		}
	}
}

func (s *ChatServer) brokerLoop() {
	defer s.wg.Done()
	<-s.shutdownCh
}

// Stop cleanly terminates the chat broker and disconnects all clients.
func (s *ChatServer) Stop() error {
	select {
	case <-s.shutdownCh:
		return ErrServerClosed
	default:
		close(s.shutdownCh)
	}

	var closeErr error
	if s.listener != nil {
		closeErr = s.listener.Close()
	}

	// Close all client network connections to unblock active readers
	s.mu.Lock()
	for _, c := range s.clients {
		_ = c.conn.Close()
	}
	s.mu.Unlock()

	s.wg.Wait()
	return closeErr
}
