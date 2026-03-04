package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/qiwi1272/ethereal-go/rest"
	"github.com/qiwi1272/ethereal-go/socketio"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx := context.Background()
	// create client and fetch products
	rest, err := rest.NewClient(ctx, os.Getenv("ETHEREAL_PK"), rest.Mainnet)
	if err != nil {
		log.Fatalf("failed to init ethereal client: %v", err)
	}
	products, err := rest.GetProductMap(ctx)
	if err != nil {
		log.Fatalf("failed to fetch products: %v", err)
	}

	eth_perp := products["ETHUSD"]

	ws := socketio.NewClient(socketio.Mainnet)

	ws.OnBookDepth(bookHandler)

	time.Sleep(3 * time.Second)

	ws.SubscribeToBook(eth_perp.ID)

	select {}
}

func bookHandler(v socketio.BookDepthStream) {
	fmt.Printf("BookDepth update: %+v\n", v.ProductID)
	fmt.Printf("BookDepth update: %+v\n", v.Bids)
	fmt.Printf("BookDepth update: %+v\n", v.Asks)
}
