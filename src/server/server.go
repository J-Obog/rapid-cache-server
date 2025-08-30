package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	r *mux.Router
}

func NewServer() *Server {
	server := &Server{
		r: mux.NewRouter(),
	}

	server.initializeRouter()
	return server
}

func (s *Server) initializeRouter() {
	s.r.HandleFunc("/set", s.handleKeySet).Methods("POST")
	s.r.HandleFunc("/del", s.handleKeyDelete).Methods("POST")
	s.r.HandleFunc("/diff", s.handleDiff).Methods("POST")
}

func (s *Server) handleKeySet(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) handleKeyDelete(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) handleDiff(w http.ResponseWriter, r *http.Request) {

}
