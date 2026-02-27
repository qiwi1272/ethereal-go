package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	wssClient "github.com/qiwi1272/ethereal-go/websocket_client"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ws := wssClient.NewClient(ctx)
	defer ws.Close()

	if err := ws.SubscribePrice(ctx, "BTCUSD"); err != nil {
		log.Fatal("SubscribePrice:", err)
	}
	if err := ws.SubscribeBook(ctx, "BTCUSD"); err != nil {
		log.Fatal("SubscribeBook:", err)
	}

	ws.OnBook(func(diff *wssClient.L2Book) {
		fmt.Println(diff)
	})
	ws.OnPrice(func(mp *wssClient.MarketPrice) {
		fmt.Println(mp)
	})

	errCh := make(chan error, 1)
	go func() {
		errCh <- ws.Listen(ctx) // blocking
	}()

	select {
	case <-ctx.Done():
		return
	case err := <-errCh:
		if err != nil && ctx.Err() == nil {
			log.Fatal(err)
		}
	}

}
