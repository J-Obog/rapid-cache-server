package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/J-Obog/rapid-cache-server/src/cachemap"
	"github.com/J-Obog/rapid-cache-server/src/filesystem"
	"github.com/gorilla/mux"
)

type Server struct {
	cacheMap *cachemap.CacheMap
	aof      *filesystem.AppendOnlyCommandFile
	cfg      ServerConfig
}

func NewServer(cfg *ServerConfig) *Server {
	server := &Server{
		cacheMap: &cachemap.CacheMap{},
		aof:      &filesystem.AppendOnlyCommandFile{},
		cfg:      *cfg,
	}

	return server
}

func (s *Server) Start() {
	r := mux.NewRouter()
	r.HandleFunc("/get", s.handleKeySet).Methods(http.MethodPost)
	r.HandleFunc("/set", s.handleKeySet).Methods(http.MethodPost)
	r.HandleFunc("/del", s.handleKeyDelete).Methods(http.MethodPost)
	http.Handle("/", r)

	log.Println("Starting up server")
	log.Fatal(http.ListenAndServe(s.cfg.Address, nil))
}

func (s *Server) handleKeySet(w http.ResponseWriter, r *http.Request) {
	timeNow := time.Now()

	var request SetKeyRequest

	bodyDecoder := json.NewDecoder(r.Body)
	bodyDecoder.Decode(&request) //TODO: Handle error

	expiresAtAsTimeObj := time.UnixMilli(int64(request.ExpiresAtMillis))

	s.cacheMap.Set(request.Key, request.Value, expiresAtAsTimeObj, timeNow)

	change := filesystem.KeyUpdate{
		Key:       request.Key,
		Val:       request.Value,
		ExpiresAt: expiresAtAsTimeObj,
		Timestamp: timeNow,
		Seed:      "", //TODO: Generate seed
	}

	s.doFileWrite(&change)
	w.WriteHeader(http.StatusAccepted)
}

func (s *Server) handleKeyDelete(w http.ResponseWriter, r *http.Request) {
	timeNow := time.Now()

	var request DeleteKeyRequest

	bodyDecoder := json.NewDecoder(r.Body)
	bodyDecoder.Decode(&request) //TODO: Handle error

	change := filesystem.KeyDelete{
		Key:       request.Key,
		Timestamp: timeNow,
		Seed:      "", //TODO: Generate seed
	}

	s.cacheMap.Delete(request.Key, timeNow)

	s.doFileWrite(&change)
	w.WriteHeader(http.StatusAccepted)
}

func (s *Server) handleKeyGet(w http.ResponseWriter, r *http.Request) {
	timeNow := time.Now()

	var request GetKeyRequest

	bodyDecoder := json.NewDecoder(r.Body)
	bodyDecoder.Decode(&request) //TODO: Handle error

	s.cacheMap.Get(request.Key, timeNow)

	w.WriteHeader(http.StatusAccepted) //TODO: return actual data
}

func (s *Server) doFileWrite(StateChange *filesystem.StateChange) {
	if s.cfg.SaveToFileSynchronously {
		s.aof.Append(StateChange)
		return
	}

	go s.aof.Append(StateChange)
}
