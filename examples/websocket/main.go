package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	restClient "github.com/qiwi1272/ethereal-go/rest_client"
	wssClient "github.com/qiwi1272/ethereal-go/websocket_client"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	rest, err := restClient.NewClient(ctx, os.Getenv("ETHEREAL_PK"), restClient.Testnet)
	if err != nil {
		panic(err)
	}
	sid := rest.Subaccount.Id

	var symbols map[string]restClient.Product
	if symbols, err = rest.GetProductMap(ctx); err != nil {
		panic(err)
	}

	ws := wssClient.NewClient(ctx)
	defer ws.Close()

	for symbolKey := range symbols {
		if err := ws.SubscribeMarketPrice(ctx, symbolKey); err != nil {
			log.Fatal("SubscribeMarketPrice:", err)
		}
		if err := ws.SubscribeBook(ctx, symbolKey); err != nil {
			log.Fatal("SubscribeBook:", err)
		}
		if err := ws.SubscribeFill(ctx, symbolKey); err != nil {
			log.Fatal("SubscribeFill:", err)
		}
	}

	if err := ws.SubscribeLiquidation(ctx, sid); err != nil {
		log.Fatal("SubscribeLiquidation:", err)
	}
	if err := ws.SubscribeOrderFill(ctx, sid); err != nil {
		log.Fatal("SubscribeOrderFill:", err)
	}
	if err := ws.SubscribeOrderUpdate(ctx, sid); err != nil {
		log.Fatal("SubscribeOrderUpdate:", err)
	}
	if err := ws.SubscribeTokenTransfer(ctx, sid); err != nil {
		log.Fatal("SubscribeTokenTransfer:", err)
	}

	ws.OnBook(func(diff *wssClient.L2Book) {
		fmt.Printf("called back L2Book: %v\n", diff)
	})
	ws.OnPrice(func(mp *wssClient.MarketPrice) {
		fmt.Printf("called back MarketPrice: %v\n", mp)
	})
	ws.OnLiquidation(func(sl *wssClient.SubaccountLiquidationEvent) {
		fmt.Printf("called back SubaccountLiquidation: %v\n", sl)
	})
	ws.OnOrderFill(func(of *wssClient.OrderFillEvent) {
		fmt.Printf("called back OrderFillEvent: %v\n", of)
	})
	ws.OnOrderUpdate(func(o *wssClient.OrderUpdateEvent) {
		fmt.Printf("called back OrderUpdateEvent: %v\n", o)
	})
	ws.OnTradeFill(func(tf *wssClient.TradeFillEvent) {
		fmt.Printf("called back TradeFillEvent: %v\n", tf)
	})
	ws.OnTransfer(func(t *wssClient.Transfer) {
		fmt.Printf("called back Transfer: %v\n", t)
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
