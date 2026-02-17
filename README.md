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
- Socket.IO streaming support

Getting started
---------------
- Requires Go 1.25+.
- Install: `go get github.com/qiwi1272/ethereal-go`

Example Usage
-------
From the client directory:
`make all`
`ETHEREAL_PK=0x0000 bin/example_account_balance`

For more complete usage examples (batching, cancel orders, websocket subscriptions, etc.),
see the [examples/](./examples/) folder in this repository.


```
Configuration Notes
-----
- If no private key is passed to `NewEtherealClient`, the library uses the `ETHEREAL_PK` environment variable.
- All signable request messages implement the `Signable` interface.
- Only one subaccount is currently supported; by default the first one discovered is used.

Status
-----
- Most of the client is complete, and easy to expand.
- Other methods will be added as needed.


Contributing
-------------
Contributions are welcome! Please open issues or pull requests as needed.

## Code Formatting
To format the code, use the following command:

```bash
make fmt
```
## Dependency Management
To tidy up dependencies, use the following command:

```bash
make tidy
```

## To install dependencies, use the following command:

```bash
make deps
```

## Testing

To run tests, use the following command:

```bash
make test
```
## Building
To build the client, use the following command:
```bash
make build
```

## All
To run all common tasks (formatting, tidying, vetting, testing, building), use the following command:

```bash
make all
```

