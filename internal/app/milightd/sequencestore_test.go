package milightd

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

var (
	n0 = "first"
	c0 = "yellow"
	b0 = 1
	s0 = "on"

	n1 = "second"
	c1 = "green"
	b1 = 2
	s1 = "off"

	tests = []Sequence{
		Sequence{
			Name: n0,
			Steps: []SequenceStep{
				SequenceStep{
					Light: Light{
						Color:      &c0,
						Brightness: &b0,
						Switch:     &s0,
					},
					Duration: 100,
				},
				SequenceStep{
					Light: Light{
						Color:      &c1,
						Brightness: &b1,
						Switch:     &s1,
					},
					Duration: 200,
				},
			},
		},
		Sequence{
			Name: n1,
			Steps: []SequenceStep{
				SequenceStep{
					Light: Light{
						Color:      &c1,
						Brightness: &b1,
						Switch:     &s1,
					},
					Duration: 300,
				},
				SequenceStep{
					Light: Light{
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

func TestSequenceStoreAddGet(t *testing.T) {
	store, dirRemove := testTempStore(t)
	defer dirRemove()

	for _, tc := range tests {
		seq, err := store.Get(tc.Name)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(tc, *seq) {
			t.Errorf("expected: %v, got: %v", tc, seq)
		}
	}
}

func TestSequenceStoreGetAll(t *testing.T) {
	store, dirRemove := testTempStore(t)
	defer dirRemove()

	sequences, err := store.GetAll()
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(tests, sequences) {
		t.Errorf("expected: %v, got: %v", tests, sequences)
	}
}

func TestSequenceStoreRemove(t *testing.T) {
	store, dirRemove := testTempStore(t)
	defer dirRemove()

	err := store.Remove(n0)
	if err != nil {
		t.Error(err)
	}

	_, err = store.Get(n0)
	if err == nil {
		t.Errorf("expected error")
	}
}

func testTempStore(t *testing.T) (*SequenceStore, func()) {
	dir, dirRemove := testTempDir(t)

	store, err := NewSequenceStore(dir)
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range tests {
		err = store.Add(tc)
		if err != nil {
			t.Error(err)
		}
	}

	return store, dirRemove
}

func testTempDir(t *testing.T) (string, func()) {
	dir, err := ioutil.TempDir("", "tmp")
	if err != nil {
		t.Fatal(err)
	}
	return dir, func() { defer os.RemoveAll(dir) }
}
