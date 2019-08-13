package main

import (
	"autosway"
	"flag"
	"log"
	"net"
	"os"
	"path/filepath"
)

func main() {
	conn := connectToSocket()
	defer conn.Close()
	ipc := autosway.NewIPC(conn)
	repo := autosway.NewRepository(getDatabasePath())

	var setup autosway.Setup
	must(autosway.Current(ipc, &setup))
	f := autosway.Fingerprint(setup)
	log.Println("current:", f)

	flag.Parse()

	switch flag.Arg(0) {
	case "", "auto":
		log.Println("running configuration commands...")
		must(autosway.Auto(ipc, repo, &setup, f))
		break
	case "save":
		log.Println("saving profile to disk...")
		must(autosway.Save(repo, &setup, f))
		break
	default:
		log.Fatalf("usage: %s %s", os.Args[0], "auto|save")
	}
}

func getDatabasePath() string {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		configDir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	return filepath.Join(configDir, "autosway")
}

func connectToSocket() net.Conn {
	conn, err := net.Dial("unix", os.Getenv("SWAYSOCK"))
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
