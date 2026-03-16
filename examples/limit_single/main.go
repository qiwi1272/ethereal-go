package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/qiwi1272/ethereal-go"
	"github.com/qiwi1272/ethereal-go/rest"
)

func main() {
	// load dotenv
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx := context.Background()

	// create client and fetch products
	client, err := rest.NewClient(ctx, os.Getenv("ETHEREAL_PK"), rest.Testnet)
	if err != nil {
		log.Fatalf("failed to init ethereal client: %v", err)
	}
	products, err := client.GetProductMap(ctx)
	if err != nil {
		log.Fatalf("failed to fetch products: %v", err)
	}

	// place an order for ethusd
	p := products["ETHUSD"]

	order := p.NewOrder(ethereal.ORDER_LIMIT, 0.123, 1000.1, false, ethereal.BUY, ethereal.TIF_GTD)
	placed, err := client.CreateOrder(ctx, order)
	if err != nil {
		log.Fatalf("failed to place limit order: %v", err)
	}

	// cancel the order we just placed
	cancel := ethereal.NewCancel(placed.Id)
	cancelled, err := client.CancelOrder(ctx, cancel)
	if err != nil {
		log.Fatalf("failed to cancel limit order: %v", err)
	}

	// again
	placed, err = rest.Send(ctx, client, order)
	if err != nil {
		log.Fatalf("failed to place limit order: %v", err)
	}

	// cancel the order we just placed
	cancelled, err = rest.Send(ctx, client, cancel)
	if err != nil {
		log.Fatalf("failed to cancel limit order: %v", err)
	}

	log.Printf("Placed and cancelled order: %s", cancelled)
	time.Sleep(time.Second * 1)
}
