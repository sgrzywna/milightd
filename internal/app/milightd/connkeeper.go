package milightd

import (
	"log"
	"time"
)

const (
	// keeperCheckPeriod defines how often keeper will try to close unused Mi-Light connection.
	keeperCheckPeriod = 15 * time.Second
)

// ConnectionManagerer repesents Mi-Light connection manager interface.
type ConnectionManagerer interface {
	Allocate() (LightController, error)
	Release()
	GetStatus() (bool, bool)
	Terminate()
}

// ConnectionKeeper is responsible for managing network connection to the Mi-Light device.
type ConnectionKeeper struct {
	connman     ConnectionManagerer
	checkPeriod time.Duration
	terminate   chan struct{}
	lastCheck   time.Time
}

// NewConnectionKeeper returns initialized ConnectionManager object.
func NewConnectionKeeper(connman ConnectionManagerer, checkPeriod time.Duration) *ConnectionKeeper {
	keeper := ConnectionKeeper{
		connman:     connman,
		checkPeriod: checkPeriod,
		terminate:   make(chan struct{}),
		lastCheck:   time.Now(),
	}
	go keeper.monitorLoop()
	return &keeper
}

// Terminate terminates keeper loop.
func (k *ConnectionKeeper) Terminate() {
	k.terminate <- struct{}{}
	<-k.terminate
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
	defer func() { k.terminate <- struct{}{} }()
	for {
		select {
		case <-k.terminate:
			k.connman.Terminate()
			return
		case <-time.After(k.checkPeriod):
			k.closeConnection()
		}
	}
}

// closeConnection closes idle connection.
func (k *ConnectionKeeper) closeConnection() {
	allocated, exists := k.connman.GetStatus()
	if exists && !allocated {
		if time.Since(k.lastCheck) > 2*k.checkPeriod {
			k.connman.Terminate()
		}
	} else {
		k.lastCheck = time.Now()
	}
}
