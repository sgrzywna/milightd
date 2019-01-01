package milightd

import (
	"errors"
	"log"
	"time"

	"github.com/sgrzywna/milight"
	"github.com/sgrzywna/milightd/pkg/models"
)

const (
	// waitForMilightTimeout is the delay between tries to communicate with Mi-Light device.
	waitForMilightTimeout = 3 * time.Second
	// commandsBufferSize is the size of commands channel.
	commandsBufferSize = 3
	// connectionTTL is the Mi-Light connection time to live.
	connectionTTL = 30 * time.Second
)

var (
	// errAllocateConnection is returned when there is an error with Mi-Light connection allocation.
	errAllocateConnection = errors.New("can't allocate connection")
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

// LightAPI represents light control interface.
type LightAPI interface {
	// Process processes light control command.
	Process(bool, models.Light) bool
}

// SequenceAPI represents sequence control interface.
type SequenceAPI interface {
	// GetSequences returns list of defined sequences.
	GetSequences() ([]models.Sequence, error)
	// GetSequence return sequence definition.
	GetSequence(string) (*models.Sequence, error)
	// AddSequence adds sequence.
	AddSequence(models.Sequence) error
	// DeleteSequence deletes sequence.
	DeleteSequence(string) error
	// GetSequenceState returns state of the running sequence.
	GetSequenceState() (*models.SequenceState, error)
	// SetSequenceState control state of the running sequence.
	SetSequenceState(models.SequenceState) (*models.SequenceState, error)
}

// Controller represents milight controller interface.
type Controller interface {
	LightAPI
	SequenceAPI
}

// MilightController controls Mi-Light device.
type MilightController struct {
	addr       string
	port       int
	cmds       chan Command
	sequencer  Sequencer
	store      *SequenceStore
	connkeeper *ConnectionKeeper
}

// NewMilightController returns initialized MilightController object.
func NewMilightController(addr string, port int, storeDir string) (*MilightController, error) {
	store, err := NewSequenceStore(storeDir)
	if err != nil {
		return nil, err
	}
	connman := NewConnectionManager(addr, port)
	connkeeper := NewConnectionKeeper(connman, connectionTTL)
	c := MilightController{
		addr:       addr,
		port:       port,
		cmds:       make(chan Command, commandsBufferSize),
		store:      store,
		connkeeper: connkeeper,
	}
	c.sequencer = NewSequenceProcessor(&c)
	go c.loop()
	return &c, nil
}

// Close terminates controller.
func (m *MilightController) Close() {
	m.sequencer.Stop()
	close(m.cmds)
	m.connkeeper.Terminate()
}

// Process processes light control command.
func (m *MilightController) Process(fromSequence bool, l models.Light) bool {
	if !fromSequence {
		m.sequencer.Stop()
	}

	res := true

	if l.Switch != nil {
		log.Printf("milightd light switch %s", *l.Switch)
		if !m.exec(&LightSwitch{on: *l.Switch}) {
			res = false
			log.Printf("milightd light switch %s failed", *l.Switch)
		}
	}

	if l.Brightness != nil {
		log.Printf("milightd brightness %d", *l.Brightness)
		if !m.exec(&LightBrightness{level: *l.Brightness}) {
			res = false
			log.Printf("milightd brightness %d failed", *l.Brightness)
		}
	}

	if l.Color != nil {
		log.Printf("milightd color %s", *l.Color)
		if !m.exec(&LightColor{color: *l.Color}) {
			res = false
			log.Printf("milightd color %s failed", *l.Color)
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
	log.Printf("milight controller loop started")
	defer log.Printf("milight controller loop terminated")

	for {
		cmd, ok := <-m.cmds
		if !ok {
			return
		}
		err := m.processCommand(cmd)
		if err != nil {
			if err == errAllocateConnection || err == milight.ErrInvalidResponse {
				time.Sleep(waitForMilightTimeout)
				continue
			}
			log.Printf("milight command error: %s", err)
		}
	}
}

// processCommand allocates connection to Mi-Light device and executes command.
func (m *MilightController) processCommand(cmd Command) error {
	ml, err := m.connkeeper.Allocate()
	if err != nil {
		log.Printf("can't allocate milight device: %s", err)
		return errAllocateConnection
	}
	defer m.connkeeper.Release()
	return cmd.Exec(ml)
}

// GetSequences returns list of defined sequences.
func (m *MilightController) GetSequences() ([]models.Sequence, error) {
	return m.store.GetAll()
}

// GetSequence return sequence definition.
func (m *MilightController) GetSequence(name string) (*models.Sequence, error) {
	return m.store.Get(name)
}

// AddSequence adds sequence.
func (m *MilightController) AddSequence(seq models.Sequence) error {
	return m.store.Add(seq)
}

// DeleteSequence deletes sequence.
func (m *MilightController) DeleteSequence(name string) error {
	return m.store.Remove(name)
}

// GetSequenceState returns state of the running sequence.
func (m *MilightController) GetSequenceState() (*models.SequenceState, error) {
	var state models.SequenceState
	sts := m.sequencer.Status()
	if sts != nil {
		state.Name = sts.Name
		state.State = models.SeqRunning
	} else {
		state.State = models.SeqStopped
	}
	return &state, nil
}

// SetSequenceState control state of the running sequence.
func (m *MilightController) SetSequenceState(state models.SequenceState) (*models.SequenceState, error) {
	switch state.State {
	case models.SeqRunning:
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
