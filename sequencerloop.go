package main

import (
	"time"
)

// SequencerLoop represents sequencer loop.
type SequencerLoop struct {
	lightCtrl LightAPI
	seq       *Sequence
	stop      chan struct{}
	finished  chan bool
	step      int
}

// NewSequencerLoop returns initialized SequencerLoop object.
func NewSequencerLoop(lightCtrl LightAPI, seq *Sequence) *SequencerLoop {
	loop := SequencerLoop{
		lightCtrl: lightCtrl,
		seq:       seq,
		stop:      make(chan struct{}),
		finished:  make(chan bool),
	}
	go loop.loop()
	return &loop
}

// Stop terminates sequencer loop.
func (l *SequencerLoop) Stop() {
	close(l.stop)
	<-l.finished
}

// loop is the sequencer main loop.
func (l *SequencerLoop) loop() {
	defer func() { l.finished <- true }()

	delay := time.Millisecond

	for {
		select {
		case <-l.stop:
			return
		case <-time.After(delay):
			delay = l.processStep()
		}
	}
}

// processStep executes next step from sequence.
func (l *SequencerLoop) processStep() time.Duration {
	defer func() { l.step++ }()
	if l.step >= len(l.seq.Steps) {
		l.step = 0
	}
	l.lightCtrl.Process(l.seq.Steps[l.step].Light)
	return time.Duration(l.seq.Steps[l.step].Duration) * time.Millisecond
}
