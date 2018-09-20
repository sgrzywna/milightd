package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

type TestController struct {
	l         Light
	sequences []Sequence
	name      string
	state     SequenceState
}

func (m *TestController) Process(fromSequence bool, l Light) bool {
	m.l = l
	return true
}

func (m *TestController) GetSequences() ([]Sequence, error) {
	return m.sequences, nil
}

func (m *TestController) GetSequence(name string) (*Sequence, error) {
	m.name = name
	return &m.sequences[0], nil
}

func (m *TestController) AddSequence(seq Sequence) error {
	m.sequences = append(m.sequences, seq)
	return nil
}

func (m *TestController) DeleteSequence(name string) error {
	m.name = name
	return nil
}

func (m *TestController) GetSequenceState() (*SequenceState, error) {
	return &m.state, nil
}

func (m *TestController) SetSequenceState(state SequenceState) (*SequenceState, error) {
	m.state = state
	return &m.state, nil
}

func TestLightHandler(t *testing.T) {
	color := "red"
	brightness := 16
	on := "on"

	data := fmt.Sprintf("{\"color\":\"%s\",\"brightness\":%d,\"switch\":\"%s\"}", color, brightness, on)

	req, err := http.NewRequest("POST", "/api/v1/light", strings.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	c := TestController{}

	rr := httptest.NewRecorder()

	newRouter(&c).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

	if c.l.Color == nil {
		t.Error("wrong color: got nil")
	} else if *c.l.Color != color {
		t.Errorf("wrong color: got %s want %s", *c.l.Color, color)
	}

	if c.l.Brightness == nil {
		t.Error("wrong brightness: got nil")
	} else if *c.l.Brightness != brightness {
		t.Errorf("wrong brightness: got %d want %d", *c.l.Brightness, brightness)
	}

	if c.l.Switch == nil {
		t.Error("wrong switch: got nil")
	} else if *c.l.Switch != on {
		t.Errorf("wrong switch: got %s want %v", *c.l.Switch, on)
	}
}

func TestGetSequences(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v1/sequence", nil)
	if err != nil {
		t.Fatal(err)
	}

	c := TestController{}
	c.sequences = tests

	rr := httptest.NewRecorder()

	newRouter(&c).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

	var sequences []Sequence

	err = json.NewDecoder(rr.Body).Decode(&sequences)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(tests, sequences) {
		t.Errorf("expected %v, got %v", tests, sequences)
	}
}

func TestGetSequence(t *testing.T) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/sequence/%s", tests[0].Name), nil)
	if err != nil {
		t.Fatal(err)
	}

	c := TestController{}
	c.sequences = tests

	rr := httptest.NewRecorder()

	newRouter(&c).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

	if c.name != tests[0].Name {
		t.Errorf("expected %s, got %s", tests[0].Name, c.name)
	}

	var seq Sequence

	err = json.NewDecoder(rr.Body).Decode(&seq)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(tests[0], seq) {
		t.Errorf("expected %v, got %v", tests[0], seq)
	}
}

func TestAddSequence(t *testing.T) {
	data, err := json.Marshal(tests[0])
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/api/v1/sequence", strings.NewReader(string(data)))
	if err != nil {
		t.Fatal(err)
	}

	c := TestController{}

	rr := httptest.NewRecorder()

	newRouter(&c).ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusCreated)
	}

	if !reflect.DeepEqual(tests[0], c.sequences[0]) {
		t.Errorf("expected %v, got %v", tests[0], c.sequences[0])
	}
}

func TestDeleteSequence(t *testing.T) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/sequence/%s", tests[0].Name), nil)
	if err != nil {
		t.Fatal(err)
	}

	c := TestController{}

	rr := httptest.NewRecorder()

	newRouter(&c).ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusNoContent)
	}

	if c.name != tests[0].Name {
		t.Errorf("expected %s, got %s", tests[0].Name, c.name)
	}
}

func TestGetSequenceState(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v1/seqctrl", nil)
	if err != nil {
		t.Fatal(err)
	}

	testState := SequenceState{
		Name:  tests[0].Name,
		State: seqRunning,
	}

	c := TestController{}
	c.state = testState

	rr := httptest.NewRecorder()

	newRouter(&c).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(testState, c.state) {
		t.Errorf("expected %v, got %v", testState, c.state)
	}
}

func TestSetSequenceState(t *testing.T) {
	testState := SequenceState{
		Name:  tests[0].Name,
		State: seqRunning,
	}

	data, err := json.Marshal(testState)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/api/v1/seqctrl", strings.NewReader(string(data)))
	if err != nil {
		t.Fatal(err)
	}

	c := TestController{}

	rr := httptest.NewRecorder()

	newRouter(&c).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(testState, c.state) {
		t.Errorf("expected %v, got %v", testState, c.state)
	}
}
