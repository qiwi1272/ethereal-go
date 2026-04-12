package main

import (
	"context"
	"log"
	"os"
	"time"

	rest "github.com/roundinternetmoney/ethereal-rest"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pk := os.Getenv("ETHEREAL_PK")
	if pk == "" {
		log.Fatal("ETHEREAL_PK is required (hex private key, with or without 0x prefix)")
	}

	client, err := rest.NewClient(ctx, pk, rest.Testnet)
	if err != nil {
		log.Fatalf("failed to init ethereal client: %v", err)
	}
	products, err := client.GetProductMap(ctx)
	if err != nil {
		log.Fatalf("failed to fetch products: %v", err)
	}

	// place 3 orders for ethusd
	eth_perp := products["ETHUSD"]

	orders := make([]*rest.Order, 3)
	for i := range orders {
		px := 1000.1 + float64(i)
		orders[i] = eth_perp.NewOrder(rest.ORDER_LIMIT, 0.123, px, false, rest.BUY, rest.TIF_GTD)
	}
	placed, err := client.CreateOrders(ctx, orders)
	if err != nil {
		log.Fatalf("failed to place limit order: %v", err)
	}

	// cancel the orders we just placed
	cancelled, err := client.CancelOrdersFromCreated(ctx, placed)
	if err != nil {
		log.Fatalf("failed to cancel limit order: %v", err)
	}

	for _, c := range cancelled {
		log.Printf("Placed and cancelled order: %v", c)
	}

	time.Sleep(time.Second * 1)
}
