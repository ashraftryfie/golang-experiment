package udppingpong

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

var (
	ErrMalformedResponse = errors.New("malformed pong response")
)

type PingPongServer struct {
	conn net.PacketConn
	done chan struct{}
	wg   sync.WaitGroup
}

func StartServer(addr string) (*PingPongServer, error) {
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed listening on udp: %w", err)
	}

	server := &PingPongServer{
		conn: conn,
		done: make(chan struct{}),
	}

	server.wg.Add(1)
	go server.serve()

	return server, nil
}

func (s *PingPongServer) serve() {
	defer s.wg.Done()
	buf := make([]byte, 512)

	for {
		n, peerAddr, err := s.conn.ReadFrom(buf)
		if err != nil {
			select {
			case <-s.done:
				// Clean shutdown
				return
			default:
				// Network error occurred
				return
			}
		}

		payload := string(buf[:n])
		if strings.HasPrefix(payload, "PING:") {
			id := strings.TrimPrefix(payload, "PING:")
			reply := "PONG:" + id
			_, _ = s.conn.WriteTo([]byte(reply), peerAddr)
		}
	}
}

func (s *PingPongServer) Addr() net.Addr {
	return s.conn.LocalAddr()
}

func (s *PingPongServer) Close() error {
	close(s.done)
	err := s.conn.Close()
	s.wg.Wait()
	return err
}

func SendPing(serverAddr string, id string, timeout time.Duration) (string, error) {
	conn, err := net.Dial("udp", serverAddr)
	if err != nil {
		return "", fmt.Errorf("failed dialing udp: %w", err)
	}
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		return "", fmt.Errorf("failed setting deadline: %w", err)
	}

	msg := "PING:" + id
	if _, err := conn.Write([]byte(msg)); err != nil {
		return "", fmt.Errorf("failed writing udp packet: %w", err)
	}

	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil {
		return "", fmt.Errorf("failed reading udp response: %w", err)
	}

	reply := string(buf[:n])
	if !strings.HasPrefix(reply, "PONG:") {
		return "", ErrMalformedResponse
	}

	return strings.TrimPrefix(reply, "PONG:"), nil
}
