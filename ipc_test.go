package autosway

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

var validMessage = []byte{
	byte(0x00), byte(0x00), byte(0x00), byte(0x00),
	byte(0x03), byte(0x00), byte(0x00), byte(0x00),
}

func TestRoundtrip(t *testing.T) {
	tests := []struct {
		name       string
		connection io.ReadWriter
		err        bool
	}{
		{
			name:       "should fail if the input message is invalid",
			connection: sway(""),
			err:        true,
		},
		{
			name:       "it should fail if socket write fails",
			connection: sway("writefail"),
			err:        true,
		},
		{
			name:       "it should fail if socket read fails",
			connection: sway("readfail"),
			err:        true,
		},
		{
			name:       "it should fail if response message is invalid",
			connection: sway("badresponse"),
			err:        true,
		},
		{
			name:       "should success if the message is valid",
			connection: sway("valid"),
			err:        false,
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

func sway(data string) *mockSocket {
	s := new(mockSocket)
	buf := bytes.NewBufferString("")
	switch data {
	case "writefail":
		s.writeErr = errors.New("foo")
		buf.WriteString("i3-ipc")
		buf.Write(validMessage)
		break
	case "readfail":
		s.readErr = errors.New("foo")
		buf.WriteString("i3-ipc")
		buf.Write(validMessage)
		break
	case "badresponse":
		buf.Write([]byte("foo"))
		break
	case "valid":
		buf.WriteString("i3-ipc")
		buf.Write(validMessage)
		break
	}
	s.rw = buf
	return s
}

type mockSocket struct {
	rw       *bytes.Buffer
	writeErr error
	readErr  error
}

func (r *mockSocket) Read(p []byte) (int, error) {
	if r.readErr == nil {
		return r.rw.Read(p)
	}
	return 0, r.readErr
}

func (r *mockSocket) Write(p []byte) (int, error) {
	if r.writeErr == nil {
		return len(p), nil
	}
	return 0, r.writeErr
}
