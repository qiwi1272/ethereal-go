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

	// place an order for ethusd
	p := products["ETHUSD"]

	order := p.NewOrder(rest.ORDER_LIMIT, 0.123, 1000.1, false, rest.BUY, rest.TIF_GTD)
	placed, err := client.CreateOrder(ctx, order)
	if err != nil {
		log.Fatalf("failed to place limit order: %v", err)
	}

	// cancel the order we just placed
	cancelReq := rest.NewCancelOrderFromCreated(&placed)
	cancelled, err := client.CancelOrder(ctx, cancelReq)
	if err != nil {
		log.Fatalf("failed to cancel limit order: %v", err)
	}

	// again
	placed, err = rest.Send(ctx, client, order)
	if err != nil {
		log.Fatalf("failed to place limit order: %v", err)
	}

	// cancel the order we just placed
	cancelReq = rest.NewCancelOrderFromCreated(&placed)
	cancelled, err = rest.Send(ctx, client, cancelReq)
	if err != nil {
		log.Fatalf("failed to cancel limit order: %v", err)
	}

	for _, c := range cancelled {
		log.Printf("Placed and cancelled order: %v", c)
	}

	time.Sleep(time.Second * 1)
}
