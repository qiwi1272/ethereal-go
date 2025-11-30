package ethereal

import (
	"fmt"
	"math/big"
	"strconv"
	"time"

	abi "github.com/ethereum/go-ethereum/signer/core/apitypes"
)

func getNonce() string {
	now := time.Now()
	return strconv.FormatInt(now.UnixNano(), 10)
}

// -------- ORDER CREATION --------
func (p *Product) NewOrder(orderType OrderType, qty float64, px float64, reduce bool, side int64, tif TimeInForce) *Order {
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

func (o *Order) toMessage() (abi.TypedDataMessage, error) {
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
	side := new(big.Int).SetInt64(o.Side)
	engine := new(big.Int).SetInt64(o.EngineType)
	id := new(big.Int).SetInt64(o.OnchainID)
	sigTs := new(big.Int).SetInt64(o.SignedAt)

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

func (o *Order) build(cl *EtherealClient) {
	o.Sender = cl.address
	o.Subaccount = cl.subaccount.Name
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

func (o *CancelOrder) toMessage() (abi.TypedDataMessage, error) {
	return abi.TypedDataMessage{
		"sender":     o.Sender,
		"subaccount": o.Subaccount,
		"nonce":      o.Nonce,
		//"orderIds":   co.OrderIDs,
	}, nil
}

func (o *CancelOrder) build(cl *EtherealClient) {
	o.Sender = cl.address
	o.Subaccount = cl.subaccount.Name

	o.Nonce = getNonce()
}
