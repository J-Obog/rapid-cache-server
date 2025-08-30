package server

import "github.com/gorilla/mux"

type Server struct {
	r mux.Router
}

func NewServer() *Server {
	return &Server{
		r: *mux.NewRouter(),
	}
}
