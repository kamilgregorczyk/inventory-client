package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Config struct {
	Timeout    time.Duration
	MaxRetries uint
}

type Client struct {
	client *http.Client
}

func New(config Config) *Client {
	return &Client{client: &http.Client{Timeout: config.Timeout},
	}
}

func (c *Client) Get(ctx context.Context, url string, model interface{}) error {
	response, err := c.execute(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	return readResponse(response, err, url, model)
}

func (c *Client) Post(ctx context.Context, url string, body interface{}, model interface{}) error {
	response, err := c.execute(ctx, "POST", url, body)
	return readResponse(response, err, url, model)
}

func (c *Client) execute(context context.Context, method string, url string, body interface{}) (resp *http.Response, err error) {
	marshaledBody, err := json.Marshal(body)

	if err != nil {
		return nil, ClientError{Message: fmt.Sprintf("marshal error %s", err.Error()), Url: url}
	}
	req, err := http.NewRequestWithContext(context, method, url, bytes.NewBuffer(marshaledBody))
	if err != nil {
		return nil, ClientError{Message: fmt.Sprintf("network error %s", err.Error()), Url: url}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return c.client.Do(req)
}

func readResponse(response *http.Response, err error, url string, model interface{}) error {
	buffer, err := ioutil.ReadAll(response.Body)

	if response.StatusCode >= 400 || response.StatusCode < 200 {
		return ClientHttpError{ResponseBody: buffer, Url: url, StatusCode: response.StatusCode}
	}

	if err != nil {
		return ClientError{Message: fmt.Sprintf("io error %s", err.Error()), Url: url}
	}

	err = json.Unmarshal(buffer, model)
	if err != nil {
		return ClientError{Message: fmt.Sprintf("parsing error %s", err.Error()), Url: url}
	}
	return nil
}
