package autosway

import (
	"bufio"
	"fmt"
	"io"
)

const (
	MAGIC_STRING = "i3-ipc"
	BUFFER_SIZE  = 1024 * 16 // 16Ko
)

type IPC struct {
	socket *bufio.ReadWriter
}

func NewIPC(socket io.ReadWriter) *IPC {
	return &IPC{socket: bufio.NewReadWriter(
		bufio.NewReaderSize(socket, BUFFER_SIZE),
		bufio.NewWriterSize(socket, BUFFER_SIZE),
	)}
}

func (ipc *IPC) Roundtrip(t int32, p ...byte) (int32, []byte, error) {
	request := NewMessage(MAGIC_STRING, t, p...)
	if err := ipc.sendRequest(request); err != nil {
		return 0, nil, err
	}
	response, err := ipc.parseReply()
	if err != nil {
		return 0, nil, err
	}
	return response.Type, response.Payload, nil
}

func (ipc *IPC) sendRequest(m *Message) error {
	if err := m.Serialize(); err != nil {
		return fmt.Errorf("request: serialize: %v", err)
	}
	if _, err := io.Copy(ipc.socket, m); err != nil {
		return fmt.Errorf("socket: write: %v", err)
	}
	if err := ipc.socket.Flush(); err != nil {
		return err
	}
	return nil
}

func (ipc *IPC) parseReply() (*Message, error) {
	m := NewMessageSize(MAGIC_STRING, BUFFER_SIZE)
	for n := 0; n == 0 || ipc.shouldHandleNextBytes(); n++ {
		if _, err := io.CopyN(m, ipc.socket, 1); err != nil {
			return nil, fmt.Errorf("socket: read: %v", err)
		}
	}
	if err := m.Unserialize(); err != nil {
		return nil, fmt.Errorf("response: unserialize: %v", err)
	}
	return m, nil
}

func (ipc *IPC) shouldHandleNextBytes() bool {
	offset, remaining := len(MAGIC_STRING), ipc.Buffered()
	switch {
	case remaining == 0:
		return false
	case remaining < offset:
		return true
	}
	next, err := ipc.socket.Peek(offset)
	if err != nil || isMagic(next) {
		return false
	}
	return true
}

func (ipc *IPC) Buffered() int {
	return ipc.socket.Reader.Buffered()
}

func isMagic(b []byte) bool {
	return string(b) == MAGIC_STRING
}
