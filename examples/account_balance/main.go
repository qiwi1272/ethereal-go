package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	restClient "github.com/qiwi1272/ethereal-go/rest_client"
)

func main() {
	// load dotenv
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx := context.Background()

	// create client and fetch products
	client, err := restClient.NewRestClient(ctx, os.Getenv("ETHEREAL_PK"), restClient.Testnet)
	if err != nil {
		log.Fatalf("failed to init ethereal client: %v", err)
	}
	// fetch balance
	balance, err := client.GetAccountBalance(ctx)
	if err != nil {
		log.Fatalf("failed to get balance: %v", err)
	}

	log.Println(balance)

}
