package main

import (
	"bufio"
	"fmt"
	"github/dyxgou/redis/pkg/serializer"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const prompt = "redis> "

func main() {
	conn, err := net.Dial("tcp", "localhost:5000")

	quitch := make(chan os.Signal, 1)
	signal.Notify(quitch, os.Interrupt, syscall.SIGTERM)

	go func() {
		select {
		case <-quitch:
			fmt.Println()
			slog.Info("redis client", "exit", conn.RemoteAddr())
			if err := conn.Close(); err != nil {
				slog.Error("closing conn", "err", err)
			}
			os.Exit(0)
		}
	}()

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(prompt)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		text := scanner.Text()

		serialized, err := serialize(text)

		if err != nil {
			fmt.Print("\t")
			fmt.Println(err)
			continue
		}

		_, err = conn.Write([]byte(serialized))

		if err != nil {
			log.Fatal(err)
		}

		resBuf := make([]byte, 1024)
		n, err := conn.Read(resBuf)
		if err != nil {
			log.Fatal(err)
		}

		msg := string(resBuf[:n])
		fmt.Printf("server> %s\n", msg)
	}
}

func serialize(text string) (string, error) {
	s := serializer.New(text)

	serialized, err := s.Serialize()

	if err != nil {
		return "", err
	}

	return serialized, nil
}
