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

// light represents command to control light.
type Light struct {
	Color      *string `json:"color"`
	Brightness *int    `json:"brightness"`
	Switch     *string `json:"switch"`
}

// Command represents command to control Mi-Light device.
type Command interface {
	Exec(*milight.Milight) error
}

// Controller represents milight controller interface.
type Controller interface {
	Process(Light) bool
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

// Process processes light control command.
func (m *MilightController) Process(l Light) bool {
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
