package milightd

import (
	"log"
	"time"
)

const (
	// keeperCheckPeriod defines how often keeper will try to close unused Mi-Light connection.
	keeperCheckPeriod = 3 * time.Second
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
	connman       ConnectionManagerer
	connectionTTL time.Duration
	terminate     chan struct{}
	lastActivity  time.Time
}

// NewConnectionKeeper returns initialized ConnectionManager object.
func NewConnectionKeeper(connman ConnectionManagerer, connectionTTL time.Duration) *ConnectionKeeper {
	keeper := ConnectionKeeper{
		connman:       connman,
		connectionTTL: connectionTTL,
		terminate:     make(chan struct{}),
		lastActivity:  time.Now(),
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
	k.lastActivity = time.Now()
	return k.connman.Allocate()
}

// Release releases Mi-Light connection.
func (k *ConnectionKeeper) Release() {
	k.lastActivity = time.Now()
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
		case <-time.After(keeperCheckPeriod):
			k.closeConnection()
		}
	}
}

// closeConnection closes idle connection.
func (k *ConnectionKeeper) closeConnection() {
	allocated, exists := k.connman.GetStatus()
	if exists && !allocated {
		if time.Since(k.lastActivity) > k.connectionTTL {
			k.connman.Terminate()
		}
	}
}
