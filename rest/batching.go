package rest

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/qiwi1272/ethereal-go"
)

// -------- BEGIN BATCH ORDER -------- //

type Batchable interface {
	ethereal.Signable
	*ethereal.Order | *ethereal.OrderCancel
}

type BatchResponseType interface {
	*ethereal.OrderCreated | *ethereal.OrderCancelled
}

type BatchOrder[T BatchResponseType] struct {
	Payload []ethereal.Signable
	resp    []T
}

func NewBatch[R BatchResponseType](
	items []ethereal.Signable,
) *BatchOrder[R] {
	return &BatchOrder[R]{
		Payload: items,
		resp:    make([]R, len(items)),
	}
}

func NewOrderBatch(orders []*ethereal.Order) *BatchOrder[*ethereal.OrderCreated] {
	var payload = make([]ethereal.Signable, len(orders))
	for i, o := range orders {
		payload[i] = o
	}
	return &BatchOrder[*ethereal.OrderCreated]{
		Payload: payload,
		resp:    make([]*ethereal.OrderCreated, len(orders)),
	}
}

func NewCancelBatch(created []*ethereal.OrderCreated) *ethereal.OrderCancel {
	ids := make([]string, len(created))
	for i, o := range created {
		ids[i] = o.Id
	}
	return &ethereal.OrderCancel{
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
	signer *ethereal.Signer,
) ([]BatchResponseType, error) {

	batchSize := len(b.Payload)
	var wg sync.WaitGroup

	wg.Add(batchSize)
	errCh := make(chan error, batchSize)

	for i, order := range b.Payload {
		go func() {
			defer wg.Done()
			order.Build(cl)
			sig, err := ethereal.Sign(order, string(Create), signer)
			if err != nil {
				errCh <- err
				return
			}
			path := batchIntentMap[intent]
			resp, err := cl.Do(ctx, "POST", path, ethereal.SignedMessage[ethereal.Signable]{
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
