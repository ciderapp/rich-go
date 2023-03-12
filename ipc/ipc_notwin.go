//go:build !windows
// +build !windows

package ipc

import (
	"net"
	"time"
)

// OpenSocket opens the discord-ipc-0 unix socket
func (i *Ipc) OpenSocket() error {
	sock, err := net.DialTimeout("unix", i.GetIpcPath()+"/discord-ipc-0", time.Second*2)
	if err != nil {
		return err
	}

	i.Socket = sock
	return nil
}
