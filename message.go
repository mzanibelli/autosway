package main

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
	Type, Length int32
	Payload      []byte
	buffer       *bytes.Buffer
}

func (m *Message) Read(p []byte) (n int, err error) {
	return m.buffer.Read(p)
}

func (m *Message) String() string {
	return fmt.Sprintf(
		"type: %d, length: %d\npayload: %s\n",
		m.Type, m.Length, string(m.Payload),
	)
}

func buildMessage(t int32, p ...byte) (*Message, error) {
	cap := MAGIC_LENGTH + 2*4 + len(p)
	m := &Message{
		Type:    t,
		Length:  int32(len(p)),
		Payload: p,
		buffer:  bytes.NewBuffer(make([]byte, 0, cap)),
	}
	if err := marshal(m.buffer, []byte(MAGIC_STRING)); err != nil {
		return nil, err
	}
	if err := marshal(m.buffer, m.Length); err != nil {
		return nil, err
	}
	if err := marshal(m.buffer, m.Type); err != nil {
		return nil, err
	}
	if err := marshal(m.buffer, m.Payload); err != nil {
		return nil, err
	}
	return m, nil
}

func readMessage(r *bytes.Buffer) (*Message, error) {
	m := new(Message)
	if err := unmarshal(r, &(m.Length)); err != nil {
		return nil, err
	}
	if err := unmarshal(r, &(m.Type)); err != nil {
		return nil, err
	}
	if int(m.Length) != r.Len() {
		return nil, errors.New("truncated payload")
	}
	m.Payload = make([]byte, r.Len())
	copy(m.Payload, r.Bytes())
	return m, nil
}

func marshal(w io.Writer, from interface{}) error {
	return binary.Write(w, binary.LittleEndian, from)
}

func unmarshal(r io.Reader, to interface{}) error {
	return binary.Read(r, binary.LittleEndian, to)
}
