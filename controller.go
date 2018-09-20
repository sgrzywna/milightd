package main

import (
	"log"
	"time"

	"github.com/sgrzywna/milight"
)

const (
	// waitForMilightTimeout is the delay between tries to communicate with Mi-Light device.
	waitForMilightTimeout = 3 * time.Second
	// commandsBufferSize is the size of commands channel.
	commandsBufferSize = 3
	// seqRunning represents state of the running sequence.
	seqRunning = "running"
	// seqRunning represents state of the stopped sequence.
	seqStopped = "stopped"
	// seqRunning represents state of the paused sequence.
	seqPaused = "paused"
)

// LightController represents API to control the light.
type LightController interface {
	// On turns light on.
	On() error
	// Off turns light off.
	Off() error
	// Color sets light color.
	Color(color byte) error
	// White sets white light.
	White() error
	// Brightness sets brightness level.
	Brightness(brightness byte) error
}

// Command represents command to control Mi-Light device.
type Command interface {
	Exec(LightController) error
}

// Sequence represents light control sequence.
type Sequence struct {
	Name  string         `json:"name"`
	Steps []SequenceStep `json:"steps"`
}

// SequenceStep represents single step from the light control sequence.
type SequenceStep struct {
	Light    Light `json:"light"`
	Duration int   `json:"duration"`
}

// Light represents command to control light.
type Light struct {
	Color      *string `json:"color"`
	Brightness *int    `json:"brightness"`
	Switch     *string `json:"switch"`
}

// SequenceState represents sequence state.
type SequenceState struct {
	Name  string `json:"name"`
	State string `json:"state"`
}

// LightAPI represents light control interface.
type LightAPI interface {
	// Process processes light control command.
	Process(bool, Light) bool
}

// SequenceAPI represents sequence control interface.
type SequenceAPI interface {
	// GetSequences returns list of defined sequences.
	GetSequences() ([]Sequence, error)
	// GetSequence return sequence definition.
	GetSequence(string) (*Sequence, error)
	// AddSequence adds sequence.
	AddSequence(Sequence) error
	// DeleteSequence deletes sequence.
	DeleteSequence(string) error
	// GetSequenceState returns state of the running sequence.
	GetSequenceState() (*SequenceState, error)
	// SetSequenceState control state of the running sequence.
	SetSequenceState(SequenceState) (*SequenceState, error)
}

// Controller represents milight controller interface.
type Controller interface {
	LightAPI
	SequenceAPI
}

// MilightController controls Mi-Light device.
type MilightController struct {
	addr      string
	port      int
	cmds      chan Command
	sequencer Sequencer
	store     *SequenceStore
}

// NewMilightController returns initialized MilightController object.
func NewMilightController(addr string, port int, storeDir string) (*MilightController, error) {
	store, err := NewSequenceStore(storeDir)
	if err != nil {
		return nil, err
	}
	c := MilightController{
		addr:  addr,
		port:  port,
		cmds:  make(chan Command, commandsBufferSize),
		store: store,
	}
	c.sequencer = NewSequenceProcessor(&c)
	go c.loop()
	return &c, nil
}

// Close terminates controller.
func (m *MilightController) Close() {
	m.sequencer.Stop()
	close(m.cmds)
}

// Process processes light control command.
func (m *MilightController) Process(fromSequence bool, l Light) bool {
	if !fromSequence {
		m.sequencer.Stop()
	}

	res := true

	if l.Switch != nil {
		log.Printf("milightd light switch %s\n", *l.Switch)
		if !m.exec(&LightSwitch{on: *l.Switch}) {
			res = false
			log.Printf("milightd light switch %s failed\n", *l.Switch)
		}
	}

	if l.Brightness != nil {
		log.Printf("milightd brightness %d\n", *l.Brightness)
		if !m.exec(&LightBrightness{level: *l.Brightness}) {
			res = false
			log.Printf("milightd brightness %d failed\n", *l.Brightness)
		}
	}

	if l.Color != nil {
		log.Printf("milightd color %s\n", *l.Color)
		if !m.exec(&LightColor{color: *l.Color}) {
			res = false
			log.Printf("milightd color %s failed\n", *l.Color)
		}
	}

	return res
}

// exec executes command.
func (m *MilightController) exec(c Command) bool {
	select {
	case m.cmds <- c:
		return true
	default:
		return false
	}
}

// loop is the main processing loop.
func (m *MilightController) loop() {
	for {
		ok := m.innerLoop()
		if !ok {
			return
		}
	}
}

// innerLoop is the communication loop.
func (m *MilightController) innerLoop() bool {
	var ml *milight.Milight
	var err error
	for {
		// Establish connection to Mi-Light device.
		ml, err = milight.NewMilight(m.addr, m.port)
		if err != nil {
			log.Printf("milight connection problem: %s\n", err)
			time.Sleep(waitForMilightTimeout)
			continue
		}

		err = ml.InitSession()
		if err != nil {
			log.Printf("milight session problem: %s\n", err)
			time.Sleep(waitForMilightTimeout)
			continue
		}

		defer ml.Close()

		log.Printf("milight connected @ %s:%d\n", m.addr, m.port)
		defer log.Printf("milight disconnected\n")

		for {
			cmd, ok := <-m.cmds
			if !ok {
				return false
			}
			err = cmd.Exec(ml)
			if err != nil {
				if err == milight.ErrInvalidResponse {
					return true
				}
				log.Printf("milight command error: %s\n", err)
			}
		}
	}
}

// GetSequences returns list of defined sequences.
func (m *MilightController) GetSequences() ([]Sequence, error) {
	return m.store.GetAll()
}

// GetSequence return sequence definition.
func (m *MilightController) GetSequence(name string) (*Sequence, error) {
	return m.store.Get(name)
}

// AddSequence adds sequence.
func (m *MilightController) AddSequence(seq Sequence) error {
	return m.store.Add(seq)
}

// DeleteSequence deletes sequence.
func (m *MilightController) DeleteSequence(name string) error {
	return m.store.Remove(name)
}

// GetSequenceState returns state of the running sequence.
func (m *MilightController) GetSequenceState() (*SequenceState, error) {
	var state SequenceState
	sts := m.sequencer.Status()
	if sts != nil {
		state.Name = sts.Name
		state.State = seqRunning
	} else {
		state.State = seqStopped
	}
	return &state, nil
}

// SetSequenceState control state of the running sequence.
func (m *MilightController) SetSequenceState(state SequenceState) (*SequenceState, error) {
	switch state.State {
	case seqRunning:
		seq, err := m.store.Get(state.Name)
		if err != nil {
			return nil, err
		}
		m.sequencer.Start(seq)
	default:
		m.sequencer.Stop()
	}

	return m.GetSequenceState()
}
