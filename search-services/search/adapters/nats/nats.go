package nats

import (
	"log/slog"

	"github.com/nats-io/nats.go"
	"yadro.com/course/search/core"
)

type Client struct {
	log  *slog.Logger
	sub  *nats.Subscription
	conn *nats.Conn
}

func NewClient(address string, log *slog.Logger) (*Client, error) {
	nc, err := nats.Connect(address)
	if err != nil {
		return nil, err
	}

	return &Client{
		log:  log,
		sub:  nil,
		conn: nc,
	}, nil
}

func (c *Client) CloseOrLog() {
	if err := c.sub.Unsubscribe(); err != nil {
		c.log.Error("nats unsubscribe", "error", err)
	}
	c.conn.Close()
}

func (c *Client) SubscribeUpdateEvent(updateEventName string, searcher core.Searcher) error {
	sub, err := c.conn.Subscribe(updateEventName, func(msg *nats.Msg) {
		searcher.UpdateIndex()
	})
	c.sub = sub
	return err
}
