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
	client, err := ethereal.NewEtherealClient(ctx, "")
	if err != nil {
		log.Fatalf("failed to init ethereal client: %v", err)
	}
	products, err := client.GetProductMap(ctx)
	if err != nil {
		log.Fatalf("failed to fetch products: %v", err)
	}

	// place an order for ethusd
	p := products["ETHUSD"]

	order := p.NewOrder(ethereal.ORDER_LIMIT, 0.123, 1000.1, false, 0, ethereal.TIF_GTD)
	placed, err := order.Send(ctx, client)
	if err != nil {
		log.Fatalf("failed to place limit order: %v", err)
	}

	// cancel the order we just placed
	cancel := ethereal.NewCancelOrder(placed.Id)
	cancelled, err := cancel.Send(ctx, client)
	if err != nil {
		log.Fatalf("failed to cancel limit order: %v", err)
	}

	log.Printf("Placed and cancelled order: %s", cancelled)
	time.Sleep(time.Second * 1)
}
