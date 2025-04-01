package main

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	mux.Router
}

func NewServer() *Server {
	s := &Server{}
	s.Router = *mux.NewRouter()
	return s
}

func (s *Server) Run(address string) error {
	return http.ListenAndServe(address, s)
}

func main(){
	s := NewServer()
	s.HandleFunc("/health", s.HealthCheckHandler).Methods("GET")
	api := s.PathPrefix("/api").Subrouter()
	users := api.PathPrefix("/users").Subrouter()
	users.HandleFunc("/create", s.CreateUser).Methods("POST")
	s.Run(":8080")
}

func (s * Server) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `${"alive": true}`)
}

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {

}