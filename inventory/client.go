package inventory

import (
	"fmt"
	"net/url"
	"test2/http"
	"time"
)

func New(config Config) *Client {
	client := http.New(http.Config{Timeout: config.Timeout})
	return &Client{
		Url:    config.Url,
		Client: client,
	}
}

type Config struct {
	Timeout time.Duration
	Url     url.URL
}

type Client struct {
	Url    url.URL
	Client *http.Client
}

func (c *Client) GetItems() ([]Inventory, error) {
	var items []Inventory
	path := fmt.Sprintf("%s/inventory", c.Url.String())
	err := c.Client.Get(path, &items)
	return items, err

}

func (c *Client) GetItem(id int) (Inventory, error) {
	var item Inventory
	path := fmt.Sprintf("%s/inventory/%d", c.Url.String(), id)
	err := c.Client.Get(path, &item)
	return item, err
}
