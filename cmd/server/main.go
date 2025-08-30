package main

import "github.com/J-Obog/rapid-cache-server/src/server"

func main() {
	cfg := &server.ServerConfig{
		Address: "localhost:8076",
	}

	s := server.NewServer(cfg)
	s.Start()
}
