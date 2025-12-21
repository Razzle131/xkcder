package nats

import (
	"log/slog"

	"github.com/nats-io/nats.go"
)

type Client struct {
	log  *slog.Logger
	conn *nats.Conn
}

func NewClient(address string, log *slog.Logger) (*Client, error) {
	nc, err := nats.Connect(address)
	if err != nil {
		return nil, err
	}

	return &Client{
		log:  log,
		conn: nc,
	}, nil
}

func (c *Client) CloseOrLog() {
	c.conn.Close()
}

func (c *Client) PublishUpdateEvent(updateEventName string) error {
	err := c.conn.Publish(updateEventName, []byte("XKCD DB has been updated"))
	if err != nil {
		return err
	}
	return c.conn.Flush()
}
