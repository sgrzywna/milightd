package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

const (
	defaultStoreFolderName string = "store"
)

func main() {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	defaultStoreFolder := filepath.Join(filepath.Dir(exePath), defaultStoreFolderName)

	var mihost = flag.String("mihost", "", "Mi-Light network address")
	var miport = flag.Int("miport", 5987, "Mi-Light network port")
	var port = flag.Int("port", 8080, "listening port")
	var storeDir = flag.String("store", defaultStoreFolder, "store folder")

	flag.Parse()

	m, err := NewMilightController(*mihost, *miport, *storeDir)
	if err != nil {
		log.Fatal(err)
	}
	defer m.Close()

	srv := NewServer(*port, m)

	log.Printf("milightd listening @ :%d\n", *port)
	log.Fatal(srv.ListenAndServe())
}
