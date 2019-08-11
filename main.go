package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
	"strings"
)

func main() {
	path := getSwaySocketPath()
	socket := connectToSocket(path)
	_, res, err := New(socket).Roundtrip(GET_OUTPUTS)
	if err != nil {
		panic(err)
	}
	var setup Setup
	json.Unmarshal(res, &setup)
	fmt.Println(setup.Outputs[0].Name)
}

func getSwaySocketPath() string {
	path, err := exec.Command("sway", "--get-socketpath").Output()
	if err != nil {
		panic(err)
	}
	return strings.Trim(string(path), "\n")
}

func connectToSocket(socket string) net.Conn {
	conn, err := net.Dial("unix", socket)
	if err != nil {
		panic(err)
	}
	return conn
}
