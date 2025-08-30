package server

import "time"

type ServerConfig struct {
	Address      string
	Port         int32
	ReindexAfter time.Duration
}
