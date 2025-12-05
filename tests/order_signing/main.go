package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	abi "github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/joho/godotenv"
	ethereal "github.com/qiwi1272/ethereal-go/client"
)

func scale1e9(s string) (*big.Int, error) {
	r := new(big.Rat)
	if _, ok := r.SetString(s); !ok {
		return nil, fmt.Errorf("bad decimal %q", s)
	}
	r.Mul(r, big.NewRat(1_000_000_000, 1))
	n := new(big.Int)
	n.Div(r.Num(), r.Denom())
	return n, nil
}

func toMessage(o ethereal.Order) (abi.TypedDataMessage, error) {
	qtyBig, err := scale1e9(o.Quantity)
	if err != nil {
		return abi.TypedDataMessage{}, err
	}
	priceBig, err := scale1e9(o.Price)
	if err != nil {
		return abi.TypedDataMessage{}, err
	}

	// even though we expect these values to be uint8 according to their signatures,
	// setting them as native uint8 raises a compiler error. strings or big ints are accepted.
	engine := big.NewInt(int64(o.EngineType))
	side := big.NewInt(int64(o.Side))
	id := big.NewInt(o.OnchainID)
	sigTs := big.NewInt(o.SignedAt)

	return abi.TypedDataMessage{
		"sender":     o.Sender,
		"subaccount": o.Subaccount,
		"quantity":   qtyBig,
		"price":      priceBig,
		"reduceOnly": o.ReduceOnly,
		"side":       side,
		"engineType": engine,
		"productId":  id,
		"nonce":      o.Nonce,
		"signedAt":   sigTs,
	}, nil
}

func TestOrders() {
	orderType := abi.TypedData{
		Types: abi.Types{"TradeOrder": []abi.Type{
			{Name: "sender", Type: "address"},
			{Name: "subaccount", Type: "bytes32"},
			{Name: "quantity", Type: "uint128"},
			{Name: "price", Type: "uint128"},
			{Name: "reduceOnly", Type: "bool"},
			{Name: "side", Type: "uint8"},
			{Name: "engineType", Type: "uint8"},
			{Name: "productId", Type: "uint32"},
			{Name: "nonce", Type: "uint64"},
			{Name: "signedAt", Type: "uint64"},
		}},
	}

	order := ethereal.Order{
		Sender:     "0xdeadbeef00000000000000000000000000000000",
		Subaccount: "0x123456789abcde00000000000000000000000000000000000000000000000000",
		Quantity:   "1",
		Price:      "3000",
		ReduceOnly: false,
		Side:       ethereal.BUY,
		EngineType: ethereal.PERPETUAL,
		OnchainID:  2, // later -> ProductId
		Nonce:      "1764897077655477722",
		SignedAt:   int64(1764897077),
	}

	message, err := toMessage(order) // clone from signing.go
	if err != nil {
		panic(err)
	}

	SenderBytes, err := orderType.EncodePrimitiveValue(orderType.Types["TradeOrder"][0].Type, message["sender"], 2)
	if err != nil {
		panic(err)
	}
	if common.Bytes2Hex(SenderBytes) != "000000000000000000000000deadbeef00000000000000000000000000000000" {
		panic("SenderBytes")
	}

	SubaccountBytes, err := orderType.EncodePrimitiveValue(orderType.Types["TradeOrder"][1].Type, message["subaccount"], 2)
	if err != nil {
		panic(err)
	}
	if common.Bytes2Hex(SubaccountBytes) != "123456789abcde00000000000000000000000000000000000000000000000000" {
		panic("SubaccountBytes")
	}

	QuantityBytes, err := orderType.EncodePrimitiveValue(orderType.Types["TradeOrder"][2].Type, message["quantity"], 2)
	if err != nil {
		panic(err)
	}
	if common.Bytes2Hex(QuantityBytes) != "000000000000000000000000000000000000000000000000000000003b9aca00" {
		panic("QuantityBytes")
	}

	PriceBytes, err := orderType.EncodePrimitiveValue(orderType.Types["TradeOrder"][3].Type, message["price"], 2)
	if err != nil {
		panic(err)
	}
	if common.Bytes2Hex(PriceBytes) != "000000000000000000000000000000000000000000000000000002ba7def3000" {
		panic("PriceBytes")
	}

	ReduceOnlyBytes, err := orderType.EncodePrimitiveValue(orderType.Types["TradeOrder"][4].Type, message["reduceOnly"], 2)
	if err != nil {
		panic(err)
	}
	if common.Bytes2Hex(ReduceOnlyBytes) != "0000000000000000000000000000000000000000000000000000000000000000" {
		panic("ReduceOnlyBytes")
	}

	SideBytes, err := orderType.EncodePrimitiveValue(orderType.Types["TradeOrder"][5].Type, message["side"], 2)
	if err != nil {
		panic(err)
	}
	if common.Bytes2Hex(SideBytes) != "0000000000000000000000000000000000000000000000000000000000000000" {
		panic("SideBytes")
	}

	EngineTypeBytes, err := orderType.EncodePrimitiveValue(orderType.Types["TradeOrder"][6].Type, message["engineType"], 2)
	if err != nil {
		panic(err)
	}
	if common.Bytes2Hex(EngineTypeBytes) != "0000000000000000000000000000000000000000000000000000000000000000" {
		panic("EngineTypeBytes")
	}

	OnchainIDBytes, err := orderType.EncodePrimitiveValue(orderType.Types["TradeOrder"][7].Type, message["productId"], 2)
	if err != nil {
		panic(err)
	}
	if common.Bytes2Hex(OnchainIDBytes) != "0000000000000000000000000000000000000000000000000000000000000002" {
		panic("OnchainIDBytes")
	}

	NonceBytes, err := orderType.EncodePrimitiveValue(orderType.Types["TradeOrder"][8].Type, message["nonce"], 2)
	if err != nil {
		panic(err)
	}
	if common.Bytes2Hex(NonceBytes) != "000000000000000000000000000000000000000000000000187e2c8a92c79dda" {
		panic("NonceBytes")
	}

	SignedAtBytes, err := orderType.EncodePrimitiveValue(orderType.Types["TradeOrder"][9].Type, message["signedAt"], 2)
	if err != nil {
		panic(err)
	}
	if common.Bytes2Hex(SignedAtBytes) != "0000000000000000000000000000000000000000000000000000000069323135" {
		panic("SignedAtBytes")
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	TestOrders()
	log.Println("Order validated")
}
