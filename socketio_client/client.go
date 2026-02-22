package socketio_client

import (
	"encoding/json"
	"fmt"
	"time"

	sio "github.com/karagenc/socket.io-go"
	eio "github.com/karagenc/socket.io-go/engine.io"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/qiwi1272/ethereal-go"
	etherealpb "github.com/qiwi1272/ethereal-go/_pb"
)

type SocketIOClient struct {
	manager *sio.Manager
	Socket  sio.ClientSocket
	pm      *protojson.UnmarshalOptions
}

func NewSocketIOClient() *SocketIOClient {
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

	wsClient := &SocketIOClient{
		manager: manager,
		Socket:  socket,
		pm:      &protojson.UnmarshalOptions{DiscardUnknown: true},
	}

	// native events + open connection
	go func(ws *SocketIOClient) {
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

	return &SocketIOClient{
		manager: manager,
		Socket:  socket,
	}
}

// avoid writing our own stream parser and use rawMessage

func (ws *SocketIOClient) SubscribeToBook(productId string) {
	req := map[string]any{
		"type":      "BookDepth",
		"productId": productId,
	}
	ws.Socket.Emit("subscribe", req)
}

const pid_prefix_len = len("{\"productId\":")
const ts_prefix_len = len(",\"timestamp\":")
const prev_ts_prefix_len = len("\"previousTimestamp\":")
const asks_prefix_len = len("\"asks\":")
const bids_prefix_len = len(",\"bids\":")

func (ws *SocketIOClient) OnBookDepth(handler func(*etherealpb.BookDiff)) {
	ws.Socket.OnEvent("BookDepth", func(bytes json.RawMessage) {

		diff := &etherealpb.BookDiff{}

		var next int
		var err error

		bytes = bytes[pid_prefix_len:] // consume {"productId":
		if next, diff.ProductId, err = etherealpb.ReadStringAt(bytes, 0); err != nil {
			panic(err)
		}

		bytes = bytes[next+ts_prefix_len:] // consume ,"timestamp":
		if next, diff.Timestamp, err = etherealpb.ReadInt64At(bytes, 0, ','); err != nil {
			panic(err)
		}

		bytes = bytes[next+prev_ts_prefix_len:] // consume "previousTimestamp":
		if next, diff.PreviousTimestamp, err = etherealpb.ReadInt64At(bytes, 0, ','); err != nil {
			panic(err)
		}

		bytes = bytes[next+asks_prefix_len:] // consume "asks":
		if next, err = diff.DecodeDiffSideMsg(bytes, true); err != nil {
			panic(err)
		}

		bytes = bytes[next+bids_prefix_len:] // consume ,"bids":
		if next, err = diff.DecodeDiffSideMsg(bytes, false); err != nil {
			panic(err)
		}

		handler(diff)
	})
}

func (ws *SocketIOClient) SubscribeToPrice(productId string) {
	req := map[string]any{
		"type":      "MarketPrice",
		"productId": productId,
	}
	ws.Socket.Emit("subscribe", req)
}

func (ws *SocketIOClient) OnPrice(handler func(MarketPriceStream)) {
	ws.Socket.OnEvent("MarketPrice", handler)
}

func (ws *SocketIOClient) SubscribeToFill(s *ethereal.Subaccount) {
	req := map[string]any{
		"type":         "OrderFill",
		"subaccountId": s.Id,
	}
	ws.Socket.Emit("subscribe", req)
}

func (ws *SocketIOClient) OnFill(handler func(OrderFillStream)) {
	ws.Socket.OnEvent("OrderFill", handler)
}

func (ws *SocketIOClient) SubscribeToOrder(s *ethereal.Subaccount) {
	req := map[string]any{
		"type":      "OrderUpdate",
		"productId": s.Id,
	}
	ws.Socket.Emit("subscribe", req)
}

func (ws *SocketIOClient) OnOrder(handler func(OrderStream)) {
	ws.Socket.OnEvent("OrderUpdate", handler)
}

func (ws *SocketIOClient) OnDisconnect(handler func(sio.Reason)) {
	ws.Socket.OnDisconnect(handler)
}
