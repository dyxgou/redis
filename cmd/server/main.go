package main

import (
	"github/dyxgou/redis/pkg/config"
	"github/dyxgou/redis/pkg/server"
	"log"
)

func main() {
	port := config.GetEnv("PORT")

	server := server.NewServer(server.Config{ListenAddr: port})

	log.Fatal(server.Start())
}
