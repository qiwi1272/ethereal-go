package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	restClient "github.com/qiwi1272/ethereal-go/rest_client"
	socketioClient "github.com/qiwi1272/ethereal-go/socketio_client"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx := context.Background()
	// create client and fetch products
	rest, err := restClient.NewRestClient(ctx, os.Getenv("ETHEREAL_PK"), restClient.Mainnet)
	if err != nil {
		log.Fatalf("failed to init ethereal client: %v", err)
	}
	products, err := rest.GetProductMap(ctx)
	if err != nil {
		log.Fatalf("failed to fetch products: %v", err)
	}

	eth_perp := products["ETHUSD"]

	ws := socketioClient.NewSocketIOClient()

	ws.SubscribeToBook(eth_perp.ID)
	ws.SubscribeToPrice(eth_perp.ID)

	ws.SubscribeToFill(rest.Subaccount)
	ws.SubscribeToOrder(rest.Subaccount)

	ws.OnBookDepth(bookHandler)
	ws.OnPrice(priceHandler)

	ws.OnFill(fillHandler)
	ws.OnOrder(orderHandler)

	select {}
}

func bookHandler(v socketioClient.BookDepthStream) {
	fmt.Printf("BookDepth update: %+v\n", v)
}

func priceHandler(v socketioClient.MarketPriceStream) {
	fmt.Printf("Price update: %+v\n", v)
}
func fillHandler(v socketioClient.OrderFillStream) {
	fmt.Printf("Fill update: %+v\n", v)
}
func orderHandler(v socketioClient.OrderStream) {
	fmt.Printf("Order update: %+v\n", v)
}
