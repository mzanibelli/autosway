package autosway

import (
	"bufio"
	"bytes"
	"errors"
	"io"
)

const (
	MAGIC_STRING = "i3-ipc"
	MAGIC_LENGTH = 6
)

type IPC struct {
	rw io.ReadWriter
}

func NewIPC(rw io.ReadWriter) *IPC {
	return &IPC{rw: rw}
}

func (ipc *IPC) Roundtrip(t int32, bs ...byte) (int32, []byte, error) {
	if ipc.rw == nil {
		return 0, nil, errors.New("invalid socket")
	}
	request, err := buildMessage(t, bs...)
	if err != nil {
		return 0, nil, err
	}
	if _, err := io.Copy(ipc.rw, request); err != nil {
		return 0, nil, err
	}
	response, err := ipc.reply()
	if err != nil {
		return 0, nil, err
	}
	return response.Type, response.Payload, nil
}

func (ipc *IPC) reply() (*Message, error) {
	data := bufio.NewReader(ipc.rw)
	if _, err := data.Discard(MAGIC_LENGTH); err != nil {
		return nil, err
	}
	tmp := bytes.NewBuffer([]byte{})
	for shouldHandleNextBytes(data) {
		if _, err := io.CopyN(tmp, data, 1); err != nil {
			return nil, err
		}
	}
	m, err := readMessage(tmp)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func shouldHandleNextBytes(r *bufio.Reader) bool {
	switch {
	case r.Buffered() == 0:
		return false
	case r.Buffered() < MAGIC_LENGTH:
		return true
	}
	next, err := r.Peek(MAGIC_LENGTH)
	if err != nil || string(next) == MAGIC_STRING {
		return false
	}
	return true
}
