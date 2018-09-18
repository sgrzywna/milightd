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
	s := Server{
		srv: &http.Server{
			Handler:      newRouter(m),
			Addr:         fmt.Sprintf(":%d", port),
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		},
	}
	return &s
}

// ListenAndServe starts HTTP server.
func (s *Server) ListenAndServe() error {
	return s.srv.ListenAndServe()
}

// newRouter returns initialized HTTP router.
func newRouter(m Controller) *mux.Router {
	r := mux.NewRouter()
	v1 := r.PathPrefix("/api/v1/").Subrouter()

	v1.HandleFunc("/light", func(w http.ResponseWriter, r *http.Request) {
		lightHandler(w, r, m)
	}).Methods("POST")

	v1.HandleFunc("/sequence", func(w http.ResponseWriter, r *http.Request) {
		listSequences(w, r, m)
	}).Methods("GET")

	v1.HandleFunc("/sequence", func(w http.ResponseWriter, r *http.Request) {
		addSequence(w, r, m)
	}).Methods("POST")

	v1.HandleFunc("/sequence/{name}", func(w http.ResponseWriter, r *http.Request) {
		getSequence(w, r, m)
	}).Methods("GET")

	v1.HandleFunc("/sequence/{name}", func(w http.ResponseWriter, r *http.Request) {
		deleteSequence(w, r, m)
	}).Methods("DELETE")

	v1.HandleFunc("/seqctrl", func(w http.ResponseWriter, r *http.Request) {
		getSequenceState(w, r, m)
	}).Methods("GET")

	v1.HandleFunc("/seqctrl", func(w http.ResponseWriter, r *http.Request) {
		setSequenceState(w, r, m)
	}).Methods("POST")

	return r
}

func lightHandler(w http.ResponseWriter, r *http.Request, c Controller) {
	var l Light

	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if !c.Process(l) {
		http.Error(w, "milightd error", http.StatusInternalServerError)
		return
	}
}

func listSequences(w http.ResponseWriter, r *http.Request, c Controller) {
	sequences, err := c.GetSequences()
	if err != nil {
		http.Error(w, "milightd error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(sequences)
	if err != nil {
		http.Error(w, "milightd error", http.StatusInternalServerError)
	}
}

func addSequence(w http.ResponseWriter, r *http.Request, c Controller) {
	var seq Sequence

	err := json.NewDecoder(r.Body).Decode(&seq)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = c.AddSequence(seq)
	if err != nil {
		http.Error(w, "milightd error", http.StatusInternalServerError)
		return
	}

	newSeq, err := c.GetSequence(seq.Name)
	if err != nil {
		http.Error(w, "milightd error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(newSeq)
	if err != nil {
		http.Error(w, "milightd error", http.StatusInternalServerError)
	}
}

func getSequence(w http.ResponseWriter, r *http.Request, c Controller) {
	vars := mux.Vars(r)
	name := vars["name"]

	seq, err := c.GetSequence(name)
	if err != nil {
		http.Error(w, "sequence not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(seq)
	if err != nil {
		http.Error(w, "milightd error", http.StatusInternalServerError)
	}
}

func deleteSequence(w http.ResponseWriter, r *http.Request, c Controller) {
	vars := mux.Vars(r)
	name := vars["name"]

	err := c.DeleteSequence(name)
	if err != nil {
		http.Error(w, "sequence not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getSequenceState(w http.ResponseWriter, r *http.Request, c Controller) {
	state, err := c.GetSequenceState()
	if err != nil {
		http.Error(w, "milightd error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(state)
	if err != nil {
		http.Error(w, "milightd error", http.StatusInternalServerError)
	}
}

func setSequenceState(w http.ResponseWriter, r *http.Request, c Controller) {
	var state SequenceState

	err := json.NewDecoder(r.Body).Decode(&state)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	newState, err := c.SetSequenceState(state)
	if err != nil {
		http.Error(w, "milightd error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(newState)
	if err != nil {
		http.Error(w, "milightd error", http.StatusInternalServerError)
	}
}
