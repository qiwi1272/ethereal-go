Ethereal Go Client
===================

Lightweight golang client for interacting with the Ethereal API.

## Features

- EIP-712 data signing
- REST order placement and cancellation
- Batch execution support (concurrent, unordered, type-safe)
- Automatic nonce and timestamp handling
- Minimal dependencies
- Easy to extend with new message types
- Socket.IO streaming support (WIP)

Getting started
---------------
- Requires Go 1.25+.
- Install: `go get github.com/qiwi1272/ethereal-go`

Example
-------
```go
package main

import (
    "context"
    "log"

    "github.com/qiwi1272/ethereal-go/ethereal"
)

func main() {
    ctx := context.Background()

    // Uses ETHEREAL_PK if private key argument is empty
    client, err := ethereal.NewEtherealClient(ctx, "")
    if err != nil {
        log.Fatal(err)
    }
	products, err := client.GetProductMap(ctx)
	if err != nil {
		log.Fatalf("failed to fetch products: %v", err)
	}

    product := products["ETHUSD"]
    order := product.NewOrder(ethereal.ORDER_LIMIT, 0.123, 1000.1, false, 0, ethereal.TIF_GTD)

	placed, err := order.Send(ctx, client)
	if err != nil {
		log.Fatalf("failed to place limit order: %v", err)
	}

    log.Printf("Order created: %+v\n", res)
}
```
For more complete usage examples (batching, cancel orders, typed data inspection, etc.),
see the [examples/](./examples/) folder in this repository.

Configuration Notes
-----
- If no private key is passed to `NewEtherealClient`, the library uses the `ETHEREAL_PK` environment variable.
- Only one subaccount is currently supported; by default the first one discovered is used.
- All signable request messages implement the `Signable` interface.


Status
-----
- Work in progress.
- Current active development: Socket.IO streaming (order books, trades, events).
