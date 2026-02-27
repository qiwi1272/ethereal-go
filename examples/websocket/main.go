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

	log.Println(sid)

	ws := wssClient.NewClient(ctx)
	defer ws.Close()

	if err := ws.SubscribeMarketPrice(ctx, "BTCUSD"); err != nil {
		log.Fatal("SubscribeMarketPrice:", err)
	}
	if err := ws.SubscribeBook(ctx, "BTCUSD"); err != nil {
		log.Fatal("SubscribeBook:", err)
	}
	// if err := ws.SubscribeFill(ctx, "BTCUSD"); err != nil {
	// 	log.Fatal("SubscribeFill:", err)
	// }
	if err := ws.SubscribeLiquidation(ctx, sid); err != nil {
		log.Fatal("SubscribeLiquidation:", err)
	}
	// if err := ws.SubscribeOrderFill(ctx, sid); err != nil {
	// 	log.Fatal("SubscribeOrderFill:", err)
	// }
	// if err := ws.SubscribeOrderUpdate(ctx, sid); err != nil {
	// 	log.Fatal("SubscribeOrderUpdate:", err)
	// }
	if err := ws.SubscribeTokenTransfer(ctx, sid); err != nil {
		log.Fatal("SubscribeTokenTransfer:", err)
	}

	ws.OnBook(func(diff *wssClient.L2Book) {
		fmt.Println(diff)
	})
	ws.OnPrice(func(mp *wssClient.MarketPrice) {
		fmt.Println(mp)
	})
	ws.OnLiquidation(func(sl *wssClient.SubaccountLiquidation) {
		fmt.Println(sl)
	})
	ws.OnOrderFill(func(of *wssClient.OrderFill) {
		fmt.Println(of)
	})
	ws.OnOrderUpdate(func(o *wssClient.Order) {
		fmt.Println(o)
	})
	ws.OnTradeFill(func(tf *wssClient.TradeFill) {
		fmt.Println(tf)
	})
	ws.OnTransfer(func(t *wssClient.Transfer) {
		fmt.Println(t)
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
