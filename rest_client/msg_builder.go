package rest_client

import (
	"fmt"
	"math/big"
	"strconv"
	"time"

	abi "github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/qiwi1272/ethereal-go"
)

func getNonce() string {
	now := time.Now()
	return strconv.FormatInt(now.UnixNano(), 10)
}

// -------- ORDER CREATION --------
func NewRawOrder(
	orderType ethereal.OrderType,
	onchainId int64,
	marketType ethereal.EngineType,
	qty float64,
	px float64,
	reduce bool,
	side ethereal.OrderSide,
	tif ethereal.TimeInForce,
) *Order {
	return &Order{
		Type:        orderType,
		Quantity:    fmt.Sprintf("%.9f", qty),
		Side:        side,
		OnchainID:   onchainId,
		EngineType:  marketType,
		ReduceOnly:  reduce,
		Price:       fmt.Sprintf("%.9f", px),
		TimeInForce: tif,
		PostOnly:    false,
	}
}

func (p *Product) NewOrder(
	orderType ethereal.OrderType,
	qty float64,
	px float64,
	reduce bool,
	side ethereal.OrderSide,
	tif ethereal.TimeInForce,
) *Order {
	return &Order{
		Type:        orderType,
		Quantity:    fmt.Sprintf("%.9f", qty),
		Side:        side,
		OnchainID:   p.OnchainID,
		EngineType:  p.EngineType,
		ReduceOnly:  reduce,
		Price:       fmt.Sprintf("%.9f", px),
		TimeInForce: tif,
		PostOnly:    false,
	}
}

func (o *Order) ToMessage() (abi.TypedDataMessage, error) {
	qtyBig, err := Scale1e9(o.Quantity)
	if err != nil {
		return abi.TypedDataMessage{}, err
	}
	priceBig, err := Scale1e9(o.Price)
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

func (o *Order) build(cl *RestClient) {
	o.Sender = cl.Address
	o.Subaccount = cl.Subaccount.Name
	nonce := getNonce()

	o.Nonce = nonce
	o.SignedAt, _ = strconv.ParseInt(nonce[:len(nonce)-9], 10, 64)
}

// -------- ORDER CANCELLATION --------
func NewCancelOrderFromCreated(orders ...OrderCreated) *CancelOrder {
	ids := make([]string, len(orders))
	for i, o := range orders {
		ids[i] = o.Id
	}
	return NewCancelOrder(ids...)
}

func NewCancelOrder(oids ...string) *CancelOrder {
	return &CancelOrder{
		OrderIDs: oids,
	}
}

func (o *CancelOrder) ToMessage() (abi.TypedDataMessage, error) {
	return abi.TypedDataMessage{
		"sender":     o.Sender,
		"subaccount": o.Subaccount,
		"nonce":      o.Nonce,
		//"orderIds":   co.OrderIDs,
	}, nil
}

func (o *CancelOrder) build(cl *RestClient) {
	o.Sender = cl.Address
	o.Subaccount = cl.Subaccount.Name

	o.Nonce = getNonce()
}
