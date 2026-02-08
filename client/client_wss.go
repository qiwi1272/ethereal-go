package ethereal

import (
	"context"
	"encoding/json"
	"log"

	"nhooyr.io/websocket"
)

const baseURL string = "wss://ws.etherealtest.net/v1/stream"

type WebsocketClient struct {
	con *websocket.Conn
}

func NewWebSocketClient(ctx context.Context) *WebsocketClient {
	c, resp, err := websocket.Dial(ctx, baseURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(resp)

	return &WebsocketClient{con: c}
}

func createSubscribe[T eventData](data T) *WssMessage[T] {
	return &WssMessage[T]{
		Event: WssSubscribeEvent,
		Data:  data,
	}
}

func createUnubscribe[T eventData](data T) *WssMessage[T] {
	return &WssMessage[T]{
		Event: WssUnsubscribeEvent,
		Data:  data,
	}
}

func (c *WebsocketClient) SubscribeBook(ctx context.Context, symbol string) error {
	var err error
	var payload []byte

	req := createSubscribe(&WssSymbolData{
		Type:   BookEventType,
		Symbol: symbol,
	})

	if payload, err = json.Marshal(req); err != nil {
		return err
	}

	if err = c.con.Write(ctx, websocket.MessageBinary, payload); err != nil {
		return err
	}

	return nil
}

func (c *WebsocketClient) SubscribePrice(ctx context.Context, symbol string) error {
	var err error
	var payload []byte

	req := createSubscribe(&WssSymbolData{
		Type:   MarketPriceEventType,
		Symbol: symbol,
	})

	if payload, err = json.Marshal(req); err != nil {
		return err
	}

	if err = c.con.Write(ctx, websocket.MessageBinary, payload); err != nil {
		return err
	}

	return nil
}

func (c *WebsocketClient) UnsubscribeBook(ctx context.Context, symbol string) error {
	var err error
	var payload []byte

	req := createUnubscribe(&WssSymbolData{
		Type:   BookEventType,
		Symbol: symbol,
	})

	if payload, err = json.Marshal(req); err != nil {
		return err
	}

	if err = c.con.Write(ctx, websocket.MessageBinary, payload); err != nil {
		return err
	}

	return nil
}

func (c *WebsocketClient) UnsubscribePrice(ctx context.Context, symbol string) error {
	var err error
	var payload []byte

	req := createUnubscribe(&WssSymbolData{
		Type:   MarketPriceEventType,
		Symbol: symbol,
	})

	if payload, err = json.Marshal(req); err != nil {
		return err
	}

	if err = c.con.Write(ctx, websocket.MessageBinary, payload); err != nil {
		return err
	}

	return nil
}

func (c *WebsocketClient) Close() {
	c.con.Close(websocket.StatusNormalClosure, "<3")
}
