package main

import (
	"bytes"
	"fmt"
	"testing"
)

var testPayload = []byte{
	byte(0x69),
	byte(0x33),
	byte(0x2d),
	byte(0x69),
	byte(0x70),
	byte(0x63),
	byte(0x00),
	byte(0x00),
	byte(0x00),
	byte(0x00),
	byte(0x03),
	byte(0x00),
	byte(0x00),
	byte(0x00),
}

func TestEncoding(t *testing.T) {
	SUT, err := buildMessage(GET_OUTPUTS)
	if err != nil {
		t.Error(err)
	}
	expected := dump(testPayload)
	actual := dump(SUT.buffer.Bytes())
	if expected != actual {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}
}

func TestDecoding(t *testing.T) {
	input := bytes.NewBuffer(testPayload)
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

func dump(in []byte) (out string) {
	for _, c := range in {
		out = fmt.Sprintf("%s0x%02x", out, c)
	}
	return
}
