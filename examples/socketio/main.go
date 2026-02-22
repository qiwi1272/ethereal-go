package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	etherealpb "github.com/qiwi1272/ethereal-go/_pb"
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
	rest, err := restClient.NewRestClient(ctx, "", restClient.Mainnet)
	if err != nil {
		log.Fatalf("failed to init ethereal client: %v", err)
	}
	products, err := rest.GetProductMap(ctx)
	if err != nil {
		log.Fatalf("failed to fetch products: %v", err)
	}

	eth_perp := products["ETHUSD"]

	ws := socketioClient.NewSocketIOClient()

	ws.OnBookDepth(bookHandler)

	time.Sleep(3 * time.Second)

	ws.SubscribeToBook(eth_perp.ID)

	select {}
}

func bookHandler(v *etherealpb.BookDiff) {
	fmt.Printf("BookDepth update: %+v\n", v.ProductId)
	fmt.Printf("BookDepth update: %+v\n", v.Bids)
	fmt.Printf("BookDepth update: %+v\n", v.Asks)
}
