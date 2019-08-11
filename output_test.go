package autosway

import (
	"reflect"
	"testing"
)

func TestCommands(t *testing.T) {
	var SUT Setup = Setup{[]Output{
		{
			Name:      "LVDS-1",
			Transform: "normal",
			Rect:      Rect{0, 0, 1440, 900},
		},
		{
			Name:      "HDMI-1",
			Transform: "90",
			Rect:      Rect{1440, 0, 800, 600},
		},
	}}
	expected := []string{
		"output LVDS-1 pos 0 0 res 1440x900",
		"output LVDS-1 transform normal",
		"output HDMI-1 pos 1440 0 res 800x600",
		"output HDMI-1 transform 90",
	}
	actual := SUT.Commands()
	if !reflect.DeepEqual(expected, actual) {
		t.Error("invalid commands generated from config:", actual)
	}
}
