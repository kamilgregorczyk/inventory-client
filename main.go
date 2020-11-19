package main

import (
	"net/url"
	"test2/client"
	"test2/logger"
	"time"
)

func main() {
	inventoryClient := client.NewInventoryClient(time.Second, &url.URL{
		Scheme: "https",
		Host:   "inventory.raspicluster.pl"})
	log := logger.NewLogger("Main")
	items, err := inventoryClient.GetItems()
	if err != nil {
		log.Error.Printf(err.Error())
	} else {
		log.Info.Printf("Items: %+v", items)
	}
}
