package main

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestRoundtrip(t *testing.T) {
	tests := []struct {
		name       string
		connection io.ReadWriter
		err        bool
	}{
		{
			name:       "should fail if the connection is nil",
			connection: nil,
			err:        true,
		},
		{
			name:       "should fail if there is no data to read",
			connection: mockSocket(""),
			err:        true,
		},
		{
			name:       "should fail if there is less than a magic string to read",
			connection: mockSocket("i3-ip"),
			err:        true,
		},
		{
			name:       "should fail if there is just the magic string",
			connection: mockSocket("i3-ipc"),
			err:        true,
		},
		{
			name: "should fail if there is no type",
			connection: mockSocket("i3-ipc",
				byte(0x00), byte(0x00), byte(0x00), byte(0x00)),
			err: true,
		},
		{
			name: "should success if length is zero and there is a type",
			connection: mockSocket("i3-ipc",
				byte(0x00), byte(0x00), byte(0x00), byte(0x00),
				byte(0x03), byte(0x00), byte(0x00), byte(0x00)),
			err: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			SUT := NewIPC(test.connection)
			_, _, err := SUT.Roundtrip(GET_OUTPUTS)
			if test.err == true && err == nil {
				t.Error("should fail")
			}
			if !test.err && err != nil {
				t.Error(err)
			}
		})
	}
}

func TestParsing(t *testing.T) {
	args := []byte{
		byte(0x05), byte(0x00), byte(0x00), byte(0x00), // length = 5
		byte(0x03), byte(0x00), byte(0x00), byte(0x00), // type = 3
	}
	payload := []byte("hello")
	args = append(args, payload...)
	SUT := NewIPC(mockSocket("i3-ipc", args...))
	typ, res, err := SUT.Roundtrip(GET_OUTPUTS)
	if err != nil {
		t.Error(err)
	}
	if typ != GET_OUTPUTS {
		t.Error("unexpected response type:", typ)
	}
	if !reflect.DeepEqual(res, payload) {
		t.Error("unexpected response payload:", res)
	}
}

func mockSocket(data string, bs ...byte) *readOnly {
	buf := bytes.NewBufferString(data)
	for _, b := range bs {
		buf.WriteByte(b)
	}
	return &readOnly{buf}
}

type readOnly struct {
	rw *bytes.Buffer
}

func (r *readOnly) Read(p []byte) (int, error) {
	return r.rw.Read(p)
}

func (r *readOnly) Write(p []byte) (int, error) {
	return len(p), nil
}
