package main

import (
	"bufio"
	"fmt"
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

		_, err := conn.Write([]byte(text))

		if err != nil {
			log.Fatal(err)
		}
	}
}
