package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Config struct {
	Timeout time.Duration
	MaxRetries uint

}

type Client struct {
	http *http.Client
}

func New(config Config) *Client {
	return &Client{http: &http.Client{Timeout: config.Timeout},
	}
}

func (c *Client) Get(url string, model interface{}) error {
	response, err := c.http.Get(url)

	if err != nil {
		return ClientError{Message: fmt.Sprintf("Failed to make a request %s", err.Error()), Url: url}
	}

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
