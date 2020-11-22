package http

import (
	"bytes"
	"context"
	"encoding/json"
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
	request, err := c.createRequest(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	response, err := c.executeWithRetry(request)
	if err != nil {
		return err
	}

	return readResponse(response, err, url, model)
}

func (c *Client) Post(ctx context.Context, url string, requestBody interface{}, responseBody interface{}) error {
	request, err := c.createRequest(ctx, "POST", url, requestBody)
	if err != nil {
		return err
	}

	response, err := c.executeWithRetry(request)
	if err != nil {
		return err
	}

	return readResponse(response, err, url, responseBody)
}

func (c *Client) createRequest(context context.Context, method string, url string, requestBody interface{}) (resp *corehttp.Request, err error) {
	marshaledBody, err := json.Marshal(requestBody)

	if err != nil {
		return nil, &ClientError{Message: "body parse error", Url: url, Err: err}
	}

	req, err := corehttp.NewRequestWithContext(context, method, url, bytes.NewBuffer(marshaledBody))
	if err != nil {
		return nil, &ClientError{Message: "network error", Url: url, Err: err}
	}
	c.setHeaders(req)
	return req, nil
}

func (c *Client) setHeaders(req *corehttp.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
}

func (c *Client) executeWithRetry(request *corehttp.Request) (*corehttp.Response, error) {
	response, err := c.retry.Execute(func() (*corehttp.Response, error) {
		response, err := c.client.Do(request)
		if shouldRetry(response, err) {
			return response, &retry.RetryableError{Err: err}
		}
		return response, err
	})

	if err != nil {
		return response, &ClientError{Message: "network error", Url: request.URL.String(), Err: err}
	}
	if response != nil && response.StatusCode >= 400 {
		return response, &ClientHttpError{Url: request.URL.String(), StatusCode: response.StatusCode}
	}

	return response, nil
}

func shouldRetry(response *corehttp.Response, err error) bool {
	return err != nil || response == nil || response.StatusCode >= 500
}

func readResponse(response *corehttp.Response, err error, url string, responseBody interface{}) error {
	buffer, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		return &ClientError{Message: "io error", Url: url, Err: err}
	}

	err = json.Unmarshal(buffer, responseBody)
	if err != nil {
		return &ClientError{Message: "parsing error", Url: url, Err: err}
	}
	return nil
}
