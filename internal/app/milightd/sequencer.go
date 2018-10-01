package milightd

import "github.com/sgrzywna/milightd/pkg/models"

// Sequencer defines sequencer interface.
type Sequencer interface {
	// Start sequence.
	Start(*models.Sequence) error
	// Stop running sequence.
	Stop() error
	// Status returns status of the running sequence.
	Status() *models.Sequence
}

// SequenceProcessor implements light control sequencer.
type SequenceProcessor struct {
	lightCtrl LightAPI
	loop      *SequencerLoop
}

// NewSequenceProcessor returns initialized SequenceProcessor object.
func NewSequenceProcessor(lightCtrl LightAPI) *SequenceProcessor {
	return &SequenceProcessor{
		lightCtrl: lightCtrl,
	}
}

// Start sequence.
func (p *SequenceProcessor) Start(seq *models.Sequence) error {
	if p.loop != nil {
		p.loop.Stop()
		p.loop = nil
	}
	p.loop = NewSequencerLoop(p.lightCtrl, seq)
	return nil
}

// Stop running sequence.
func (p *SequenceProcessor) Stop() error {
	if p.loop != nil {
		p.loop.Stop()
		p.loop = nil
	}
	return nil
}

// Status returns status of the running sequence.
func (p *SequenceProcessor) Status() *models.Sequence {
	if p.loop != nil {
		return p.loop.seq
	}
	return nil
}
