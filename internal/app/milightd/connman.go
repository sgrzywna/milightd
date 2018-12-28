package milightd

import (
	"log"
	"sync"

	"github.com/sgrzywna/milight"
)

// ConnectionManager repesents Mi-Light connection manager interface.
type ConnectionManager struct {
	addr      string
	port      int
	ml        *milight.Milight
	allocated bool
	mux       sync.Mutex
}

// NewConnectionManager returns initialized ConnectionManager object.
func NewConnectionManager(addr string, port int) *ConnectionManager {
	man := ConnectionManager{
		addr: addr,
		port: port,
	}
	return &man
}

// Allocate allocates Mi-Light connection.
func (m *ConnectionManager) Allocate() (LightController, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	if m.ml != nil {
		m.allocated = true
		return m.ml, nil
	}

	ml, err := milight.NewMilight(m.addr, m.port)
	if err != nil {
		return nil, err
	}

	err = ml.InitSession()
	if err != nil {
		return nil, err
	}

	log.Printf("milight connected @ %s:%d", m.addr, m.port)

	m.ml = ml
	m.allocated = true

	return m.ml, nil
}

// GetStatus returns status of the Mi-Light connection.
func (m *ConnectionManager) GetStatus() (bool, bool) {
	m.mux.Lock()
	defer m.mux.Unlock()
	return m.allocated, m.ml != nil
}

// Release releases Mi-Light connection.
func (m *ConnectionManager) Release() {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.allocated = false
}

// Terminate terminates keeper loop.
func (m *ConnectionManager) Terminate() {
	m.mux.Lock()
	defer m.mux.Unlock()
	log.Printf("milight connection terminated")
	if m.ml != nil {
		m.ml.Close()
		m.ml = nil
	}
	m.allocated = false
}
