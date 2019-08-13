package autosway

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	RUN_COMMAND int32 = iota
	GET_WORKSPACES
	SUBSCRIBE
	GET_OUTPUTS
	GET_TREE
	GET_MARKS
	GET_BAR_CONFIG
	GET_VERSION
	GET_BINDING_MODES
	GET_CONFIG
	SEND_TICK
	SYNC
)

type Message struct {
	Header       []byte
	Type, Length int32
	Payload      []byte
	buffer       *bytes.Buffer
}

func (m *Message) String() string {
	return fmt.Sprintf("type: %d, length: %d\n", m.Type, m.Length)
}

func NewMessage(header string, t int32, p ...byte) *Message {
	m := NewMessageSize(header, len(header)+2*4+len(p))
	m.Length = int32(len(p))
	m.Type = t
	m.Payload = p
	return m
}

func NewMessageSize(header string, size int) *Message {
	return &Message{
		Header: []byte(header),
		buffer: bytes.NewBuffer(make([]byte, 0, size)),
	}
}

func (m *Message) Len() int {
	return m.buffer.Len()
}

func (m *Message) Write(p []byte) (n int, err error) {
	return m.buffer.Write(p)
}

func (m *Message) Read(p []byte) (n int, err error) {
	return m.buffer.Read(p)
}

func (m *Message) Serialize() error {
	m.buffer.Reset()
	if err := serialize(m.buffer, m.Header); err != nil {
		return err
	}
	if err := serialize(m.buffer, m.Length); err != nil {
		return err
	}
	if err := serialize(m.buffer, m.Type); err != nil {
		return err
	}
	if err := serialize(m.buffer, m.Payload); err != nil {
		return err
	}
	return nil
}

func (m *Message) Unserialize() error {
	if err := unserialize(m.buffer, &(m.Header)); err != nil {
		return err
	}
	if err := unserialize(m.buffer, &(m.Length)); err != nil {
		return err
	}
	if err := unserialize(m.buffer, &(m.Type)); err != nil {
		return err
	}
	if int(m.Length) != m.buffer.Len() {
		return errors.New("truncated payload")
	}
	m.Payload = make([]byte, m.buffer.Len())
	copy(m.Payload, m.buffer.Bytes())
	return nil
}

func serialize(w io.Writer, from interface{}) error {
	return binary.Write(w, binary.LittleEndian, from)
}

func unserialize(r io.Reader, to interface{}) error {
	return binary.Read(r, binary.LittleEndian, to)
}
