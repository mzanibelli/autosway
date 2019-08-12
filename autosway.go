package autosway

import (
	"encoding/json"
	"errors"
	"fmt"
)

func Current(ipc *IPC, setup *Setup) error {
	_, res, err := ipc.Roundtrip(GET_OUTPUTS)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(res, setup); err != nil {
		return err
	}
	return nil
}

func Auto(ipc *IPC, repo *Repository, setup *Setup, f string) error {
	if err := repo.Load(setup, f); err != nil {
		return err
	}
	if Fingerprint(*setup) != f {
		return errors.New("corrupted profile")
	}
	for _, c := range setup.Commands() {
		if err := RunCmd(ipc, c); err != nil {
			return err
		}
	}
	return nil
}

func RunCmd(ipc *IPC, c string) error {
	_, res, err := ipc.Roundtrip(RUN_COMMAND, []byte(c)...)
	if err != nil {
		return err
	}
	var resp []struct{ Success bool }
	if err := json.Unmarshal(res, &resp); err != nil {
		return err
	}
	for _, r := range resp {
		if !r.Success {
			return fmt.Errorf("error running command: %v", r)
		}
	}
	return nil
}

func Save(repo *Repository, setup *Setup, f string) error {
	if err := repo.Save(setup, f); err != nil {
		return err
	}
	return nil
}
