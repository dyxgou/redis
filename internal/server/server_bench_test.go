package server

import (
	"net"
	"testing"
)

func BenchmarkSendingFromSingleConn(b *testing.B) {
	conn, err := net.Dial(tcpMethod, ts.addr)

	if err != nil {
		b.Error(err)
	}

	b.SetParallelism(1)

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_, err := conn.Write([]byte("hello form benchmark"))

			if err != nil {
				b.Error(err)
			}
		}
	})
}
