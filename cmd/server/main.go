package main

import (
	"github/dyxgou/redis/internal/server"
	"github/dyxgou/redis/pkg/config"
	"log"
)

func main() {
	port := config.GetEnv("PORT")

	server := server.New(server.Config{ListenAddr: port})

	log.Fatal(server.Start())
}
