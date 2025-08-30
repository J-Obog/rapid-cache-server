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
	aol command.AppendOnlyCommandList
	cfg ServerConfig
}

func NewServer(cfg *ServerConfig) *Server {
	server := &Server{
		aol: *command.NewAppendOnlyCommandList(),
		cfg: *cfg,
	}

	return server
}

func (s *Server) Start() {
	r := mux.NewRouter()
	r.HandleFunc("/set", s.handleKeySet).Methods(http.MethodPost)
	r.HandleFunc("/del", s.handleKeyDelete).Methods(http.MethodPost)
	r.HandleFunc("/diff", s.handleDiff).Methods(http.MethodPost)

	http.Handle("/", r)
	log.Println("Starting up server")
	log.Fatal(http.ListenAndServe(s.cfg.Address, nil))
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
