package ipc

import (
	"bytes"
	"encoding/binary"
	"net"
	"os"
)

type Ipc struct {
	Socket net.Conn
}

func New() *Ipc {
	return &Ipc{}
}

func (i *Ipc) IsOpen() bool {
	if i.Socket == nil {
		return false
	}
	buf := make([]byte, 1024)
	// i.Socket.SetReadDeadline(time.Now().Add(1 * time.Second))
	if _, err := i.Socket.Read(buf); err != nil {
		return false
	}
	return true
}

// Choose the right directory to the ipc socket and return it
func (i *Ipc) GetIpcPath() string {
	variablesnames := []string{"XDG_RUNTIME_DIR", "TMPDIR", "TMP", "TEMP"}

	for _, variablename := range variablesnames {
		path, exists := os.LookupEnv(variablename)

		if exists {
			return path
		}
	}

	return "/tmp"
}

func (i *Ipc) CloseSocket() error {
	if i.Socket != nil {
		i.Socket.Close()
		i.Socket = nil
	}
	return nil
}

// Read the socket response
func (i *Ipc) Read() (string, error) {
	buf := make([]byte, 512)
	payloadlength, err := i.Socket.Read(buf)
	if err != nil {
		//fmt.Println("Nothing to read")
	}

	buffer := new(bytes.Buffer)
	for i := 8; i < payloadlength; i++ {
		buffer.WriteByte(buf[i])
	}

	return buffer.String(), nil
}

// Send opcode and payload to the unix socket
func (i *Ipc) Send(opcode int, payload string) (string, error) {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, int32(opcode))

	if err := binary.Write(buf, binary.LittleEndian, int32(len(payload))); err != nil {
		return "", err
	}

	buf.Write([]byte(payload))
	if _, err := i.Socket.Write(buf.Bytes()); err != nil {
		return "", err
	}

	return i.Read()
}
