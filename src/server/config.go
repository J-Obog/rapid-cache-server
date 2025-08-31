package server

import "time"

type ServerConfig struct {
	Address                 string
	Port                    int32
	ReindexInterval         time.Duration
	OutputFilePath          string
	SaveToFileSynchronously bool
}
