package main

// LightSwitch represents command to switch on/off the light.
type LightSwitch struct {
	on bool
}

// Exec executes command.
func (c *LightSwitch) Exec() error {
	return nil
}
