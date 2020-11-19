package main

import (
	"log"
	"net/url"
	"test2/client"
	"time"
)

func main() {
	inventoryClient := client.NewInventoryClient(time.Second, &url.URL{
		Scheme: "https",
		Host:   "inventory.raspicluster.pl"})
	items, err := inventoryClient.GetItems()
	if err != nil {
		log.Panicf(err.Error())
	} else {
		log.Printf("Items: %+v", items)
	}
}
