package server

import (
	"log/slog"
	"net"
)

const tcpListener = "tcp"

type Config struct {
	ListenAddr string
}

type Server struct {
	Config
	ln     net.Listener
	quitch chan struct{}

	msgch chan []byte

	peers     map[*Peer]struct{}
	addPeerCh chan *Peer
}

func NewServer(cfg Config) *Server {
	return &Server{
		Config:    cfg,
		quitch:    make(chan struct{}),
		msgch:     make(chan []byte),
		peers:     make(map[*Peer]struct{}),
		addPeerCh: make(chan *Peer),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen(tcpListener, s.Config.ListenAddr)

	if err != nil {
		return err
	}

	slog.Info("server started at", "addr", s.Config.ListenAddr)
	s.ln = ln

	go s.loop()

	return s.acceptLoop()
}

func (s *Server) loop() {
	for {
		select {
		case peer := <-s.addPeerCh:
			s.peers[peer] = struct{}{}
		case msg := <-s.msgch:
			slog.Info("new message", "msg", string(msg))
		case <-s.quitch:
			return
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
		slog.Error("connection failed ", "err", err)
	}
}
