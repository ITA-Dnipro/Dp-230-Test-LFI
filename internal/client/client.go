package client

import (
	"io"
	"net/http"
	"time"
)

type Client struct {
	http.Client
}

func New() *Client {
	return &Client{
		http.Client{
			Timeout: time.Second * 5,
		},
	}
}

func (c *Client) GetResponseBodyFrom(url string) ([]byte, error) {
	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
