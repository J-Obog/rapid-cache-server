package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/J-Obog/rapid-cache-server/src/command"
	"github.com/gorilla/mux"
)

type Server struct {
	r   *mux.Router
	aol command.AppendOnlyCommandList
}

func NewServer() *Server {
	server := &Server{
		r:   mux.NewRouter(),
		aol: *command.NewAppendOnlyCommandList(),
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
	timeNow := time.Now()

	data := make(map[string]interface{})

	bodyDecoder := json.NewDecoder(r.Body)
	bodyDecoder.Decode(data) //TODO: Handle error

	key := data["k"].(string)
	val := data["v"].(string)
	exp := data["expMs"].(int64)

	s.aol.Append(command.Command{
		Name:      command.CommandNameSet,
		Key:       key,
		Timestamp: timeNow,
		Seed:      "", //TODO: Generate seed
		Params: map[string]string{
			"value":     val,
			"expiresAt": strconv.FormatInt(exp, 10),
		},
	})

	w.WriteHeader(http.StatusAccepted)
}

func (s *Server) handleKeyDelete(w http.ResponseWriter, r *http.Request) {
	timeNow := time.Now()

	data := make(map[string]interface{})

	bodyDecoder := json.NewDecoder(r.Body)
	bodyDecoder.Decode(data) //TODO: Handle error

	key := data["k"].(string)

	s.aol.Append(command.Command{
		Name:      command.CommandNameDelete,
		Key:       key,
		Timestamp: timeNow,
		Seed:      "", //TODO: Generate seed
		Params:    map[string]string{},
	})

	w.WriteHeader(http.StatusAccepted)
}

func (s *Server) handleDiff(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	bodyDecoder := json.NewDecoder(r.Body)
	bodyDecoder.Decode(data) //TODO: Handle error

	var l []command.Command

	if afterMsIface, exists := data["afterMs"]; exists {
		afterMs, _ := afterMsIface.(int64)
		l = s.aol.GetAllAfterTimestamp(time.UnixMilli(afterMs))
	} else {
		l = s.aol.GetAll()
	}

	bodyEncoder := json.NewEncoder(w)
	bodyEncoder.Encode(l)
}
