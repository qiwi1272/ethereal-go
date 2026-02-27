package restClient

import (
	"context"
	"encoding/json"
	"sync"
)

func (o *Order) Send(ctx context.Context, client *Client) (OrderCreated, error) {
	var err error
	var created OrderCreated

	o.build(client)

	sig, err := Sign(o, "TradeOrder", client)
	if err != nil {
		return created, err
	}

	resp, err := client.do(ctx, "POST", "/v1/order", SignedGenericMessage{
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

func (o *CancelOrder) Send(ctx context.Context, client *Client) ([]OrderCancelled, error) {
	var cancelled Response[[]OrderCancelled]

	o.build(client)

	sig, err := Sign(o, "CancelOrder", client)
	if err != nil {
		return cancelled.Data, err
	}

	resp, err := client.do(ctx, "POST", "/v1/order/cancel", SignedGenericMessage{
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

type BatchResponse interface {
	OrderCreated | OrderCancelled
}

type intent string

const (
	Create intent = "TradeOrder"
	Cancel intent = "CancelOrder"
)

var intentMap = map[intent]string{
	Create: "/v1/order",
	Cancel: "/v1/order/cancel",
}

func SendBatch[ResponseType BatchResponse](
	ctx context.Context,
	cl *Client,
	intent intent,
	payload []Signable,
) ([]ResponseType, error) {

	batchSize := len(payload)
	var wg sync.WaitGroup

	wg.Add(batchSize)
	receipts := make([]ResponseType, batchSize)
	errCh := make(chan error, batchSize)

	for i, order := range payload {
		go func() {
			defer wg.Done()
			order.build(cl)
			sig, err := Sign(order, string(Create), cl)
			if err != nil {
				errCh <- err
				return
			}
			resp, err := cl.do(ctx, "POST", intentMap[intent], SignedGenericMessage{
				Data:      order,
				Signature: sig,
			})
			if err != nil {
				errCh <- err
				return
			}
			if err := json.Unmarshal(resp, &receipts[i]); err != nil {
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

	return receipts, nil
}
