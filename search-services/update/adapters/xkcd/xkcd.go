package xkcd

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"yadro.com/course/update/core"
)

type Client struct {
	log    *slog.Logger
	client http.Client
	url    string
}

func NewClient(url string, timeout time.Duration, log *slog.Logger) (*Client, error) {
	if url == "" {
		return nil, fmt.Errorf("empty base url specified")
	}
	return &Client{
		client: http.Client{Timeout: timeout},
		log:    log,
		url:    url,
	}, nil
}

type xkcdGetReply struct {
	ID         int    `json:"num"`
	URL        string `json:"img"`
	Title      string `json:"title"`
	Safe_title string `json:"safe_title"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
}

func (c Client) Get(ctx context.Context, id int) (core.XKCDInfo, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%v/info.0.json", c.url, id), nil)
	if err != nil {
		c.log.Error("xkcd get new request", "error", err)
		return core.XKCDInfo{}, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		c.log.Error("xkcd get do request", "error", err)
		return core.XKCDInfo{}, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return core.XKCDInfo{}, core.ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		c.log.Error("xkcd bad get response", "status", resp.StatusCode)
		return core.XKCDInfo{}, core.ErrBadResponseStatus
	}

	reply := xkcdGetReply{}
	err = json.NewDecoder(resp.Body).Decode(&reply)
	if err != nil {
		c.log.Error("decode xkcd get json", "error", err)
		return core.XKCDInfo{}, err
	}

	return core.XKCDInfo{
		ID:          reply.ID,
		URL:         reply.URL,
		Title:       reply.Title,
		Description: strings.TrimSpace(fmt.Sprintf("%s %s %s", reply.Transcript, reply.Safe_title, reply.Alt)),
	}, nil
}

type xkcdLastIdReply struct {
	LastId int `json:"num"`
}

func (c Client) LastID(ctx context.Context) (int, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/info.0.json", c.url), nil)
	if err != nil {
		c.log.Error("xkcd last id new request", "error", err)
		return 0, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		c.log.Error("xkcd last id do request", "error", err)
		return 0, err
	}
	if resp.StatusCode != http.StatusOK {
		c.log.Error("xkcd bad last id response", "status", resp.StatusCode)
		return 0, core.ErrNotFound
	}

	reply := xkcdLastIdReply{}
	err = json.NewDecoder(resp.Body).Decode(&reply)
	if err != nil {
		c.log.Error("decode xkcd last id json", "error", err)
		return 0, err
	}

	return reply.LastId, nil
}
