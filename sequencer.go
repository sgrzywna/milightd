package main

// Sequencer defines sequencer interface.
type Sequencer interface {
	// Start sequence.
	Start(*Sequence) error
	// Stop running sequence.
	Stop() error
	// Status returns status of the running sequence.
	Status() *Sequence
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
func (p *SequenceProcessor) Start(seq *Sequence) error {
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
func (p *SequenceProcessor) Status() *Sequence {
	if p.loop != nil {
		return p.loop.seq
	}
	return nil
}
