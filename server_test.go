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

	req, err := http.NewRequest("POST", "/light", strings.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	c := TestController{}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lightHandler(w, r, &c)
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if *c.l.Color != color {
		t.Errorf("wrong color: got %s want %s", *c.l.Color, color)
	}

	if *c.l.Brightness != brightness {
		t.Errorf("wrong brightness: got %d want %d", *c.l.Brightness, brightness)
	}

	if *c.l.Switch != on {
		t.Errorf("wrong switch: got %s want %v", *c.l.Switch, on)
	}
}
