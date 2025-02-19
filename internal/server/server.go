package server

import (
	"context"
	"github/dyxgou/redis/internal/evaluator"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const tcpMethod = "tcp"

type Config struct {
	ListenAddr string
}

type Server struct {
	Config
	quitch chan os.Signal

	ln net.Listener

	e *evaluator.Evaluator
}

func New(cfg Config) *Server {
	return &Server{
		Config: cfg,
		quitch: make(chan os.Signal, 1),
		e:      evaluator.New(context.Background()),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen(tcpMethod, s.Config.ListenAddr)
	defer ln.Close()

	signal.Notify(s.quitch, os.Interrupt, syscall.SIGTERM)

	if err != nil {
		return err
	}

	slog.Info("server started at", "addr", s.Config.ListenAddr, "PID", os.Getpid())
	s.ln = ln

	go s.loop()

	return s.acceptLoop()
}

func (s *Server) close() {
	slog.Info("server closed succesfully")
	os.Exit(0)
}

func (s *Server) loop() {
	select {
	case <-s.quitch:
		s.close()
	}
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()

		if err != nil {
			slog.Error("accept conn err", "err", err)
			continue
		}

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	peer := NewPeer(conn, s.e)

	slog.Info("connection stablished", "remoteAddr", conn.RemoteAddr())

	if err := peer.readConn(); err != nil {
		slog.Error("read conn", "addr", conn.RemoteAddr(), "err", err)
	}
}
