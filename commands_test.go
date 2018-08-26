package main

import (
	"testing"

	"github.com/sgrzywna/milight"
)

type TestLightController struct {
	on         bool
	off        bool
	color      byte
	white      bool
	brightness byte
}

func (lc *TestLightController) On() error {
	lc.on = true
	return nil
}

func (lc *TestLightController) Off() error {
	lc.off = true
	return nil
}

func (lc *TestLightController) Color(color byte) error {
	lc.color = color
	return nil
}

func (lc *TestLightController) White() error {
	lc.white = true
	return nil
}

func (lc *TestLightController) Brightness(brightness byte) error {
	lc.brightness = brightness
	return nil
}

func TestLightSwitchOn(t *testing.T) {
	c := LightSwitch{on}
	lc := TestLightController{}
	err := c.Exec(&lc)
	if err != nil {
		t.Error(err)
	}
	if !lc.on {
		t.Error("LightSwitchOn failed")
	}
	if lc.off {
		t.Error("Expected LightSwitchOn, got LightSwitchOff")
	}
}

func TestLightSwitchOff(t *testing.T) {
	c := LightSwitch{off}
	lc := TestLightController{}
	err := c.Exec(&lc)
	if err != nil {
		t.Error(err)
	}
	if !lc.off {
		t.Error("LightSwitchOff failed")
	}
	if lc.on {
		t.Error("Expected LightSwitchOff, got LightSwitchOn")
	}
}

func TestLightSwitchInvalid(t *testing.T) {
	c := LightSwitch{"noronoroff"}
	lc := TestLightController{}
	err := c.Exec(&lc)
	if err != nil {
		t.Error(err)
	}
	if !lc.off {
		t.Error("LightSwitchOff failed")
	}
	if lc.on {
		t.Error("Expected LightSwitchOff, got LightSwitchOn")
	}
}

func TestLightBrightness(t *testing.T) {
	b := 12
	c := LightBrightness{b}
	lc := TestLightController{}
	err := c.Exec(&lc)
	if err != nil {
		t.Error(err)
	}
	if int(lc.brightness) != b {
		t.Errorf("LightBrightness expected: %d, got: %d", b, int(lc.brightness))
	}
}

func TestLightColor(t *testing.T) {
	c := LightColor{yellow}
	lc := TestLightController{}
	err := c.Exec(&lc)
	if err != nil {
		t.Error(err)
	}
	if lc.color != milight.Yellow {
		t.Errorf("LightColor expected: %d, got: %d", milight.Yellow, lc.color)
	}
}

func TestLightColorWhite(t *testing.T) {
	c := LightColor{white}
	lc := TestLightController{}
	err := c.Exec(&lc)
	if err != nil {
		t.Error(err)
	}
	if !lc.white {
		t.Error("TestLightColorWhite failed")
	}
}

func TestLightInvalid(t *testing.T) {
	c := LightColor{"notexisting"}
	lc := TestLightController{}
	err := c.Exec(&lc)
	if err == nil {
		t.Error("LightColor expected error, got nil")
	}
}
