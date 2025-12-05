package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	ethereal "github.com/qiwi1272/ethereal-go/client"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx := context.Background()
	// create client and fetch products
	rest, err := ethereal.NewEtherealClient(ctx, "", ethereal.Testnet)
	if err != nil {
		log.Fatalf("failed to init ethereal client: %v", err)
	}
	products, err := rest.GetProductMap(ctx)
	if err != nil {
		log.Fatalf("failed to fetch products: %v", err)
	}

	eth_perp := products["ETHUSD"]

	ws := ethereal.NewWebSocketClient()

	ws.SubscribeToBook(&eth_perp)
	ws.SubscribeToPrice(&eth_perp)

	ws.SubscribeToFill(rest.Subaccount)
	ws.SubscribeToOrder(rest.Subaccount)

	ws.OnBookDepth(bookHandler)
	ws.OnPrice(priceHandler)

	ws.OnFill(fillHandler)
	ws.OnOrder(orderHandler)

	select {}
}

func bookHandler(v ethereal.BookDepthStream) {
	fmt.Printf("BookDepth update: %+v\n", v)
}

func priceHandler(v ethereal.MarketPriceStream) {
	fmt.Printf("Price update: %+v\n", v)
}
func fillHandler(v ethereal.OrderFillStream) {
	fmt.Printf("Fill update: %+v\n", v)
}
func orderHandler(v ethereal.OrderStream) {
	fmt.Printf("Order update: %+v\n", v)
}
