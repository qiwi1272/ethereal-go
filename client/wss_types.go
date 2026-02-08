package ethereal

import "encoding/json"

type requestType string

const SubscribeRequest requestType = "subscribe"
const UnsubscribeRequest requestType = "unsubscribe"

type WssMessage[T WssData] struct {
	Event requestType `json:"event"` // subscribe
	Data  T
}

type eventType string

const BookEvent eventType = "L2Book"
const MarketPriceEvent eventType = "MarketPrice"

type WssData struct {
	Type eventType `json:"type"`
}

type WssSymbolData struct {
	WssData
	Symbol string `json:"symbol"`
}

type BookDepthL2WssStream struct {
	Event             string          `json:"e"`
	Symbol            string          `json:"s"`
	Timestamp         float64         `json:"t"`
	PreviousTimestamp float64         `json:"pt"`
	Asks              [][]json.Number `json:"a"`
	Bids              [][]json.Number `json:"b"`
}

type MarketPriceWssStream struct {
	MarketPrice eventType `json:"e"`
	Symbol      string    `json:"s"`
	Timestamp   float64   `json:"t"`
	BidPx       string    `json:"bidPx"`
	AskPx       string    `json:"askPx"`
	MarkPx      string    `json:"markPx"`
	Mark24hPx   string    `json:"mark24hPx"` // Price24hAgo
}

type ErrorWss struct {
	Code string `json:"code"`
}

// type MarketPriceStream struct {
// 	ProductID    string `json:"productId"`
// 	BestBidPrice string `json:"bestBidPrice"`
// 	BestAskPrice string `json:"bestAskPrice"`
// 	OraclePrice  string `json:"oraclePrice"`
// 	Price24hAgo  string `json:"price24hAgo"`
// }

// type OrderFillObject struct {
// 	Id            string    `json:"id"`      // uuid
// 	OrderId       string    `json:"orderId"` // uuid
// 	ClientOrderID string    `json:"bestAskPrice"`
// 	Price         string    `json:"oraclePrice"`
// 	Filled        string    `json:"filled"`
// 	Type          OrderType `json:"type"`
// 	Side          OrderSide `json:"side"`
// 	ReduceOnly    bool      `json:"reduceOnly"`
// 	FeeUsd        string    `json:"feeUsd"`
// 	IsMaker       bool      `json:"isMaker"`
// 	ProductId     string    `json:"productId"`    // uuid
// 	SubaccountId  string    `json:"subaccountId"` // uuid
// 	CreatedAt     int64     `json:"createdAt"`
// }

// type OrderFillStream struct {
// 	Data []OrderFillObject `json:"data"`
// }

// type OrderStreamObject struct {
// 	Id                   string      `json:"id"` // uuid
// 	ClientOrderID        string      `json:"clientOrderId,omitempty"`
// 	AvailableQuantity    string      `json:"availableQuantity"`
// 	Quantity             string      `json:"quantity"`
// 	Side                 OrderSide   `json:"side"`         // 0 BUY, 1 SELL
// 	ProductId            string      `json:"productId"`    // uuid
// 	SubaccountId         string      `json:"subaccountId"` // uuid
// 	Status               string      `json:"status"`       // TODO: eunm status
// 	ReduceOnly           bool        `json:"reduceOnly"`
// 	Close                bool        `json:"close"`
// 	UpdatedAt            int64       `json:"updatedAt"` // epoch
// 	CreatedAt            int64       `json:"createdAt"` // epoch
// 	Sender               string      `json:"sender"`
// 	Price                string      `json:"price"`
// 	Filled               string      `json:"filled"`
// 	StopPrice            string      `json:"stopPrice,omitempty"`
// 	StopType             string      `json:"stopType,omitempty"`
// 	StopPriceType        string      `json:"stopPriceType,omitempty"`
// 	TimeInForce          TimeInForce `json:"timeInForce"`
// 	ExpiresAt            int64       `json:"expiresAt,omitempty"`
// 	PostOnly             bool        `json:"postOnly"`
// 	GroupID              string      `json:"groupId,omitempty"` // UUID
// 	GroupContingencyType int         `json:"groupContingencyType,omitempty"`
// }

// type OrderStream struct {
// 	Data []OrderStreamObject `json:"data"`
// }
