package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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
	rest, err := restClient.NewClient(ctx, os.Getenv("ETHEREAL_PK"), restClient.Mainnet)
	if err != nil {
		log.Fatalf("failed to init ethereal client: %v", err)
	}
	products, err := rest.GetProductMap(ctx)
	if err != nil {
		log.Fatalf("failed to fetch products: %v", err)
	}

	eth_perp := products["ETHUSD"]

	ws := socketioClient.NewClient()

	ws.OnBookDepth(bookHandler)

	time.Sleep(3 * time.Second)

	ws.SubscribeToBook(eth_perp.ID)

	select {}
}

func bookHandler(v *socketioClient.BookDiff) {
	fmt.Printf("BookDepth update: %+v\n", v.ProductId)
	fmt.Printf("BookDepth update: %+v\n", v.Bids)
	fmt.Printf("BookDepth update: %+v\n", v.Asks)
}
