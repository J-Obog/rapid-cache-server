package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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
	r.HandleFunc("/set", s.handleKeySet).Methods(http.MethodPost)
	r.HandleFunc("/del", s.handleKeyDelete).Methods(http.MethodPost)
	http.Handle("/", r)

	log.Println("Starting up server")
	log.Fatal(http.ListenAndServe(s.cfg.Address, nil))
}

func (s *Server) handleKeySet(w http.ResponseWriter, r *http.Request) {
	timeNow := time.Now()

	var setRequest SetKeyRequest

	bodyDecoder := json.NewDecoder(r.Body)
	bodyDecoder.Decode(&setRequest) //TODO: Handle error

	commandToAppend := filesystem.Command{
		Name:      filesystem.CommandNameSet,
		Key:       setRequest.Key,
		Timestamp: timeNow,
		Seed:      "", //TODO: Generate seed
		Params: map[filesystem.CommandParamKey]string{
			filesystem.CommandParamKeyValue:     setRequest.Value,
			filesystem.CommandParamKeyExpiresAt: strconv.FormatUint(setRequest.ExpiresAtMillis, 10),
		},
	}

	s.doFileWrite(&commandToAppend)
	w.WriteHeader(http.StatusAccepted)
}

func (s *Server) handleKeyDelete(w http.ResponseWriter, r *http.Request) {
	timeNow := time.Now()

	var deleteRequest DeleteKeyRequest

	bodyDecoder := json.NewDecoder(r.Body)
	bodyDecoder.Decode(&deleteRequest) //TODO: Handle error

	commandToAppend := filesystem.Command{
		Name:      filesystem.CommandNameDelete,
		Key:       deleteRequest.Key,
		Timestamp: timeNow,
		Seed:      "", //TODO: Generate seed
	}

	s.doFileWrite(&commandToAppend)
	w.WriteHeader(http.StatusAccepted)
}

func (s *Server) doFileWrite(newCommand *filesystem.Command) {
	if s.cfg.SaveToFileSynchronously {
		s.aof.Append(newCommand)
		return
	}

	go s.aof.Append(newCommand)
}
