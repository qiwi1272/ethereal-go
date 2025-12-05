package main

import (
	"context"
	"log"
	"time"

	"github.com/joho/godotenv"
	ethereal "github.com/qiwi1272/ethereal-go/client"
)

func main() {
	// load dotenv
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx := context.Background()

	// create client and fetch products
	client, err := ethereal.NewEtherealClient(ctx, "", ethereal.Mainnet)
	if err != nil {
		log.Fatalf("failed to init ethereal client: %v", err)
	}
	products, err := client.GetProductMap(ctx)
	if err != nil {
		log.Fatalf("failed to fetch products: %v", err)
	}

	// place 3 orders for ethusd
	eth_perp := products["ETHUSD"]

	orders := make([]*ethereal.Order, 3)
	for i := range orders {
		px := 1000.1 + float64(i)
		orders[i] = eth_perp.NewOrder(ethereal.ORDER_LIMIT, 0.123, px, false, ethereal.BUY, ethereal.TIF_GTD)
	}
	placed, err := client.BatchOrder(ctx, orders)
	if err != nil {
		log.Fatalf("failed to place limit order: %v", err)
	}

	// cancel the orders we just placed
	cancel := ethereal.NewCancelOrderFromCreated(placed...)
	cancelled, err := cancel.Send(ctx, client)
	if err != nil {
		log.Fatalf("failed to cancel limit order: %v", err)
	}

	log.Printf("Placed and cancelled orders: %s", cancelled)
	time.Sleep(time.Second * 1)
}
