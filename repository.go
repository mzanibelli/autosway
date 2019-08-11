package main

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

type Repository struct {
	root string
}

func NewRepository(root string) *Repository {
	return &Repository{root}
}

func (r *Repository) Save(s *Setup, name string) error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	path := filepath.Join(r.root, name)
	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return err
	}
	return nil
}

func (r *Repository) Load(s *Setup, name string) error {
	path := filepath.Join(r.root, name)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, s); err != nil {
		return err
	}
	return nil
}
