package milightdclient

import "github.com/sgrzywna/milightd/internal/app/milightd"

// Light represents light control command.
type Light struct {
	color      string
	brightness int
	state      string

	colorPtr      *string
	brightnessPtr *int
	statePtr      *string
}

// SetColor sets color name.
func (l *Light) SetColor(color string) {
	l.color = color
	l.colorPtr = &l.color
}

// SetBrightness sets light brightness.
func (l *Light) SetBrightness(brightness int) {
	l.brightness = brightness
	l.brightnessPtr = &l.brightness
}

// SetSwitch sets light state.
func (l *Light) SetSwitch(state bool) {
	if state {
		l.state = milightd.On
	} else {
		l.state = milightd.Off
	}
	l.statePtr = &l.state
}

// GetColor returns color name.
func (l *Light) GetColor() *string {
	return l.colorPtr
}

// GetBrightness returns light brightness.
func (l *Light) GetBrightness() *int {
	return l.brightnessPtr
}

// GetSwitch returns light state.
func (l *Light) GetSwitch() *string {
	return l.statePtr
}
