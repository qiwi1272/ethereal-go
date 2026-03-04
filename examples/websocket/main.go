package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/qiwi1272/ethereal-go"
	"github.com/qiwi1272/ethereal-go/pb"
	"github.com/qiwi1272/ethereal-go/rest"
	"github.com/qiwi1272/ethereal-go/websocket"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	rest, err := rest.NewClient(ctx, os.Getenv("ETHEREAL_PK"), rest.Testnet)
	if err != nil {
		panic(err)
	}
	sid := rest.Subaccount.Id

	var symbols map[string]ethereal.Product
	if symbols, err = rest.GetProductMap(ctx); err != nil {
		panic(err)
	}

	ws := websocket.NewClient(ctx, websocket.Testnet)
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

	ws.OnBook(func(diff *pb.L2Book) {
		//fmt.Printf("called back L2Book: %v\n", diff)
	})
	ws.OnPrice(func(mp *pb.MarketPrice) {
		//fmt.Printf("called back MarketPrice: %v\n", mp)
	})
	ws.OnLiquidation(func(sl *pb.SubaccountLiquidationEvent) {
		fmt.Printf("called back SubaccountLiquidation: %v\n", sl)
	})
	ws.OnOrderFill(func(of *pb.OrderFillEvent) {
		//fmt.Printf("called back OrderFillEvent: %v\n", of)
	})
	ws.OnOrderUpdate(func(o *pb.OrderUpdateEvent) {
		//fmt.Printf("called back OrderUpdateEvent: %v\n", o)
	})
	ws.OnTradeFill(func(tf *pb.TradeFillEvent) {
		//fmt.Printf("called back TradeFillEvent: %v\n", tf)
	})
	ws.OnTransfer(func(t *pb.Transfer) {
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
