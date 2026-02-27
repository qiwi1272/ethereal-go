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

type SymbolData struct {
	Type   string `json:"type"`
	Symbol string `json:"symbol"`
}

type SymbolMessage struct {
	Event string     `json:"event"`
	Data  SymbolData `json:"data"`
}

func (c *Client) SubscribeBook(ctx context.Context, symbol string) error {
	var err error
	var payload []byte

	req := &SymbolMessage{
		Event: "subscribe",
		Data: SymbolData{
			Type:   "L2Book",
			Symbol: symbol,
		},
	}

	if payload, err = json.Marshal(req); err != nil {
		return err
	}

	if err = c.Con.Write(ctx, websocket.MessageBinary, payload); err != nil {
		return err
	}

	return nil
}

func (c *Client) SubscribePrice(ctx context.Context, symbol string) error {
	var err error
	var payload []byte

	req := &SymbolMessage{
		Event: "subscribe",
		Data: SymbolData{
			Type:   "MarketPrice",
			Symbol: symbol,
		},
	}

	if payload, err = json.Marshal(req); err != nil {
		return err
	}

	if err = c.Con.Write(ctx, websocket.MessageBinary, payload); err != nil {
		return err
	}

	return nil
}

func (c *Client) UnsubscribeBook(ctx context.Context, symbol string) error {
	var err error
	var payload []byte

	req := &SymbolMessage{
		Event: "unsubscribe",
		Data: SymbolData{
			Type:   "L2Book",
			Symbol: symbol,
		},
	}

	if payload, err = json.Marshal(req); err != nil {
		return err
	}

	if err = c.Con.Write(ctx, websocket.MessageBinary, payload); err != nil {
		return err
	}

	return nil
}

func (c *Client) UnsubscribePrice(ctx context.Context, symbol string) error {
	var err error
	var payload []byte

	req := &SymbolMessage{
		Event: "unsubscribe",
		Data: SymbolData{
			Type:   "MarketPrice",
			Symbol: symbol,
		},
	}

	if payload, err = json.Marshal(req); err != nil {
		return err
	}

	if err = c.Con.Write(ctx, websocket.MessageBinary, payload); err != nil {
		return err
	}

	return nil
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
