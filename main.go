package main

import (
	"flag"
	"log"
)

func main() {
	var mihost = flag.String("mihost", "", "Mi-Light network address")
	var miport = flag.Int("miport", 5987, "Mi-Light network port")
	var port = flag.Int("port", 8080, "listening port")

	flag.Parse()

	m := NewMilightController(*mihost, *miport)
	defer m.Close()

	srv := NewServer(*port, m)

	log.Printf("milightd listening @ :%d\n", *port)
	log.Fatal(srv.ListenAndServe())
}
