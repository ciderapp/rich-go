//go:build windows
// +build windows

package ipc

import (
	winio "github.com/Microsoft/go-winio"
	"time"
)

// OpenSocket opens the discord-ipc-0 named pipe
func (i *Ipc) OpenSocket() error {
	// Connect to the Windows named pipe, this is a well known name
	// We use DialTimeout since it will block forever (or very very long) on Windows
	// if the pipe is not available (Discord not running)
	duration := time.Second * 2
	sock, err := winio.DialPipe(`\\.\pipe\discord-ipc-0`, &duration)
	if err != nil {
		return err
	}

	i.Socket = sock
	return nil
}
