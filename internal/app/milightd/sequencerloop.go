package milightd

import (
	"time"

	"github.com/sgrzywna/milightd/pkg/models"
)

// SequencerLoop represents sequencer loop.
type SequencerLoop struct {
	lightCtrl LightAPI
	seq       *models.Sequence
	stop      chan struct{}
	step      int
}

// NewSequencerLoop returns initialized SequencerLoop object.
func NewSequencerLoop(lightCtrl LightAPI, seq *models.Sequence) *SequencerLoop {
	loop := SequencerLoop{
		lightCtrl: lightCtrl,
		seq:       seq,
		stop:      make(chan struct{}),
	}
	go loop.loop()
	return &loop
}

// Stop terminates sequencer loop.
func (l *SequencerLoop) Stop() {
	l.stop <- struct{}{}
	<-l.stop
}

// loop is the sequencer main loop.
func (l *SequencerLoop) loop() {
	defer func() { l.stop <- struct{}{} }()

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
	l.lightCtrl.Process(true, l.seq.Steps[l.step].Light)
	return time.Duration(l.seq.Steps[l.step].Duration) * time.Millisecond
}
