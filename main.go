package main

import (
	"context"
	"log"
	"net/url"
	"test2/http/retry"
	"test2/inventory"
	"time"
)

func main() {
	inventoryClient := inventory.New(inventory.Config{
		Timeout: time.Second,
		Url: url.URL{
			Scheme: "https",
			Host:   "inventory.raspicluster.pl"},
		RetriesConfig: retry.RetriesConfig{
			MaxRetries: 3,
			Delay:      time.Second,
			Factor:     1.3,
		},
	})
	items, err := inventoryClient.GetItems(context.Background())
	if err != nil {
		log.Panicf(err.Error())
	} else {
		log.Printf("Items: %+v", items)
		item, err := inventoryClient.GetItem(context.Background(), items[0].Id)
		if err != nil {
			log.Panicf(err.Error())
		} else {
			log.Printf("Item: %+v", item)
		}

		item2, err := inventoryClient.CreateItem(context.Background(), inventory.CreateInventory{Name: "aa", Description: "cc"})
		if err != nil {
			log.Panicf(err.Error())
		} else {
			log.Printf("Item: %+v", item2)
		}
	}
}
