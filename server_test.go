package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type TestController struct {
	l Light
}

func (c *TestController) Process(l Light) bool {
	c.l = l
	return true
}

func TestLightHandler(t *testing.T) {
	color := "red"
	brightness := 16
	on := "on"

	data := fmt.Sprintf("{\"color\":\"%s\",\"brightness\":%d,\"switch\":\"%s\"}", color, brightness, on)

	req, err := http.NewRequest("POST", "/api/v1/light", strings.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	c := TestController{}

	rr := httptest.NewRecorder()

	newRouter(&c).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

	if c.l.Color == nil {
		t.Error("wrong color: got nil")
	} else if *c.l.Color != color {
		t.Errorf("wrong color: got %s want %s", *c.l.Color, color)
	}

	if c.l.Brightness == nil {
		t.Error("wrong brightness: got nil")
	} else if *c.l.Brightness != brightness {
		t.Errorf("wrong brightness: got %d want %d", *c.l.Brightness, brightness)
	}

	if c.l.Switch == nil {
		t.Error("wrong switch: got nil")
	} else if *c.l.Switch != on {
		t.Errorf("wrong switch: got %s want %v", *c.l.Switch, on)
	}
}
