package main

import (
	"time"

	"github.com/J-Obog/rapid-cache-server/src/server"
)

func main() {
	cfg := &server.ServerConfig{
		Address:      "localhost:8076",
		ReindexAfter: 5 * time.Minute,
	}

	s := server.NewServer(cfg)
	s.Start()
}
