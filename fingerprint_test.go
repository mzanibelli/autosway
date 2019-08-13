package autosway

import "testing"

var (
	screen1 Setup = Setup{[]Output{{
		Make:   "Apple Computer Inc",
		Model:  "Color LCD",
		Serial: "0x00000000",
	}}}
	screen2 Setup = Setup{[]Output{{
		Make:   "Samsung",
		Model:  "Color LCD",
		Serial: "0x00000000",
	}}}
	screen3 Setup = Setup{[]Output{{
		Make:   "Apple Computer Inc",
		Model:  "Apple TV",
		Serial: "0x00000000",
	}}}
	screen4 Setup = Setup{[]Output{{
		Make:   "Apple Computer Inc",
		Model:  "Color LCD",
		Serial: "0x00000001",
	}}}
	screen11 Setup = Setup{[]Output{
		{
			Make:   "Apple Computer Inc",
			Model:  "Color LCD",
			Serial: "0x00000000",
		},
		{
			Make:   "Apple Computer Inc",
			Model:  "Color LCD",
			Serial: "0x00000000",
		}}}
	screen12 Setup = Setup{[]Output{
		{
			Make:   "Apple Computer Inc",
			Model:  "Color LCD",
			Serial: "0x00000000",
		},
		{
			Make:   "Samsung",
			Model:  "Color LCD",
			Serial: "0x00000000",
		}}}
	screen21 Setup = Setup{[]Output{
		{
			Make:   "Samsung",
			Model:  "Color LCD",
			Serial: "0x00000000",
		},
		{
			Make:   "Apple Computer Inc",
			Model:  "Color LCD",
			Serial: "0x00000000",
		}}}
)

func TestFingerprint(t *testing.T) {
	tests := []struct {
		name   string
		setup1 Setup
		setup2 Setup
		match  bool
	}{
		{
			name:   "it should match when the setup is the same",
			setup1: screen1,
			setup2: screen1,
			match:  true,
		},
		{
			name:   "it should not match when vendors are different",
			setup1: screen1,
			setup2: screen2,
			match:  false,
		},
		{
			name:   "it should not match when models are different",
			setup1: screen1,
			setup2: screen3,
			match:  false,
		},
		{
			name:   "it should not match when serials are different",
			setup1: screen1,
			setup2: screen4,
			match:  false,
		},
		{
			name:   "it should not match when the number of monitor is different",
			setup1: screen1,
			setup2: screen12,
			match:  false,
		},
		{
			name:   "it should not match when one monitor is different in a multi-monitor setup",
			setup1: screen11,
			setup2: screen12,
			match:  false,
		},
		{
			name:   "it should match regardless to the order of the monitors",
			setup1: screen12,
			setup2: screen21,
			match:  true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f1 := Fingerprint(test.setup1)
			f2 := Fingerprint(test.setup2)
			if !test.match && f1 == f2 {
				t.Error("fingerprints should different")
			}
			if test.match && f1 != f2 {
				t.Error("fingerprints should be the same")
			}
		})
	}
}
