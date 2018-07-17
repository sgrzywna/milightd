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

	res := true

	if l.Switch != nil {
		log.Printf("milightd light switch %s\n", *l.Switch)
		if !m.Process(&LightSwitch{on: *l.Switch}) {
			res = false
			log.Printf("milightd light switch %s failed\n", *l.Switch)
		}
	}

	if l.Brightness != nil {
		log.Printf("milightd brightness %d\n", *l.Brightness)
		if !m.Process(&LightBrightness{level: *l.Brightness}) {
			res = false
			log.Printf("milightd brightness %d failed\n", *l.Brightness)
		}
	}

	if l.Color != nil {
		log.Printf("milightd color %s\n", *l.Color)
		if !m.Process(&LightColor{color: *l.Color}) {
			res = false
			log.Printf("milightd color %s failed\n", *l.Color)
		}
	}

	if !res {
		http.Error(w, "milightd error", http.StatusInternalServerError)
		return
	}
}
