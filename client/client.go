package client

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"

	"github.com/ciderapp/rich-go/ipc"
)

type Client struct {
	Session *interface{}
	Ipc     ipc.Ipc
}

func New() *Client {
	return &Client{}
}

var logged bool

// Login sends a handshake in the socket and returns an error or nil
func (c *Client) Login(clientid string) error {
	if !logged {
		payload, err := json.Marshal(Handshake{"1", clientid})
		if err != nil {
			return err
		}

		err = c.Ipc.OpenSocket()
		if err != nil {
			return err
		}

		// TODO: Response should be parsed
		c.Ipc.Send(0, string(payload))
	}
	logged = true

	return nil
}

func (c *Client) Logout() {
	logged = false

	err := c.Ipc.CloseSocket()
	if err != nil {
		panic(err)
	}
}

func (c *Client) SetActivity(activity Activity) error {
	if !logged {
		return nil
	}

	payload, err := json.Marshal(Frame{
		"SET_ACTIVITY",
		Args{
			os.Getpid(),
			mapActivity(&activity),
		},
		c.getNonce(),
	})

	if err != nil {
		return err
	}

	// TODO: Response should be parsed
	c.Ipc.Send(1, string(payload))
	return nil
}

func (c *Client) getNonce() string {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		fmt.Println(err)
	}

	buf[6] = (buf[6] & 0x0f) | 0x40

	return fmt.Sprintf("%x-%x-%x-%x-%x", buf[0:4], buf[4:6], buf[6:8], buf[8:10], buf[10:])
}
