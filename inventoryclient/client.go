package inventoryclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func New(timeout time.Duration, url *url.URL) *InventoryClient {
	return &InventoryClient{
		Url:  url,
		Http: &http.Client{Timeout: timeout},
	}
}

type InventoryClient struct {
	Url  *url.URL
	Http *http.Client
}

func (c *InventoryClient) GetItems() ([]Inventory, error) {
	fullPath := c.Url.String() + "/inventory"
	response, err := c.Http.Get(fullPath)

	if err != nil {
		return nil, InventoryClientError{Message: fmt.Sprintf("Failed to make a request %s", err.Error()), Url: fullPath}
	}

	buffer, err := ioutil.ReadAll(response.Body)

	if response.StatusCode >= 400 || response.StatusCode < 200 {
		return nil, InventoryClientHttpError{ResponseBody: buffer, Url: fullPath, StatusCode: response.StatusCode}
	}

	if err != nil {
		return nil, InventoryClientError{Message: fmt.Sprintf("Failed due to IO error %s", err.Error()), Url: fullPath}
	}

	var items []Inventory
	err = json.Unmarshal(buffer, &items)
	if err != nil {
		return nil, InventoryClientError{Message: fmt.Sprintf("Failed due to parsing error %s", err.Error()), Url: fullPath}
	}

	return items, nil
}
