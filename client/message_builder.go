package ethereal

import "fmt"

func (p *Product) NewLimitOrder(qty float64, px float64, reduce bool, side int64) *LimitOrder {
	return &LimitOrder{
		Type:        "LIMIT",
		Quantity:    fmt.Sprintf("%.9f", qty),
		Side:        side,
		OnchainID:   p.OnchainID,
		EngineType:  p.EngineType,
		ReduceOnly:  reduce,
		Price:       fmt.Sprintf("%.9f", px),
		TimeInForce: "GTD",
		PostOnly:    false,
	}
}

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
