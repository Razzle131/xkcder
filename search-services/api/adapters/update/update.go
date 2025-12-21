package update

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"yadro.com/course/api/core"
	updatepb "yadro.com/course/proto/update"
)

type Client struct {
	log    *slog.Logger
	client updatepb.UpdateClient
	conn   *grpc.ClientConn
}

func NewClient(address string, log *slog.Logger) (*Client, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{
		client: updatepb.NewUpdateClient(conn),
		log:    log,
		conn:   conn,
	}, nil
}

func (c *Client) CloseOrLog() {
	err := c.conn.Close()
	if err != nil {
		c.log.Error("update close connection", "error", err)
	}
}

func (c Client) Ping(ctx context.Context) error {
	_, err := c.client.Ping(ctx, nil)
	return err
}

func (c Client) Status(ctx context.Context) (core.UpdateStatus, error) {
	resp, err := c.client.Status(ctx, nil)
	if err != nil {
		return core.UpdateStatus(""), err
	}

	var status core.UpdateStatus
	switch resp.Status {
	case updatepb.Status_STATUS_IDLE:
		status = core.StatusUpdateIdle
	case updatepb.Status_STATUS_RUNNING:
		status = core.StatusUpdateRunning
	default:
		status = core.StatusUpdateUnknown
	}

	return status, nil
}

func (c Client) Stats(ctx context.Context) (core.UpdateStats, error) {
	resp, err := c.client.Stats(ctx, nil)
	if err != nil {
		return core.UpdateStats{}, err
	}

	return core.UpdateStats{
		WordsTotal:    int(resp.WordsTotal),
		WordsUnique:   int(resp.WordsUnique),
		ComicsFetched: int(resp.ComicsFetched),
		ComicsTotal:   int(resp.ComicsTotal),
	}, nil
}

func (c Client) Update(ctx context.Context) error {
	_, err := c.client.Update(ctx, nil)
	if status.Code(err) == codes.AlreadyExists {
		return core.ErrAlreadyExists
	}

	return err
}

func (c Client) Drop(ctx context.Context) error {
	_, err := c.client.Drop(ctx, nil)
	return err
}
