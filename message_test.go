package autosway

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"testing"
)

var testPayloadWithoutBody = []byte{
	// magic
	byte(0x66),
	byte(0x6f),
	byte(0x6f),
	// length
	byte(0x00),
	byte(0x00),
	byte(0x00),
	byte(0x00),
	// type
	byte(0x03),
	byte(0x00),
	byte(0x00),
	byte(0x00),
}

var testPayloadWithBody = []byte{
	// magic
	byte(0x66),
	byte(0x6f),
	byte(0x6f),
	// length
	byte(0x01), // payload is not empty
	byte(0x00),
	byte(0x00),
	byte(0x00),
	// type
	byte(0x03),
	byte(0x00),
	byte(0x00),
	byte(0x00),
	// payload
	byte(0x00), // and that's true
}

var testInconsistentPayload = []byte{
	// magic
	byte(0x66),
	byte(0x6f),
	byte(0x6f),
	// length
	byte(0x01), // payload is not empty
	byte(0x00),
	byte(0x00),
	byte(0x00),
	// type
	byte(0x03),
	byte(0x00),
	byte(0x00),
	byte(0x00),
	// but that's not true
}

func TestSerializeWithoutBody(t *testing.T) {
	SUT := NewMessage("foo", GET_OUTPUTS)
	if err := SUT.Serialize(); err != nil {
		t.Error(err)
	}
	expected := dump(testPayloadWithoutBody)
	content, _ := ioutil.ReadAll(SUT)
	actual := dump(content)
	if expected != actual {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}
}

func TestSerializeWithBody(t *testing.T) {
	SUT := NewMessage("foo", GET_OUTPUTS, byte(0x00))
	if err := SUT.Serialize(); err != nil {
		t.Error(err)
	}
	expected := dump(testPayloadWithBody)
	content, _ := ioutil.ReadAll(SUT)
	actual := dump(content)
	if expected != actual {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}
}

func TestUnserializeWithoutBody(t *testing.T) {
	input := bytes.NewBuffer(testPayloadWithoutBody)
	SUT := NewMessageSize("foo", input.Len())
	if _, err := io.Copy(SUT, input); err != nil {
		t.Error(err)
	}
	if err := SUT.Unserialize(); err != nil {
		t.Error(err)
	}
	if SUT.Type != GET_OUTPUTS {
		t.Error("invalid type:", SUT.Type)
	}
	if SUT.Length != 0 {
		t.Error("invalid length:", SUT.Length)
	}
}

func TestUnserializeWithBody(t *testing.T) {
	input := bytes.NewBuffer(testPayloadWithBody)
	SUT := NewMessageSize("foo", input.Len())
	if _, err := io.Copy(SUT, input); err != nil {
		t.Error(err)
	}
	if err := SUT.Unserialize(); err != nil {
		t.Error(err)
	}
	if SUT.Type != GET_OUTPUTS {
		t.Error("invalid type:", SUT.Type)
	}
	if SUT.Length != 1 {
		t.Error("invalid length:", SUT.Length)
	}
	if !reflect.DeepEqual([]byte{byte(0x00)}, SUT.Payload) {
		t.Error("invalid payload")
	}
}

func TestConsistency(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		err   bool
	}{
		{
			name:  "it should not fail with zero length and empty body",
			input: testPayloadWithoutBody,
			err:   false,
		},
		{
			name:  "it should not fail with non-zero length and non-empty body",
			input: testPayloadWithBody,
			err:   false,
		},
		{
			name:  "it should fail if length is not zero and body is empty",
			input: testInconsistentPayload,
			err:   true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := bytes.NewBuffer(test.input)
			SUT := NewMessageSize("foo", input.Len())
			if _, err := io.Copy(SUT, input); err != nil {
				t.Error(err)
			}
			err := SUT.Unserialize()
			if test.err == true && err == nil {
				t.Error("should fail")
			}
			if !test.err && err != nil {
				t.Error(err)
			}
		})
	}
}

func dump(in []byte) (out string) {
	for _, c := range in {
		out = fmt.Sprintf("%s0x%02x", out, c)
	}
	return
}
