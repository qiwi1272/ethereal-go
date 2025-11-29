Ethereal Go Client
===================

Lightweight golang client for interacting with the Ethereal API.

Getting started
---------------
- Require Go 1.25+.
- Install: `go get github.com/qiwi1272/ethereal-go`

Example
-------
```go
package main

import (
	"context"
	"log"
	"time"

	client "github.com/qiwi1272/ethereal-go/client"
)

func main() {
	ctx := context.Background()

	// Supplying the key via environment: ETHEREAL_PK=0xabc...
	// Alternativley, use github.com/joho/godotenv
	cl, err := client.NewEtherealClient(ctx, "")
	if err != nil {
		log.Fatal(err)
	}

	products, err := cl.GetProductMap(ctx)
	if err != nil {
		log.Fatal(err)
	}
	eth_perp := products["ETH-PERP"]

	order := perp.NewLimitOrder(0.123, 1000.1, false, 0) // side 0 = buy, 1 = sell
	order.ClientOrderID = "exampleOrder"

	created, err := order.Send(ctx, cl)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("order created: %+v", created)

	cancel := client.NewCancelOrderFromCreated(placed)
	canceled, err := cancel.Send(ctx, cl)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("order canceled: %+v", canceled)
}
```

Notes
-----
- The client uses `ETHEREAL_PK` when no private key string is passed to `NewEtherealClient`.
- Currently only one subaccount is supported per client. By default the first one is used.
- This is a WIP, I'll expand the scope as I find time. 
- Build to be generic enough that implementing new methods is trivial.
- Immediate focus is on batch submissions and socket.io