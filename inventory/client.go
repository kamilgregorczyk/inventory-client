package inventory

import (
	"context"
	"fmt"
	"net/url"
	"test2/http"
	"test2/http/retry"
	"time"
)

type ClientConfig struct {
	Timeout       time.Duration
	Url           url.URL
	RetriesConfig retry.RetriesConfig
}

type Client struct {
	Url    url.URL
	Client *http.Client
}

func NewClient(config ClientConfig) (*Client, error) {
	client, err := http.NewClient(http.ClientConfig{
		Timeout: config.Timeout,
		Retries: config.RetriesConfig,
		Headers: http.Headers{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		}})
	if err != nil {
		return nil, err
	}
	return &Client{
		Url:    config.Url,
		Client: client,
	}, nil
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
