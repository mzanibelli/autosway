package main

import "encoding/json"

type Output struct {
	Id   int
	Name string
}

type Setup struct {
	Outputs []Output
}

func (s *Setup) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Outputs)
}

func (s *Setup) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &(s.Outputs))
}
