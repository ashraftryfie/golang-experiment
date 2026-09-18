package udppingpong

import (
	"errors"
	"net"
	"time"
)

var ErrNotImplemented = errors.New("exercise not implemented yet")

type PingPongServer struct {
	conn net.PacketConn
}

func StartServer(addr string) (*PingPongServer, error) {
	return nil, ErrNotImplemented
}

func (s *PingPongServer) Addr() net.Addr {
	if s.conn == nil {
		return nil
	}
	return s.conn.LocalAddr()
}

func (s *PingPongServer) Close() error {
	return ErrNotImplemented
}

func SendPing(serverAddr string, id string, timeout time.Duration) (string, error) {
	return "", ErrNotImplemented
}
