package main

import (
	"reflect"
	"testing"
	"time"
)

type LightAPIRecorder struct {
	calls []Light
}

func (r *LightAPIRecorder) Process(fromSequence bool, l Light) bool {
	r.calls = append(r.calls, l)
	return true
}

func TestSequencerLoop(t *testing.T) {
	var (
		n0 = "first"

		c0 = "yellow"
		b0 = 1
		s0 = "on"

		c1 = "green"
		b1 = 2
		s1 = "off"
	)

	seq := Sequence{
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
	}

	rec := LightAPIRecorder{}

	loop := NewSequencerLoop(&rec, &seq)
	time.Sleep(3 * time.Second)
	loop.Stop()

	if len(rec.calls) < 3 {
		t.Fatalf("expected at least %d calls, got %d", 3, len(rec.calls))
	}

	if !reflect.DeepEqual(seq.Steps[0].Light, rec.calls[0]) {
		t.Errorf("expected %v, got %v", seq.Steps[0].Light, rec.calls[0])
	}

	if !reflect.DeepEqual(seq.Steps[1].Light, rec.calls[1]) {
		t.Errorf("expected %v, got %v", seq.Steps[1].Light, rec.calls[1])
	}
}
