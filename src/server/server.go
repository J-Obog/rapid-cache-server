package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/J-Obog/rapid-cache-server/src/cachemap"
	"github.com/J-Obog/rapid-cache-server/src/filesystem"
	"github.com/gorilla/mux"
)

type Server struct {
	cacheMap *cachemap.CacheMap
	aof      *filesystem.WriteOperationAOF
	cfg      ServerConfig
}

func NewServer(cfg *ServerConfig) *Server {
	aof := &filesystem.WriteOperationAOF{}
	cache := cachemap.NewCacheMap()

	server := &Server{
		cacheMap: cache,
		aof:      aof,
		cfg:      *cfg,
	}

	return server
}

func (s *Server) Start() {
	if err := s.aof.Open(s.cfg.OutputFilePath); err != nil {
		log.Fatalf("Error while opening data file: %v", err)
	}

	stateChanges, _ := s.aof.Read()

	fmt.Println(stateChanges)

	for _, change := range stateChanges {
		switch v := change.(type) {
		case filesystem.SetKeyOperation:
			s.cacheMap.SetWithoutLock(v.Key, v.Val, v.ExpiresAt, v.Timestamp)
		case filesystem.DeleteKeyOperation:
			s.cacheMap.DeleteWithoutLock(v.Key, v.Timestamp)
		}
	}

	r := mux.NewRouter()
	r.HandleFunc("/get", s.handleKeySet).Methods(http.MethodPost)
	r.HandleFunc("/set", s.handleKeySet).Methods(http.MethodPost)
	r.HandleFunc("/del", s.handleKeyDelete).Methods(http.MethodPost)
	http.Handle("/", r)

	server := &http.Server{
		Addr: s.cfg.Address,
	}

	log.Println("Starting up server")

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error while running server: %v", err)
		}
		log.Println("Server finished serving requests")
	}()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 1*time.Minute)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Error while shutting down server: %v", err)
	}

	log.Println("Server shutdown")

	if err := s.aof.Close(); err != nil {
		log.Fatalf("Error while closing aof: %v", err) //TODO: Maybe dont do fatal log
	}

	log.Println("Aof has been closed")
}

func (s *Server) handleKeySet(w http.ResponseWriter, r *http.Request) {
	timeNow := time.Now()

	var request SetKeyRequest

	bodyDecoder := json.NewDecoder(r.Body)
	bodyDecoder.Decode(&request) //TODO: Handle error

	expiresAtAsTimeObj := time.UnixMilli(int64(request.ExpiresAtMillis))

	s.cacheMap.Set(request.Key, request.Value, expiresAtAsTimeObj, timeNow)

	change := filesystem.SetKeyOperation{
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

	change := filesystem.DeleteKeyOperation{
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

func (s *Server) doFileWrite(StateChange filesystem.WriteOperation) {
	if s.cfg.SaveToFileSynchronously {
		s.aof.Append(StateChange)
		return
	}

	go s.aof.Append(StateChange)
}
