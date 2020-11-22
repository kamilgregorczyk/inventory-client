package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	corehttp "net/http"
	"test2/http/retry"
	"time"
)

type ClientConfig struct {
	Timeout time.Duration
	Retries retry.RetriesConfig
}

type Client struct {
	client *corehttp.Client
	retry  *retry.Retry
}

func NewClient(config ClientConfig) (*Client, error) {
	if config.Timeout.Milliseconds() <= 0 {
		return nil, TimeoutBellowZeroError
	}

	retry, err := retry.NewRetries(config.Retries)

	if err != nil {
		return nil, err
	}

	return &Client{
		client: &corehttp.Client{Timeout: config.Timeout},
		retry:  retry,
	}, nil
}

func (c *Client) Get(ctx context.Context, url string, model interface{}) error {
	response, err := c.execute(ctx, "GET", url, nil)
	if err != nil {
		return &ClientError{Message: "network error", Url: url, Err: err}
	}
	defer response.Body.Close()
	return readResponse(response, err, url, model)
}

func (c *Client) Post(ctx context.Context, url string, body interface{}, model interface{}) error {
	response, err := c.execute(ctx, "POST", url, body)
	if err != nil {
		return &ClientError{Message: "network error", Url: url, Err: err}
	}
	defer response.Body.Close()
	return readResponse(response, err, url, model)
}

func (c *Client) execute(context context.Context, method string, url string, body interface{}) (resp *corehttp.Response, err error) {
	marshaledBody, err := json.Marshal(body)

	if err != nil {
		return nil, fmt.Errorf("marshal error on url %s %w", url, err)
	}

	req, err := corehttp.NewRequestWithContext(context, method, url, bytes.NewBuffer(marshaledBody))
	if err != nil {
		return nil, &ClientError{Message: "network error", Url: url, Err: err}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return c.executeWithRetry(req)
}

func (c *Client) executeWithRetry(req *corehttp.Request) (*corehttp.Response, error) {
	return c.retry.Execute(func() (*corehttp.Response, error) {
		response, err := c.client.Do(req)
		if shouldRetry(response) {
			return response, &retry.RetryableError{Err: err}
		}
		return response, err
	})
}

func shouldRetry(response *corehttp.Response) bool {
	return response == nil || response.StatusCode >= 500
}

func readResponse(response *corehttp.Response, err error, url string, model interface{}) error {
	buffer, err := ioutil.ReadAll(response.Body)

	if response.StatusCode >= 400 || response.StatusCode < 200 {
		return &ClientHttpError{ResponseBody: buffer, Url: url, StatusCode: response.StatusCode}
	}

	if err != nil {
		return &ClientError{Message: "io error", Url: url, Err: err}
	}

	err = json.Unmarshal(buffer, model)
	if err != nil {
		return &ClientError{Message: "parsing error", Url: url, Err: err}
	}
	return nil
}
