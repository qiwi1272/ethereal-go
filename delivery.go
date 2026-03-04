package ethereal

import (
	"context"
	"encoding/json"
	"sync"
)

type BatchOrderClient interface {
	Signer
	SubaccountHolder
	Do(ctx context.Context, method, path string, body any) ([]byte, error)
}

type Response[T any] struct {
	Data T `json:"data"`
}

func (o *Order) Send(ctx context.Context, cl BatchOrderClient) (OrderCreated, error) {
	var err error
	var created OrderCreated

	o.build(cl)

	sig, err := Sign(o, "TradeOrder", cl)
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

func (o *OrderCancel) Send(ctx context.Context, cl BatchOrderClient) ([]OrderCancelled, error) {
	var cancelled Response[[]OrderCancelled]

	o.build(cl)

	sig, err := Sign(o, "CancelOrder", cl)
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

type batchIntent string

const (
	Create batchIntent = "TradeOrder"
	Cancel batchIntent = "CancelOrder"
)

var batchIntentMap = map[batchIntent]string{
	Create: "/v1/order",
	Cancel: "/v1/order/cancel",
}

func (b *BatchOrder[BatchResponseType]) SendBatch(
	ctx context.Context,
	cl *BatchOrderClient,
	intent batchIntent,
) ([]BatchResponseType, error) {

	batchSize := len(b.Payload)
	var wg sync.WaitGroup

	wg.Add(batchSize)
	errCh := make(chan error, batchSize)

	for i, o := range b.Payload {
		go func() {
			defer wg.Done()
			order := *o
			Cl := *cl
			order.build(Cl)
			sig, err := Sign(order, string(Create), Cl)
			if err != nil {
				errCh <- err
				return
			}
			path := batchIntentMap[intent]
			resp, err := Cl.Do(ctx, "POST", path, SignedMessage[Signable]{
				Data:      order,
				Signature: sig,
			})
			if err != nil {
				errCh <- err
				return
			}
			if err := json.Unmarshal(resp, &b.resp[i]); err != nil {
				errCh <- err
				return
			}
		}()
	}

	wg.Wait()
	close(errCh)

	for err := range errCh { // first error is fine, there should be none
		if err != nil {
			return nil, err
		}
	}

	return b.resp, nil
}
