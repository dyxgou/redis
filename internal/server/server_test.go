package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"testing"
	"time"
)

var ts *TestSuite

type TestSuite struct {
	server *Server
	donech chan struct{}
	ctx    context.Context
	addr   string
}

func NewTestSuite(c Config) *TestSuite {
	s := NewServer(c)

	go s.Start()

	return &TestSuite{
		donech: make(chan struct{}),
		server: s,
		ctx:    context.Background(),
		addr:   "localhost" + c.ListenAddr,
	}
}

func (ts *TestSuite) close() {
	ts.server.quitch <- os.Interrupt
}

func TestMain(m *testing.M) {
	ts = NewTestSuite(Config{ListenAddr: ":5000"})

	code := m.Run()

	ts.close()
	os.Exit(code)
}

func (ts *TestSuite) sendMessage(msg string, conn net.Conn) error {
	now := time.Now()
	ctx, cancel := context.WithTimeout(ts.ctx, 20*time.Millisecond)
	defer cancel()

	_, err := conn.Write([]byte(msg))

	if err != nil {
		return err
	}

	go func() {
		ts.donech <- struct{}{}
	}()

	select {
	case <-ts.donech:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("conn sending message timeout. took=%s", time.Since(now).String())
	}
}

func TestSendMessage(t *testing.T) {
	conn, err := net.Dial(tcpMethod, ts.addr)
	defer conn.Close()

	if err != nil {
		t.Error(err)
	}

	err = ts.sendMessage("hello from test!", conn)

	if err != nil {
		t.Error(err)
	}

	t.Log("message sent succesfully!")
}
