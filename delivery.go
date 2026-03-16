package ethereal

import (
	"context"
	"encoding/json"
)

type OrderClient interface {
	SubaccountHolder
	Do(ctx context.Context, method, path string, body any) ([]byte, error)
}

type Response[T any] struct {
	Data T `json:"data"`
}

type OrderCreated struct { // TODO: missed one
	Id     string `json:"id"`
	Cloid  string `json:"clientOrderId"`
	Filled string `json:"filled"`
	Result string `json:"result"`
}

func (o *Order) Send(ctx context.Context, cl OrderClient, signer *Signer) (OrderCreated, error) {
	var err error
	var created OrderCreated

	o.Build(cl)

	sig, err := Sign(o, "TradeOrder", signer)
	if err != nil {
		return created, err
	}

	resp, err := cl.Do(ctx, "POST", "/v1/order", SignedMessage[*Order]{
		Data:      o,
		Signature: sig,
	})
	if err != nil {
		return created, err
	}

	if err := json.Unmarshal(resp, &created); err != nil {
		return created, err
	}

	return created, nil
}

func (o *OrderCancel) Send(ctx context.Context, cl OrderClient, signer *Signer) ([]*OrderCancelled, error) {
	var cancelled Response[[]*OrderCancelled]

	o.Build(cl)

	sig, err := Sign(o, "CancelOrder", signer)
	if err != nil {
		return cancelled.Data, err
	}

	resp, err := cl.Do(ctx, "POST", "/v1/order/cancel", SignedMessage[*OrderCancel]{
		Data:      o,
		Signature: sig,
	})
	if err != nil {
		return cancelled.Data, err
	}

	if err := json.Unmarshal(resp, &cancelled); err != nil {
		return cancelled.Data, err
	}

	return cancelled.Data, nil
}
