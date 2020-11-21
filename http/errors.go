package http

import (
	"errors"
	"fmt"
)

var (
	TimeoutBellowZeroError = errors.New("timeout has to be larger than 0ms")
)

type ClientError struct {
	Url     string
	Message string
	Err     error
}

func (e *ClientError) Error() string {
	return fmt.Sprintf("failed to call %s due to %s: %s", e.Url, e.Message, e.Err)
}

func (e *ClientError) Unwrap() error {
	return e.Err
}

type ClientHttpError struct {
	Url          string
	StatusCode   int
	ResponseBody []byte
}

func (e *ClientHttpError) Error() string {
	return fmt.Sprintf("failed to call %s due to HTTP error %d %s", e.Url, e.StatusCode, e.ResponseBody)
}
