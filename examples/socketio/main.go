package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	pb "github.com/qiwi1272/ethereal-go/_pb"
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
	rest, err := restClient.NewRestClient(ctx, "", restClient.Testnet)
	if err != nil {
		log.Fatalf("failed to init ethereal client: %v", err)
	}
	products, err := rest.GetProductMap(ctx)
	if err != nil {
		log.Fatalf("failed to fetch products: %v", err)
	}

	eth_perp := products["ETHUSD"]

	fmt.Println(eth_perp.ID)

	ws := socketioClient.NewSocketIOClient()

	ws.Socket.OnEvent("BookDepth", func(args ...any) {
		fmt.Printf("argc=%d\n", len(args))
		for i, a := range args {
			fmt.Printf("  arg[%d] type=%T val=%#v\n", i, a, a)
		}
	})
	ws.OnPrice(priceHandler)

	time.Sleep(5 * time.Second)

	ws.SubscribeToBook(eth_perp.ID)
	ws.SubscribeToPrice(eth_perp.ID)

	select {}
}

func bookHandler(v *pb.BookDiff) {
	fmt.Printf("BookDepth update: %+v\n", v.Symbol)
	fmt.Printf("BookDepth update: %+v\n", v.Bids)
	fmt.Printf("BookDepth update: %+v\n", v.Asks)
}

func priceHandler(v socketioClient.MarketPriceStream) {
	fmt.Printf("priceHandler update: %+v\n", v)
}
