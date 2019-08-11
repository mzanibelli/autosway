package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	path := getSwaySocketPath()
	repo := NewRepository(getDatabasePath())
	conn := connectToSocket(path)
	_, res, err := NewIPC(conn).Roundtrip(GET_OUTPUTS)
	if err != nil {
		log.Fatal(err)
	}

	var setup Setup
	if err := json.Unmarshal(res, &setup); err != nil {
		log.Fatal(err)
	}

	f := Fingerprint(setup)

	switch flag.Arg(1) {
	case "auto":
		if err := repo.Load(&setup, f); err != nil {
			log.Fatal(err)
		}
		if Fingerprint(setup) != f {
			log.Fatal("corrupted profile:", f)
		}
		// TODO: set configuration
	case "save":
		if err := repo.Save(&setup, f); err != nil {
			log.Fatal(err)
		}
		break
	default:
		log.Fatal("not implemented:", flag.Arg(1))
	}
}

func getDatabasePath() string {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		configDir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	return filepath.Join(configDir, "wl-monitors")
}

func getSwaySocketPath() string {
	path, err := exec.Command("sway", "--get-socketpath").Output()
	if err != nil {
		log.Fatal(err)
	}
	return strings.Trim(string(path), "\n")
}

func connectToSocket(socket string) net.Conn {
	conn, err := net.Dial("unix", socket)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}
