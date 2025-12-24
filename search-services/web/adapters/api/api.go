package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"yadro.com/course/web/core"
)

type Client struct {
	logger  *slog.Logger
	baseUrl string
	client  *http.Client
}

func New(logger *slog.Logger, address string) *Client {
	return &Client{
		logger:  logger,
		baseUrl: address,
		client:  http.DefaultClient,
	}
}

func (c *Client) Search(phrase string) (core.SearchResponse, error) {
	resp, err := http.Get(c.baseUrl + "/api/search?phrase=" + strings.ReplaceAll(phrase, " ", "_"))
	if err != nil {
		return core.SearchResponse{}, err
	}
	if resp.StatusCode != http.StatusOK {
		msg, err := io.ReadAll(resp.Body)
		if err != nil {
			return core.SearchResponse{}, err
		}
		return core.SearchResponse{}, errors.New(string(msg))
	}

	var searchResponse core.SearchResponse
	if err = json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
		return core.SearchResponse{}, err
	}

	return searchResponse, nil
}

func (c *Client) SearchRandom() (core.SearchResponse, error) {
	resp, err := http.Get(c.baseUrl + "/api/search/random")
	if err != nil {
		return core.SearchResponse{}, err
	}
	if resp.StatusCode != http.StatusOK {
		msg, err := io.ReadAll(resp.Body)
		if err != nil {
			return core.SearchResponse{}, err
		}
		return core.SearchResponse{}, errors.New(string(msg))
	}

	var searchResponse core.SearchResponse
	if err = json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
		return core.SearchResponse{}, err
	}

	return searchResponse, nil
}

func (c *Client) DbStats() (core.UpdateStatsResponse, error) {
	resp, err := http.Get(c.baseUrl + "/api/db/stats")
	if err != nil {
		return core.UpdateStatsResponse{}, err
	}
	if resp.StatusCode != http.StatusOK {
		msg, err := io.ReadAll(resp.Body)
		if err != nil {
			return core.UpdateStatsResponse{}, err
		}
		return core.UpdateStatsResponse{}, errors.New(string(msg))
	}

	var stats core.UpdateStatsResponse
	if err = json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return core.UpdateStatsResponse{}, err
	}

	return stats, nil
}

func (c *Client) DbUpdate(tokenHeader, token string) error {
	req, err := http.NewRequest("POST", c.baseUrl+"/api/db/update", nil)
	if err != nil {
		c.logger.Error(err.Error())
		return err
	}
	req.Header.Add(tokenHeader, token)

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error(err.Error())
		return err
	}
	if resp.StatusCode != http.StatusOK {
		msg, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(string(msg))
	}

	return nil
}

func (c *Client) DbDrop(tokenHeader, token string) error {
	req, err := http.NewRequest("DELETE", c.baseUrl+"/api/db", nil)
	if err != nil {
		return err
	}
	req.Header.Add(tokenHeader, token)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		msg, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(string(msg))
	}

	return nil
}

func (c *Client) Login(name, password string) (string, error) {
	req := core.LoginRequest{
		Name:     name,
		Password: password,
	}

	json, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(c.baseUrl+"/api/login", "application/json", bytes.NewBuffer(json))
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		msg, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return "", errors.New(string(msg))
	}

	token, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(token), nil
}

func (c *Client) Verify(tokenHeader, token string) error {
	req, err := http.NewRequest("POST", c.baseUrl+"/api/verify", nil)
	if err != nil {
		c.logger.Error(err.Error())
		return err
	}
	req.Header.Add(tokenHeader, token)

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error(err.Error())
		return err
	}
	if resp.StatusCode != http.StatusOK {
		msg, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(string(msg))
	}

	return nil
}
