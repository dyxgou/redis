package main

import (
	"github/dyxgou/redis/internal/server"
	"github/dyxgou/redis/pkg/config"
	"log"
)

func main() {
	port := config.GetEnv("PORT")

	server := server.NewServer(server.Config{ListenAddr: port})

	log.Fatal(server.Start())
}
