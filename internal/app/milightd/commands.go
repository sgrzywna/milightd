package milightd

import (
	"fmt"

	"github.com/sgrzywna/milight"
)

const (
	white           = "white"
	red             = "red"
	orange          = "orange"
	yellow          = "yellow"
	chartreuseGreen = "chartreusegreen"
	green           = "green"
	springGreen     = "springgreen"
	cyan            = "cyan"
	azure           = "azure"
	blue            = "blue"
	violet          = "violet"
	magenta         = "magenta"
	rose            = "rose"

	// On turns light on.
	On = "on"
	// Off turns light off.
	Off = "off"
)

// colors maps color name with corresponding color value.
var colors = map[string]byte{
	red:             milight.Red,
	orange:          milight.Orange,
	yellow:          milight.Yellow,
	chartreuseGreen: milight.ChartreuseGreen,
	green:           milight.Green,
	springGreen:     milight.SpringGreen,
	cyan:            milight.Cyan,
	azure:           milight.Azure,
	blue:            milight.Blue,
	violet:          milight.Violet,
	magenta:         milight.Magenta,
	rose:            milight.Rose,
}

// LightSwitch represents command to switch on/off the light.
type LightSwitch struct {
	on string
}

// Exec executes command.
func (c *LightSwitch) Exec(lc LightController) error {
	if c.on == On {
		return lc.On()
	}
	return lc.Off()
}

// LightBrightness represents command to control light brightness.
type LightBrightness struct {
	level int
}

// Exec executes command.
func (c *LightBrightness) Exec(lc LightController) error {
	return lc.Brightness(byte(c.level))
}

// LightColor represents command to control light color.
type LightColor struct {
	color string
}

// Exec executes command.
func (c *LightColor) Exec(lc LightController) error {
	if c.color == white {
		return lc.White()
	}
	color, ok := colors[c.color]
	if !ok {
		return fmt.Errorf("unsupported color")
	}
	return lc.Color(color)
}
