package main

import (
	"context"
	"log"
	"time"

	rest "roundinternet.money/ethereal-rest"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create client and fetch products
	client, err := rest.NewClient(ctx, "0bb5d63b84421e1268dda020818ae30cf26e7f10e321fb820a8aa69216dea92a", rest.Testnet)
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
