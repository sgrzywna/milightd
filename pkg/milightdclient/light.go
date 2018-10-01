package milightdclient

import "github.com/sgrzywna/milightd/internal/app/milightd"

// Light represents light control command.
type Light struct {
	milightd.Light
	// actual values
	color      string
	brightness int
	state      string
}

// SetColor sets color name.
func (l *Light) SetColor(color string) {
	l.color = color
	l.Color = &l.color
}

// SetBrightness sets light brightness.
func (l *Light) SetBrightness(brightness int) {
	l.brightness = brightness
	l.Brightness = &l.brightness
}

// SetSwitch sets light state.
func (l *Light) SetSwitch(state bool) {
	if state {
		l.state = milightd.On
	} else {
		l.state = milightd.Off
	}
	l.Switch = &l.state
}

// Clear sets all attributes to their zero values.
func (l *Light) Clear() {
	l.Color = nil
	l.Brightness = nil
	l.Switch = nil
	l.color = ""
	l.brightness = 0
	l.state = ""
}

// Assign assign milightd.Light structure to the Light structure.
func (l *Light) Assign(light milightd.Light) {
	l.Clear()
	if light.Color != nil {
		l.SetColor(*light.Color)
	}
	if light.Brightness != nil {
		l.SetBrightness(*light.Brightness)
	}
	if light.Switch != nil {
		l.SetSwitch(*light.Switch == milightd.On)
	}
}
