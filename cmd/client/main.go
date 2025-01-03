package main

import (
	"bufio"
	"fmt"
	"github/dyxgou/redis/pkg/serializer"
	"log"
	"net"
	"os"
)

const prompt = "redis> "

func main() {
	conn, err := net.Dial("tcp", "localhost:5000")

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
