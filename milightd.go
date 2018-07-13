package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// sequence: {
//   name: "alert",
//     steps: [
//       {
//         "switch": "on",
//         "color": "white"
//       },
//       "sleep:0.3",
//       {
//         "color": "green"
//       },
//       "sleep:0.3",
//       {
//         "switch": "off"
//       },
//    ]
// }

const (
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

	on  = "on"
	off = "off"
)

// light represents command to control light.
type light struct {
	Color      *string `json:"color"`
	Brightness *int    `json:"brightness"`
	Switch     *string `json:"switch"`
}

func main() {
	var mihost = flag.String("mihost", "", "Mi-Light network address")
	var miport = flag.Int("miport", 5987, "Mi-Light network port")
	var port = flag.Int("port", 8080, "listening port")

	flag.Parse()

	m := NewMilightController(*mihost, *miport)
	defer m.Close()

	r := mux.NewRouter()
	r.HandleFunc("/light", func(w http.ResponseWriter, r *http.Request) {
		lightHandler(w, r, m)
	}).Methods("PUT")

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%d", *port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("milightd listening @ :%d\n", *port)
	log.Fatal(srv.ListenAndServe())
}

func lightHandler(w http.ResponseWriter, r *http.Request, m *MilightController) {
	var l light
	if r.Body == nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if l.Switch != nil {
		log.Printf("milightd light switch %s\n", *l.Switch)
		m.Process(&LightSwitch{on: *l.Switch == on})
	}
}
