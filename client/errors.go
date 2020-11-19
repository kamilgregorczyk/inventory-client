package client

import "fmt"

type InventoryClientError struct {
	Message string
	Url     string
}

func (e InventoryClientError) Error() string {
	return fmt.Sprintf("Failed to call %s due to %s", e.Url, e.Message)
}

type InventoryClientHttpError struct {
	Url          string
	StatusCode   int
	ResponseBody []byte
}

func (e InventoryClientHttpError) Error() string {
	return fmt.Sprintf("Failed to call %s due to HTTP error %d %s", e.Url, e.StatusCode, e.ResponseBody)
}
