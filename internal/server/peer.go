package server

import (
	"github/dyxgou/redis/internal/evaluator"
	"github/dyxgou/redis/pkg/lexer"
	"github/dyxgou/redis/pkg/parser"
	"log/slog"
	"net"
)

type Peer struct {
	conn net.Conn

	e *evaluator.Evaluator
	p *parser.Parser
}

func NewPeer(conn net.Conn, e *evaluator.Evaluator) *Peer {
	return &Peer{
		conn: conn,
		e:    e,
	}
}

func (p *Peer) InitParser(input string) {
	if p.p == nil {
		p.p = parser.New(lexer.New(input))
	} else {
		p.p.Reset(input)
	}
}

func (p *Peer) closeConn(err error) {
	slog.Error("conn closed", "err", err, "at", p.conn.RemoteAddr())
	p.conn.Close()
}

func (p *Peer) send(s string) error {
	_, err := p.conn.Write([]byte(s))
	if err != nil {
		return err
	}

	return nil
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
		p.InitParser(string(msgBuf))

		cmd, err := p.p.Parse()
		if err != nil {
			p.send(err.Error())
			continue
		}

		res, err := p.e.Eval(cmd)
		if err != nil {
			p.send(err.Error())
			continue
		}

		p.send(res)
	}
}
