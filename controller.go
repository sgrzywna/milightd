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
)

// Command represents command to control Mi-Light device.
type Command interface {
	Exec() error
}

// MilightController controls Mi-Light device.
type MilightController struct {
	addr string
	port int
	cmds chan Command
}

// NewMilightController returns initialized MilightController object.
func NewMilightController(addr string, port int) *MilightController {
	c := MilightController{
		addr: addr,
		port: port,
		cmds: make(chan Command, commandsBufferSize),
	}
	go c.loop()
	return &c
}

// Close terminates controller.
func (m *MilightController) Close() {
	close(m.cmds)
}

// Process executes command.
func (m *MilightController) Process(c Command) bool {
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

// innerLoop processes commands.
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
		defer ml.Close()

		log.Printf("milight connected @ %s:%d\n", m.addr, m.port)
		defer log.Printf("milight disconnected\n")

		for {
			cmd, ok := <-m.cmds
			if !ok {
				return false
			}
			err = cmd.Exec()
			if err != nil {
				log.Printf("milight command error: %s\n", err)
				return true
			}
		}
	}
}
