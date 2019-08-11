package main

import (
	"bytes"
	"fmt"
	"testing"
)

var testPayloadWithoutBody = []byte{
	// magic
	byte(0x69),
	byte(0x33),
	byte(0x2d),
	byte(0x69),
	byte(0x70),
	byte(0x63),
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
	byte(0x69),
	byte(0x33),
	byte(0x2d),
	byte(0x69),
	byte(0x70),
	byte(0x63),
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
	byte(0x69),
	byte(0x33),
	byte(0x2d),
	byte(0x69),
	byte(0x70),
	byte(0x63),
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

func TestEncoding(t *testing.T) {
	SUT, err := buildMessage(GET_OUTPUTS)
	if err != nil {
		t.Error(err)
	}
	expected := dump(testPayloadWithoutBody)
	actual := dump(SUT.buffer.Bytes())
	if expected != actual {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}
}

func TestDecoding(t *testing.T) {
	input := bytes.NewBuffer(testPayloadWithoutBody)
	input.Next(MAGIC_LENGTH)
	SUT, err := readMessage(input)
	if err != nil {
		t.Error(err)
	}
	if SUT.Type != GET_OUTPUTS {
		t.Error("invalid type:", SUT.Type)
	}
	if SUT.Length != 0 {
		t.Error("invalid length:", SUT.Length)
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
			input.Next(MAGIC_LENGTH)
			_, err := readMessage(input)
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
