package http

import "fmt"

type ClientError struct {
	Message string
	Url     string
}

func (e ClientError) Error() string {
	return fmt.Sprintf("Failed to call %s due to %s", e.Url, e.Message)
}

type ClientHttpError struct {
	Url          string
	StatusCode   int
	ResponseBody []byte
}

func (e ClientHttpError) Error() string {
	return fmt.Sprintf("Failed to call %s due to HTTP error %d %s", e.Url, e.StatusCode, e.ResponseBody)
}
