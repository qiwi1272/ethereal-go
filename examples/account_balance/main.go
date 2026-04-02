package main

import (
	"context"
	"fmt"
	"log"

	rest "github.com/roundinternetmoney/ethereal-rest"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create client and fetch products
	client, err := rest.NewClient(ctx, "0bb5d63b84421e1268dda020818ae30cf26e7f10e321fb820a8aa69216dea92a", rest.Testnet)
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
