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
