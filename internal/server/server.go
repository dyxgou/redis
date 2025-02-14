package server

import (
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

	msgch chan []byte
	ln    net.Listener

	peers     map[*Peer]struct{}
	addPeerCh chan *Peer
}

func New(cfg Config) *Server {
	return &Server{
		Config:    cfg,
		quitch:    make(chan os.Signal, 1),
		msgch:     make(chan []byte),
		peers:     make(map[*Peer]struct{}),
		addPeerCh: make(chan *Peer),
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
	for {
		select {
		case peer := <-s.addPeerCh:
			s.peers[peer] = struct{}{}
		case msg := <-s.msgch:
			slog.Info("new message", "msg", string(msg))
		case <-s.quitch:
			s.close()
		}
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
	peer := NewPeer(conn, s.msgch)

	slog.Info("connection stablished", "remoteAddr", conn.RemoteAddr())
	s.addPeerCh <- peer

	if err := peer.readConn(); err != nil {
		peer.closeConn(err)
	}
}
