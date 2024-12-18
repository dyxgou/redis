package server

import (
	"log/slog"
	"net"
)

type Peer struct {
	conn  net.Conn
	msgch chan<- []byte
}

func NewPeer(conn net.Conn, msgch chan<- []byte) *Peer {
	return &Peer{
		conn:  conn,
		msgch: msgch,
	}
}

func (p *Peer) closeConn(err error) {
	slog.Error("conn closed due to", "err", err, "at", p.conn.RemoteAddr())
	p.conn.Close()
}

func (p *Peer) readConn() error {
	buf := make([]byte, 1024)

	for {
		n, err := p.conn.Read(buf)

		if err != nil {
			return err
		}

		msgBuf := make([]byte, n)
		copy(msgBuf, buf[:n])
		p.msgch <- msgBuf
	}
}
