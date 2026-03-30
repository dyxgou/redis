package server_test

import (
	"github/dyxgou/redis/internal/server"
	"github/dyxgou/redis/pkg/serializer"
	"log/slog"
	"net"
	"os"
	"testing"
)

type testSuite struct {
	server  *server.Server
	command string
}

var ts *testSuite

func serialize(text string) (string, error) {
	s := serializer.New(text)

	serialized, err := s.Serialize()

	if err != nil {
		return "", err
	}

	return serialized, nil
}

func New(addr, cmd string) *testSuite {
	return &testSuite{
		server:  server.New(server.Config{ListenAddr: addr}),
		command: cmd,
	}
}

func TestMain(m *testing.M) {
	ts = New(":5000", "GET username")

	go func() {
		err := ts.server.Start()

		if err != nil {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	code := m.Run()

	os.Exit(code)
}

func TestGetNilCommand(t *testing.T) {
	slog.Info("listen addr from test get", "addr", ts.server.Config.ListenAddr)
	conn, err := net.Dial("tcp", "localhost"+ts.server.Config.ListenAddr)
	if err != nil {
		t.Fatal("connection failed", "err", err)
	}

	t.Cleanup(func() {
		if err := conn.Close(); err != nil {
			t.Fatal(err)
		}
	})

	serialized, err := serialize(ts.command)
	if err != nil {
		t.Fatal(err)
	}

	_, err = conn.Write([]byte(serialized))
	if err != nil {
		t.Fatal(err)
	}

	resBuf := make([]byte, 1024)
	n, err := conn.Read(resBuf)
	if err != nil {
		t.Fatal(err)
	}

	msg := string(resBuf[:n])
	const nilRes = "(nil)"

	if msg != nilRes {
		t.Fatalf("server response expected=%s. got=%s", nilRes, msg)
	}
}

func TestSetCommand(t *testing.T) {
	slog.Info("listen addr from test get", "addr", ts.server.Config.ListenAddr)
	tests := []struct {
		cmd string
		res string
	}{
		{
			cmd: "SET username Alejandro",
			res: "OK",
		},
		{
			cmd: "GET username",
			res: "Alejandro",
		},
	}

	conn, err := net.Dial("tcp", "localhost"+ts.server.Config.ListenAddr)
	if err != nil {
		t.Fatal("connection failed", "err", err)
	}

	t.Cleanup(func() {
		if err := conn.Close(); err != nil {
			t.Fatal(err)
		}
	})

	for _, tt := range tests {
		serialized, err := serialize(tt.cmd)
		if err != nil {
			t.Fatal(err)
		}

		_, err = conn.Write([]byte(serialized))
		if err != nil {
			t.Fatal(err)
		}

		resBuf := make([]byte, 1024)
		n, err := conn.Read(resBuf)
		if err != nil {
			t.Fatal(err)
		}

		res := string(resBuf[:n])

		if res != tt.res {
			t.Fatalf("server response expected=%s. got=%s", tt.res, res)
		}
	}
}
