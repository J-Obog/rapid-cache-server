package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/J-Obog/rapid-cache-server/src/command"
	"github.com/gorilla/mux"
)

type Server struct {
	r   *mux.Router
	aol command.AppendOnlyCommandList
	cfg ServerConfig
}

func NewServer(cfg *ServerConfig) *Server {
	server := &Server{
		r:   mux.NewRouter(),
		aol: *command.NewAppendOnlyCommandList(),
		cfg: *cfg,
	}

	server.initializeRouter()
	return server
}

func (s *Server) Start() {
	http.Handle("/", s.r)
	log.Println("Starting up server")
	log.Fatal(http.ListenAndServe(s.cfg.Address, nil))
}

func (s *Server) initializeRouter() {
	s.r.HandleFunc("/set", s.handleKeySet).Methods("POST")
	s.r.HandleFunc("/del", s.handleKeyDelete).Methods("POST")
	s.r.HandleFunc("/diff", s.handleDiff).Methods("POST")
}

func (s *Server) handleKeySet(w http.ResponseWriter, r *http.Request) {
	timeNow := time.Now()

	var setRequest *SetKeyRequest

	bodyDecoder := json.NewDecoder(r.Body)
	bodyDecoder.Decode(setRequest) //TODO: Handle error

	s.aol.Append(command.Command{
		Name:      command.CommandNameSet,
		Key:       setRequest.Key,
		Timestamp: timeNow,
		Seed:      "", //TODO: Generate seed
		Params: map[command.CommandParamKey]string{
			command.CommandParamKeyValue:     setRequest.Value,
			command.CommandParamKeyExpiresAt: strconv.FormatUint(setRequest.ExpiresAtMillis, 10),
		},
	})

	w.WriteHeader(http.StatusAccepted)
}

func (s *Server) handleKeyDelete(w http.ResponseWriter, r *http.Request) {
	timeNow := time.Now()

	var deleteRequest *DeleteKeyRequest

	bodyDecoder := json.NewDecoder(r.Body)
	bodyDecoder.Decode(deleteRequest) //TODO: Handle error

	s.aol.Append(command.Command{
		Name:      command.CommandNameDelete,
		Key:       deleteRequest.Key,
		Timestamp: timeNow,
		Seed:      "", //TODO: Generate seed
	})

	w.WriteHeader(http.StatusAccepted)
}

func (s *Server) handleDiff(w http.ResponseWriter, r *http.Request) {
	var diffRequest *DiffRequest

	bodyDecoder := json.NewDecoder(r.Body)
	bodyDecoder.Decode(diffRequest) //TODO: Handle error

	var l []command.Command

	if diffRequest.AfterTimestampMillis != nil {
		l = s.aol.GetAllAfterTimestamp(time.UnixMilli(int64(*diffRequest.AfterTimestampMillis)))
	} else {
		l = s.aol.GetAll()
	}

	bodyEncoder := json.NewEncoder(w)
	bodyEncoder.Encode(l)
}
