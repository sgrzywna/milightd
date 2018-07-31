package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Server represents milightd HTTP server.
type Server struct {
	srv *http.Server
}

// NewServer returns initialized HTTP server.
func NewServer(port int, m Controller) *Server {
	var s Server

	r := mux.NewRouter()
	v1 := r.PathPrefix("/api/v1/").Subrouter()

	v1.HandleFunc("/light", func(w http.ResponseWriter, r *http.Request) {
		lightHandler(w, r, m)
	}).Methods("PUT")

	s.srv = &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return &s
}

// ListenAndServe starts HTTP server.
func (s *Server) ListenAndServe() error {
	return s.srv.ListenAndServe()
}

func lightHandler(w http.ResponseWriter, r *http.Request, c Controller) {
	if r.Body == nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	var l Light

	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if !c.Process(l) {
		http.Error(w, "milightd error", http.StatusInternalServerError)
		return
	}
}
