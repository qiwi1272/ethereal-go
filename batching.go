package etherealRest

import (
	"context"
	"encoding/json"
	"sync"
)

// -------- BEGIN BATCH ORDER -------- //

type Batchable interface {
	Signable
	*Order | *OrderCancel
}

type BatchResponseType interface {
	*OrderCreated | *OrderCancelled
}

type BatchOrder[T BatchResponseType] struct {
	Payload []Signable
	resp    []T
}

func NewBatch[R BatchResponseType](
	items []Signable,
) *BatchOrder[R] {
	return &BatchOrder[R]{
		Payload: items,
		resp:    make([]R, len(items)),
	}
}

func NewOrderBatch(orders []*Order) *BatchOrder[*OrderCreated] {
	var payload = make([]Signable, len(orders))
	for i, o := range orders {
		payload[i] = o
	}
	return &BatchOrder[*OrderCreated]{
		Payload: payload,
		resp:    make([]*OrderCreated, len(orders)),
	}
}

func NewCancelBatch(created []*OrderCreated) *OrderCancel {
	ids := make([]string, len(created))
	for i, o := range created {
		ids[i] = o.Id
	}
	return &OrderCancel{
		OrderIDs: ids,
	}
}

// -------- END BATCH ORDER -------- //

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
	cl *Client,
	intent batchIntent,
	signer *Signer,
) ([]BatchResponseType, error) {

	batchSize := len(b.Payload)
	var wg sync.WaitGroup

	wg.Add(batchSize)
	errCh := make(chan error, batchSize)

	for i, order := range b.Payload {
		go func() {
			defer wg.Done()
			order.Build(cl)
			sig, err := Sign(order, string(intent), signer)
			if err != nil {
				errCh <- err
				return
			}
			path := batchIntentMap[intent]
			resp, err := cl.Do(ctx, "POST", path, SignedMessage[Signable]{
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
