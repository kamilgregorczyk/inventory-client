package main

import (
	"log"
	"net/url"
	"test2/inventory"
	"time"
)

func main() {
	inventoryClient := inventory.New(inventory.Config{
		Timeout: time.Second,
		Url: url.URL{
			Scheme: "https",
			Host:   "inventory.raspicluster.pl"},
	})
	items, err := inventoryClient.GetItems()
	if err != nil {
		log.Panicf(err.Error())
	} else {
		log.Printf("Items: %+v", items)
		item, err := inventoryClient.GetItem(items[0].Id)
		if err != nil {
			log.Panicf(err.Error())
		} else {
			log.Printf("Item: %+v", item)
		}

	}
}
