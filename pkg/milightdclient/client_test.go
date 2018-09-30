package milightdclient

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/sgrzywna/milightd/internal/app/milightd"
)

func TestSetLight(t *testing.T) {
	// empty
	l0 := Light{}
	// all set
	l1 := Light{}
	l1.SetColor("color")
	l1.SetBrightness(3)
	l1.SetSwitch(true)
	// just color
	l2 := Light{}
	l2.SetColor("color")
	// just brightness
	l3 := Light{}
	l3.SetBrightness(3)
	// just switch
	l4 := Light{}
	l4.SetSwitch(true)

	cases := []Light{l0, l1, l2, l3, l4}

	var expected Light

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var l milightd.Light
		err := json.NewDecoder(r.Body).Decode(&l)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		expected.Assign(l)
	}))
	defer server.Close()

	c := NewClient(server.URL)

	for _, tc := range cases {
		err := c.SetLight(tc)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(tc, expected) {
			t.Errorf("expected %s, got %s", tc.String(), expected.String())
		}
	}
}

var (
	n0 = "first"
	c0 = "yellow"
	b0 = 1
	s0 = "on"

	n1 = "second"
	c1 = "green"
	b1 = 2
	s1 = "off"

	tests = []milightd.Sequence{
		milightd.Sequence{
			Name: n0,
			Steps: []milightd.SequenceStep{
				milightd.SequenceStep{
					Light: milightd.Light{
						Color:      &c0,
						Brightness: &b0,
						Switch:     &s0,
					},
					Duration: 100,
				},
				milightd.SequenceStep{
					Light: milightd.Light{
						Color:      &c1,
						Brightness: &b1,
						Switch:     &s1,
					},
					Duration: 200,
				},
			},
		},
		milightd.Sequence{
			Name: n1,
			Steps: []milightd.SequenceStep{
				milightd.SequenceStep{
					Light: milightd.Light{
						Color:      &c1,
						Brightness: &b1,
						Switch:     &s1,
					},
					Duration: 300,
				},
				milightd.SequenceStep{
					Light: milightd.Light{
						Color:      &c0,
						Brightness: &b0,
						Switch:     &s0,
					},
					Duration: 400,
				},
			},
		},
	}
)

func TestGetSequences(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(tests)
		if err != nil {
			http.Error(w, "error", http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	c := NewClient(server.URL)

	ss, err := c.GetSequences()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(tests, ss) {
		t.Errorf("expected %v, got %v", tests, ss)
	}
}

func TestAddSequence(t *testing.T) {
	var expected milightd.Sequence

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := json.NewDecoder(r.Body).Decode(&expected)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	c := NewClient(server.URL)

	err := c.AddSequence(tests[0])
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(tests[0], expected) {
		t.Errorf("expected %v, got %v", expected, tests[0])
	}
}

func TestGetSequence(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(tests[0])
		if err != nil {
			http.Error(w, "error", http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	c := NewClient(server.URL)

	seq, err := c.GetSequence(n0)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(tests[0], *seq) {
		t.Errorf("expected %v, got %v", tests[0], *seq)
	}
}

func TestDeleteSequence(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	c := NewClient(server.URL)

	err := c.DeleteSequence(n0)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetSequenceState(t *testing.T) {
	testState := milightd.SequenceState{
		Name:  tests[0].Name,
		State: milightd.SeqRunning,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		err := json.NewEncoder(w).Encode(testState)
		if err != nil {
			http.Error(w, "error", http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c := NewClient(server.URL)

	seqState, err := c.GetSequenceState()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(testState, *seqState) {
		t.Errorf("expected %v, got %v", testState, *seqState)
	}
}

func TestSetSequenceState(t *testing.T) {
	var expected milightd.SequenceState

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := json.NewDecoder(r.Body).Decode(&expected)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
	}))
	defer server.Close()

	c := NewClient(server.URL)

	testState := milightd.SequenceState{
		Name:  tests[0].Name,
		State: milightd.SeqRunning,
	}

	err := c.SetSequenceState(testState)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(testState, expected) {
		t.Errorf("expected %v, got %v", expected, testState)
	}
}
