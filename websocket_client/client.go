package wssClient

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"nhooyr.io/websocket"
)

const baseURL string = "wss://ws2.etherealtest.net/v1/stream"

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

type EventType int

const (
	EventUnknown EventType = iota
	EventL2Book
	EventMarketPrice
	EventTradeFill
	EventSubaccountLiquidation
	EventOrderUpdate
	EventOrderFill
	EventTokenTransfer
)

var eventTypeMap = map[string]EventType{
	"L2Book":                EventL2Book,
	"MarketPrice":           EventMarketPrice,
	"TradeFill":             EventTradeFill,
	"SubaccountLiquidation": EventSubaccountLiquidation,
	"OrderUpdate":           EventOrderUpdate,
	"OrderFill":             EventOrderFill,
	"TokenTransfer":         EventTokenTransfer,
}

type Client struct {
	Con                *websocket.Conn
	bookHandler        func(*L2Book)      // non-array
	priceHandler       func(*MarketPrice) // non-array
	tradeFillHandler   func(*TradeFillEvent)
	liquidationHandler func(*SubaccountLiquidationEvent)
	orderUpdateHandler func(*OrderUpdateEvent)
	orderFillHandler   func(*OrderFillEvent)
	transferHandler    func(*Transfer) // non-array
}

func NewClient(ctx context.Context) *Client {
	c, _, err := websocket.Dial(ctx, baseURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &Client{Con: c}
}

func marshalSubscribe[T eventData](data T) (b []byte, err error) {
	req := &SubIntent[T]{I: sub, D: data}
	return json.Marshal(req)
}

func marshalUnsubscribe[T eventData](data T) (b []byte, err error) {
	req := &SubIntent[T]{I: unsub, D: data}
	return json.Marshal(req)
}

func marshalToValueCallback[T proto.Message](data []byte, pb T, cb func(T)) (err error) {
	if err := protojson.Unmarshal(data, pb); err != nil {
		return err
	}
	cb(pb)
	return
}

// func marshalToArrayCallback[T proto.Message](data []byte, pb EventMessageArray, parse func() T, cb func(T)) (err error) {
// 	if err := protojson.Unmarshal(data, &pb); err != nil {
// 		return err
// 	}
// 	for _, eventMsg := range pb.Data {
// 		a := parse()
// 		a = *eventMsg
// 	}
// 	parse()
// 	cb(pb)
// 	return
// }

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

func (c *Client) OnBook(callback func(*L2Book)) {
	c.bookHandler = callback
}

func (c *Client) OnPrice(callback func(*MarketPrice)) {
	c.priceHandler = callback
}

func (c *Client) OnTradeFill(callback func(*TradeFillEvent)) {
	c.tradeFillHandler = callback
}

func (c *Client) OnLiquidation(callback func(*SubaccountLiquidationEvent)) {
	c.liquidationHandler = callback
}

func (c *Client) OnOrderUpdate(callback func(*OrderUpdateEvent)) {
	c.orderUpdateHandler = callback
}

func (c *Client) OnOrderFill(callback func(*OrderFillEvent)) {
	c.orderFillHandler = callback
}

func (c *Client) OnTransfer(callback func(*Transfer)) {
	c.transferHandler = callback
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

		var msg wssMsg
		if err := json.Unmarshal(data, &msg); err != nil {
			cancel(err)
			return err
		}

		var event EventType = EventUnknown
		var ok bool
		if event, ok = eventTypeMap[msg.Event]; !ok {
			event = EventUnknown
		}

		switch event {
		case EventL2Book:
			var diff L2Book
			if err := marshalToValueCallback(msg.Data, &diff, c.bookHandler); err != nil {
				cancel(err)
				return err
			}

		case EventMarketPrice:
			var mp MarketPrice
			if err := marshalToValueCallback(msg.Data, &mp, c.priceHandler); err != nil {
				cancel(err)
				return err
			}

		case EventSubaccountLiquidation:
			var lq SubaccountLiquidationEvent
			if err := marshalToValueCallback(msg.Data, &lq, c.liquidationHandler); err != nil {
				fmt.Println(string(data))
				cancel(err)
				return err
			}

		case EventOrderFill:
			var ou OrderFillEvent
			if err := marshalToValueCallback(data, &ou, c.orderFillHandler); err != nil {
				cancel(err)
				return err
			}

		case EventOrderUpdate:
			var ou OrderUpdateEvent
			if err := marshalToValueCallback(data, &ou, c.orderUpdateHandler); err != nil {
				cancel(err)
				return err
			}

		case EventTradeFill:
			var tf TradeFillEvent
			if err := marshalToValueCallback(data, &tf, c.tradeFillHandler); err != nil {
				cancel(err)
				return err
			}

		case EventTokenTransfer:
			var t Transfer
			if err := marshalToValueCallback(data, &t, c.transferHandler); err != nil {
				fmt.Println(string(data))
				cancel(err)
				return err
			}

		default:
			fmt.Printf("unknown event, raw: %s\n", string(data))
		}
	}
}

func (c *Client) Close() {
	c.Con.Close(websocket.StatusNormalClosure, "<3")
}
