package autosway

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"
)

func Fingerprint(s Setup) string {
	outputs := make([]string, len(s.Outputs), len(s.Outputs))
	for i := range outputs {
		outputs[i] = slug(s.Outputs[i])
	}
	sort.Strings(outputs)
	return hash(outputs)
}

func slug(o Output) string {
	return fmt.Sprintf("%s|%s|%s", o.Make, o.Model, o.Serial)
}

func hash(outputs []string) string {
	h := sha256.New()
	h.Write([]byte(strings.Join(outputs, "+++")))
	return fmt.Sprintf("%x", h.Sum(nil))
}
