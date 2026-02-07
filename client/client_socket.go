package ethereal

import (
	"fmt"
	"time"

	sio "github.com/karagenc/socket.io-go"
	eio "github.com/karagenc/socket.io-go/engine.io"
)

type WebsocketClient struct {
	manager *sio.Manager
	Socket  sio.ClientSocket
}

func NewWebSocketClient() *WebsocketClient {
	baseURL := "wss://ws.ethereal.trade/socket.io/"
	retryDelay := time.Minute
	config := &sio.ManagerConfig{
		EIO: eio.ClientConfig{
			Transports: []string{"websocket"},
		},
		ReconnectionDelay:    &retryDelay,
		ReconnectionAttempts: 11,
	}

	manager := sio.NewManager(baseURL, config)
	socket := manager.Socket("/v1/stream", nil)

	wsClient := &WebsocketClient{
		manager: manager,
		Socket:  socket,
	}

	// native events + open connection
	go func(ws *WebsocketClient) {
		ws.Socket.OnConnect(func() {
			fmt.Println("connected via socket to ethereal")
		})

		ws.Socket.OnDisconnect(func(reason sio.Reason) {
			fmt.Println("ethereal socket disconnected: ", reason)
		})

		ws.manager.OnError(func(err error) {
			fmt.Printf("ethereal socket manager error: %v\n", err)
		})

		ws.Socket.Connect()
	}(wsClient)

	return &WebsocketClient{
		manager: manager,
		Socket:  socket,
	}
}

func (ws *WebsocketClient) SubscribeToBook(productId string) {
	req := map[string]any{
		"type":      "BookDepth",
		"productId": productId,
	}
	ws.Socket.Emit("subscribe", req)
}

func (ws *WebsocketClient) OnBookDepth(handler func(BookDepthStream)) {
	ws.Socket.OnEvent("BookDepth", handler)
}

func (ws *WebsocketClient) SubscribeToPrice(productId string) {
	req := map[string]any{
		"type":      "MarketPrice",
		"productId": productId,
	}
	ws.Socket.Emit("subscribe", req)
}

func (ws *WebsocketClient) OnPrice(handler func(MarketPriceStream)) {
	ws.Socket.OnEvent("MarketPrice", handler)
}

func (ws *WebsocketClient) SubscribeToFill(s *Subaccount) {
	req := map[string]any{
		"type":         "OrderFill",
		"subaccountId": s.Id,
	}
	ws.Socket.Emit("subscribe", req)
}

func (ws *WebsocketClient) OnFill(handler func(OrderFillStream)) {
	ws.Socket.OnEvent("OrderFill", handler)
}

func (ws *WebsocketClient) SubscribeToOrder(s *Subaccount) {
	req := map[string]any{
		"type":      "OrderUpdate",
		"productId": s.Id,
	}
	ws.Socket.Emit("subscribe", req)
}

func (ws *WebsocketClient) OnOrder(handler func(OrderStream)) {
	ws.Socket.OnEvent("OrderUpdate", handler)
}

func (ws *WebsocketClient) OnDisconnect(handler func(sio.Reason)) {
	ws.Socket.OnDisconnect(handler)
}
