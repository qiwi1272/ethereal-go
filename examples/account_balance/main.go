package main

import (
	"context"
	"fmt"
	"log"
	"os"

	rest "github.com/roundinternetmoney/ethereal-rest"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pk := os.Getenv("ETHEREAL_PK")
	if pk == "" {
		log.Fatal("ETHEREAL_PK is required (hex private key, with or without 0x prefix)")
	}

	client, err := rest.NewClient(ctx, pk, rest.Testnet)
	if err != nil {
		log.Fatalf("failed to init ethereal client: %v", err)
	}
	// fetch balance
	balance, err := client.GetAccountBalance(ctx)
	if err != nil {
		log.Fatalf("failed to get balance: %v", err)
	}

	for _, b := range balance {
		fmt.Println(b)
	}

}
