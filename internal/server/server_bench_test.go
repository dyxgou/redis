package server_test

import (
	"net"
	"testing"
)

func BenchmarkGetNil(b *testing.B) {
	const cmd = "GET nilKey"

	b.ResetTimer()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {

			conn, err := net.Dial("tcp", "localhost"+ts.server.Config.ListenAddr)

			if err != nil {
				b.Fatal("connection failed", "err", err)
			}

			b.Cleanup(func() {
				if err := conn.Close(); err != nil {
					b.Fatal(err)
				}
			})

			serialized, err := serialize(cmd)
			if err != nil {
				b.Fatal(err)
			}

			_, err = conn.Write([]byte(serialized))
			if err != nil {
				b.Fatal(err)
			}

			resBuf := make([]byte, 1024)
			n, err := conn.Read(resBuf)
			if err != nil {
				b.Fatal(err)
			}

			msg := string(resBuf[:n])
			const nilRes = "(nil)"

			if msg != nilRes {
				b.Fatalf("server response expected=%s. got=%s", nilRes, msg)
			}
		}
	})
}
