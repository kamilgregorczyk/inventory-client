package inventory

import (
	"context"
	"fmt"
	"net/url"
	"test2/http"
	"test2/http/retry"
	"time"
)

func New(config Config) *Client {
	client := http.New(http.Config{Timeout: config.Timeout, Retries: config.RetriesConfig})
	return &Client{
		Url:    config.Url,
		Client: client,
	}
}

type Config struct {
	Timeout       time.Duration
	Url           url.URL
	RetriesConfig retry.RetriesConfig
}

type Client struct {
	Url    url.URL
	Client *http.Client
}

func (c *Client) GetItems(ctx context.Context) ([]Inventory, error) {
	var items []Inventory
	path := fmt.Sprintf("%s/inventory", c.Url.String())
	err := c.Client.Get(ctx, path, &items)
	return items, err

}

func (c *Client) GetItem(ctx context.Context, id int) (Inventory, error) {
	var item Inventory
	path := fmt.Sprintf("%s/inventory/%d", c.Url.String(), id)
	err := c.Client.Get(ctx, path, &item)
	return item, err
}

func (c *Client) CreateItem(ctx context.Context, createInventory CreateInventory) (Inventory, error) {
	var item Inventory
	path := fmt.Sprintf("%s/inventory", c.Url.String())
	err := c.Client.Post(ctx, path, createInventory, &item)
	return item, err
}
