package main

import (
	"autosway"
	"flag"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	ipc := autosway.NewIPC(connectToSocket())
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

func getSwaySocketPath() string {
	path, err := exec.Command("sway", "--get-socketpath").Output()
	if err != nil {
		log.Fatal(err)
	}
	return strings.Trim(string(path), "\n")
}

func connectToSocket() net.Conn {
	conn, err := net.Dial("unix", getSwaySocketPath())
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
