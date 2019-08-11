package main

import (
	"encoding/json"
	"fmt"
)

const DEFAULT_TRANSFORM = "normal"

type Output struct {
	Name      string
	Make      string
	Model     string
	Serial    string
	Transform string
	Rect      Rect
}

type Rect struct {
	X, Y, Width, Height int
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

func (s Setup) Commands() []string {
	result := make([]string, 0, len(s.Outputs)*2)
	for _, o := range s.Outputs {
		result = append(result, fmt.Sprintf("output %s %s",
			o.Name, o.String()))
		result = append(result, fmt.Sprintf("output %s transform %s",
			o.Name, o.Transform))
	}
	return result
}

func (o Output) String() string {
	return fmt.Sprintf("pos %d %d res %dx%d",
		o.Rect.X, o.Rect.Y, o.Rect.Width, o.Rect.Height)
}
