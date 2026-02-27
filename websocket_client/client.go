package wssClient

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"google.golang.org/protobuf/encoding/protojson"
	"nhooyr.io/websocket"
)

const baseURL string = "wss://ws2.etherealtest.net/v1/stream"

type BookHandler = func(*L2Book)
type PriceHandler = func(*MarketPrice)

type Client struct {
	Con *websocket.Conn
	bh  BookHandler
	ph  PriceHandler
}

func NewClient(ctx context.Context) *Client {
	c, _, err := websocket.Dial(ctx, baseURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &Client{Con: c}
}

type Intent string

const (
	sub   Intent = "subscribe"
	unsub Intent = "unsubscribe"
)

type eventData interface{}

type SymbolEvent struct {
	eventData
	T string `json:"type"`
	S string `json:"symbol"`
}

type SubaccountEvent struct {
	eventData
	T string `json:"type"`
	S string `json:"subaccountId"`
}

type SubIntent[T eventData] struct {
	I Intent    `json:"event"`
	D eventData `json:"data"`
}

func marshalSubscribe[T eventData](data T) (b []byte, err error) {
	req := &SubIntent[T]{I: sub, D: data}
	return json.Marshal(req)
}

func marshalUnsubscribe[T eventData](data T) (b []byte, err error) {
	req := &SubIntent[T]{I: unsub, D: data}
	return json.Marshal(req)
}

func (c *Client) req(ctx context.Context, payload []byte) (err error) {
	return c.Con.Write(ctx, websocket.MessageBinary, payload)
}

func (c *Client) SubscribeBook(ctx context.Context, symbol string) (err error) {
	var b []byte
	if b, err = marshalSubscribe(&SymbolEvent{
		T: "L2Book",
		S: symbol,
	}); err != nil {
		return err
	}
	return c.req(ctx, b)
}

func (c *Client) UnsubscribeBook(ctx context.Context, symbol string) (err error) {
	var b []byte
	if b, err = marshalUnsubscribe(&SymbolEvent{
		T: "L2Book",
		S: symbol,
	}); err != nil {
		return err
	}
	return c.req(ctx, b)
}

func (c *Client) SubscribeMarketPrice(ctx context.Context, symbol string) (err error) {
	var b []byte
	if b, err = marshalSubscribe(&SymbolEvent{
		T: "MarketPrice",
		S: symbol,
	}); err != nil {
		return err
	}
	return c.req(ctx, b)
}

func (c *Client) UnsubscribeMarketPrice(ctx context.Context, symbol string) (err error) {
	var b []byte
	if b, err = marshalUnsubscribe(&SymbolEvent{
		T: "MarketPrice",
		S: symbol,
	}); err != nil {
		return err
	}
	return c.req(ctx, b)
}

func (c *Client) SubscribeFill(ctx context.Context, symbol string) (err error) {
	var b []byte
	if b, err = marshalSubscribe(&SymbolEvent{
		T: "TradeFill",
		S: symbol,
	}); err != nil {
		return err
	}
	return c.req(ctx, b)
}

func (c *Client) UnsubscribeFill(ctx context.Context, symbol string) (err error) {
	var b []byte
	if b, err = marshalUnsubscribe(&SymbolEvent{
		T: "TradeFill",
		S: symbol,
	}); err != nil {
		return err
	}
	return c.req(ctx, b)
}

func (c *Client) SubscribeLiquidation(ctx context.Context, subaccountUuid string) (err error) {
	var b []byte
	if b, err = marshalSubscribe(&SubaccountEvent{
		T: "SubaccountLiquidation",
		S: subaccountUuid,
	}); err != nil {
		return err
	}
	return c.req(ctx, b)
}

func (c *Client) UnsubscribeLiquidation(ctx context.Context, subaccountUuid string) (err error) {
	var b []byte
	if b, err = marshalUnsubscribe(&SubaccountEvent{
		T: "SubaccountLiquidation",
		S: subaccountUuid,
	}); err != nil {
		return err
	}
	return c.req(ctx, b)
}

func (c *Client) SubscribeOrderUpdate(ctx context.Context, subaccountUuid string) (err error) {
	var b []byte
	if b, err = marshalSubscribe(&SubaccountEvent{
		T: "OrderUpdate",
		S: subaccountUuid,
	}); err != nil {
		return err
	}
	return c.req(ctx, b)
}

func (c *Client) UnsubscribeOrderUpdate(ctx context.Context, subaccountUuid string) (err error) {
	var b []byte
	if b, err = marshalUnsubscribe(&SubaccountEvent{
		T: "OrderUpdate",
		S: subaccountUuid,
	}); err != nil {
		return err
	}
	return c.req(ctx, b)
}

func (c *Client) SubscribeOrderFill(ctx context.Context, subaccountUuid string) (err error) {
	var b []byte
	if b, err = marshalSubscribe(&SubaccountEvent{
		T: "OrderFill",
		S: subaccountUuid,
	}); err != nil {
		return err
	}
	return c.req(ctx, b)
}

func (c *Client) UnsubscribeOrderFill(ctx context.Context, subaccountUuid string) (err error) {
	var b []byte
	if b, err = marshalUnsubscribe(&SubaccountEvent{
		T: "OrderFill",
		S: subaccountUuid,
	}); err != nil {
		return err
	}
	return c.req(ctx, b)
}

func (c *Client) SubscribeTokenTransfer(ctx context.Context, subaccountUuid string) (err error) {
	var b []byte
	if b, err = marshalSubscribe(&SubaccountEvent{
		T: "TokenTransfer",
		S: subaccountUuid,
	}); err != nil {
		return err
	}
	return c.req(ctx, b)
}

func (c *Client) UnsubscribeTokenTransfer(ctx context.Context, subaccountUuid string) (err error) {
	var b []byte
	if b, err = marshalUnsubscribe(&SubaccountEvent{
		T: "TokenTransfer",
		S: subaccountUuid,
	}); err != nil {
		return err
	}
	return c.req(ctx, b)
}

func (c *Client) OnBook(callback BookHandler) {
	c.bh = callback
}

func (c *Client) OnPrice(callback PriceHandler) {
	c.ph = callback
}

type wssMsg struct {
	Event string          `json:"e"`
	Ts    int64           `json:"t"`
	Data  json.RawMessage `json:"data"`
}

func (c *Client) Listen(parent context.Context) error {
	ctx, cancel := context.WithCancelCause(parent)
	defer cancel(nil)
	defer c.Close()

	for {
		_, data, err := c.Con.Read(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return context.Cause(ctx)
			}
			cancel(err)
			return err
		}

		fmt.Println(string(data))

		var msg wssMsg
		if err := json.Unmarshal(data, &msg); err != nil {
			cancel(err)
			return err
		}

		switch msg.Event {
		case "L2Book":
			var diff L2Book
			if err := protojson.Unmarshal(msg.Data, &diff); err != nil {
				cancel(err)
				return err
			}
			c.bh(&diff)

		case "MarketPrice":
			var mp MarketPrice
			if err := protojson.Unmarshal(msg.Data, &mp); err != nil {
				cancel(err)
				return err
			}
			c.ph(&mp)

		default:
			fmt.Printf("unknown event, raw: %s\n", string(data))
		}
	}
}

func (c *Client) Close() {
	c.Con.Close(websocket.StatusNormalClosure, "<3")
}
