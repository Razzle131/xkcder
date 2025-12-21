package search

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"yadro.com/course/api/core"
	searchpb "yadro.com/course/proto/search"
)

type Client struct {
	log    *slog.Logger
	client searchpb.SearchClient
	conn   *grpc.ClientConn
}

func NewClient(address string, log *slog.Logger) (*Client, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		client: searchpb.NewSearchClient(conn),
		log:    log,
		conn:   conn,
	}, nil
}

func (c *Client) CloseOrLog() {
	err := c.conn.Close()
	if err != nil {
		c.log.Error("search close connection", "error", err)
	}
}

func (c Client) Ping(ctx context.Context) error {
	_, err := c.client.Ping(ctx, nil)
	return err
}

func (c Client) Search(ctx context.Context, phrase string, searchLimit int) ([]core.Comics, error) {
	req := searchpb.SearchRequest{
		Limit:  int64(searchLimit),
		Phrase: phrase,
	}

	resp, err := c.client.Search(ctx, &req)
	if err != nil {
		return nil, err
	}

	res := make([]core.Comics, 0, len(resp.Comics))
	for _, c := range resp.Comics {
		res = append(res, core.Comics{
			ID:  int(c.Id),
			URL: c.Url,
		})
	}

	return res, nil
}

func (c Client) SearchIndex(ctx context.Context, phrase string, searchLimit int) ([]core.Comics, error) {
	req := searchpb.SearchRequest{
		Limit:  int64(searchLimit),
		Phrase: phrase,
	}

	resp, err := c.client.SearchIndex(ctx, &req)
	if err != nil {
		return nil, err
	}

	res := make([]core.Comics, 0, len(resp.Comics))
	for _, c := range resp.Comics {
		res = append(res, core.Comics{
			ID:  int(c.Id),
			URL: c.Url,
		})
	}

	return res, nil
}
