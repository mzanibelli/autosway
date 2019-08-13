package autosway

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

var testSetup Setup = Setup{[]Output{{
	Name:   "LVDS-1",
	Make:   "Apple Computer Inc",
	Model:  "Color LCD",
	Serial: "0x00000000",
}}}

func TestRepository(t *testing.T) {
	tmpRoot, err := ioutil.TempDir("./testdata", "test-repository-")
	if err != nil {
		t.Fatal("could not initialize test environment:", err)
	}
	defer os.RemoveAll(tmpRoot)

	SUT := NewRepository(tmpRoot)
	if err := SUT.Save(&testSetup, "foo"); err != nil {
		t.Error(err)
	}

	var readSetup Setup
	err = SUT.Load(&readSetup, "foo")
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(testSetup, readSetup) {
		t.Error("serialization inconsistency")
	}
}
