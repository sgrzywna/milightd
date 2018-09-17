package main

import (
	"encoding/json"

	scribble "github.com/nanobox-io/golang-scribble"
)

// SequenceStorer represents sequence store interface.
type SequenceStorer interface {
	// GetAll retrieves all sequences from store.
	GetAll() ([]Sequence, error)
	// Get retrieve single sequence from store.
	Get(string) (*Sequence, error)
	// Add stores single sequence into store.
	Add(Sequence) error
	// Remove removes single sequence from store.
	Remove(string) error
}

const (
	collection string = "sequence"
)

// SequenceStore represents sequence store.
type SequenceStore struct {
	db *scribble.Driver
}

// NewSequenceStore returns initialized NewSequenceStore object.
func NewSequenceStore(dir string) (*SequenceStore, error) {
	db, err := scribble.New(dir, nil)
	if err != nil {
		return nil, err
	}
	return &SequenceStore{db: db}, nil
}

// GetAll retrieves all sequences from store.
func (s *SequenceStore) GetAll() ([]Sequence, error) {
	records, err := s.db.ReadAll(collection)
	if err != nil {
		return nil, err
	}
	var sequences []Sequence
	for _, r := range records {
		var seq Sequence
		if err := json.Unmarshal([]byte(r), &seq); err != nil {
			return nil, err
		}
		sequences = append(sequences, seq)
	}
	return sequences, nil
}

// Get retrieve single sequence from store.
func (s *SequenceStore) Get(name string) (*Sequence, error) {
	var seq Sequence
	if err := s.db.Read(collection, name, &seq); err != nil {
		return nil, err
	}
	return &seq, nil
}

// Add stores single sequence into store.
func (s *SequenceStore) Add(seq Sequence) error {
	return s.db.Write(collection, seq.Name, seq)
}

// Remove removes single sequence from store.
func (s *SequenceStore) Remove(name string) error {
	return s.db.Delete(collection, name)
}
