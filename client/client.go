package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"test2/logger"
	"time"
)

func NewInventoryClient(timeout time.Duration, url *url.URL) *InventoryClient {
	return &InventoryClient{
		url:  url,
		http: &http.Client{Timeout: timeout},
		log:  logger.NewLogger("InventoryClient"),
	}
}

type InventoryClient struct {
	url  *url.URL
	http *http.Client
	log  *logger.Logger
}

func (c *InventoryClient) GetItems() (*[]Inventory, error) {
	fullPath := c.url.String() + "/inventory"
	startTime := time.Now()
	c.logNewRequest("GET", fullPath)
	response, err := c.http.Get(fullPath)

	if err != nil {
		c.logFailedRequest("GET", fullPath, startTime)
		return nil, InventoryClientError{Message: fmt.Sprintf("Failed to make a request %s", err.Error()), Url: fullPath}
	}

	buffer, err := ioutil.ReadAll(response.Body)

	if response.StatusCode >= 400 || response.StatusCode < 200 {
		c.logFailedRequest("GET", fullPath, startTime)
		return nil, InventoryClientHttpError{ResponseBody: buffer, Url: fullPath, StatusCode: response.StatusCode}
	}

	if err != nil {
		c.logFailedRequest("GET", fullPath, startTime)
		return nil, InventoryClientError{Message: fmt.Sprintf("Failed due to IO error %s", err.Error()), Url: fullPath}
	}

	var items []Inventory
	err = json.Unmarshal(buffer, &items)
	if err != nil {
		c.logFailedRequest("GET", fullPath, startTime)
		return nil, InventoryClientError{Message: fmt.Sprintf("Failed due to parsing error %s", err.Error()), Url: fullPath}
	}

	c.logFinishedRequest("GET", fullPath, startTime)
	return &items, nil

}

func (c *InventoryClient) logNewRequest(method string, path string) {
	c.log.Info.Printf("Outgoing request [%s] %s", method, path)
}

func (c *InventoryClient) logFinishedRequest(method string, path string, startTime time.Time) {
	c.log.Info.Printf("Outoing request [%s] %s finished in %d ms", method, path, elapsedTime(startTime).Milliseconds())
}

func (c *InventoryClient) logFailedRequest(method string, path string, startTime time.Time) {
	c.log.Error.Printf("Outoing request [%s] %s failed in %d ms", method, path, elapsedTime(startTime).Milliseconds())
}

func elapsedTime(startTime time.Time) time.Duration {
	return time.Now().Sub(startTime)
}
