package milightd

import (
	"log"
	"time"
)

const (
	// keeperCheckPeriod says how often keeper will try to close unused Mi-Light connection.
	keeperCheckPeriod = 15 * time.Second
)

// ConnectionManagerer repesents Mi-Light connection manager interface.
type ConnectionManagerer interface {
	Allocate() (LightController, error)
	Release()
	IsAllocated() bool
	Terminate()
}

// ConnectionKeeper is responsible for managing network connection to the Mi-Light device.
type ConnectionKeeper struct {
	connman   ConnectionManagerer
	done      chan struct{}
	lastCheck time.Time
}

// NewConnectionKeeper returns initialized ConnectionManager object.
func NewConnectionKeeper(connman ConnectionManagerer) *ConnectionKeeper {
	keeper := ConnectionKeeper{
		connman:   connman,
		done:      make(chan struct{}),
		lastCheck: time.Now(),
	}
	go keeper.monitorLoop()
	return &keeper
}

// Terminate terminates keeper loop.
func (k *ConnectionKeeper) Terminate() {
	k.done <- struct{}{}
}

// Allocate allocates Mi-Light connection.
func (k *ConnectionKeeper) Allocate() (LightController, error) {
	return k.connman.Allocate()
}

// Release releases Mi-Light connection.
func (k *ConnectionKeeper) Release() {
	k.connman.Release()
}

// monitorLoop monitors Mi-Light connection.
func (k *ConnectionKeeper) monitorLoop() {
	log.Printf("milight monitoring loop started")
	defer log.Printf("milight monitoring loop terminated")
	for {
		select {
		case <-k.done:
			return
		case <-time.After(keeperCheckPeriod):
			k.closeConnection()
		}
	}
}

// closeConnection closes idle connection.
func (k *ConnectionKeeper) closeConnection() {
	if !k.connman.IsAllocated() {
		if time.Since(k.lastCheck) > 2*keeperCheckPeriod {
			k.connman.Terminate()
		}
	} else {
		k.lastCheck = time.Now()
	}
}
