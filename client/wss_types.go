package ethereal

import "encoding/json"

type BookDepthStream struct {
	Bids              [][]json.Number `json:"bids"`
	Asks              [][]json.Number `json:"asks"`
	ProductID         string          `json:"productId"`
	Timestamp         float64         `json:"timestamp"`
	PreviousTimestamp float64         `json:"previousTimestamp"`
}

type MarketPriceStream struct {
	ProductID    string `json:"productId"`
	BestBidPrice string `json:"bestBidPrice"`
	BestAskPrice string `json:"bestAskPrice"`
	OraclePrice  string `json:"oraclePrice"`
	Price24hAgo  string `json:"price24hAgo"`
}

type OrderFillStream struct {
	Data []struct {
		Id            string    `json:"id"`      // uuid
		OrderId       string    `json:"orderId"` // uuid
		ClientOrderID string    `json:"bestAskPrice"`
		Price         string    `json:"oraclePrice"`
		Filled        string    `json:"filled"`
		Type          OrderType `json:"type"`
		Side          OrderSide `json:"side"`
		ReduceOnly    bool      `json:"reduceOnly"`
		FeeUsd        string    `json:"feeUsd"`
		IsMaker       bool      `json:"isMaker"`
		ProductId     string    `json:"productId"`    // uuid
		SubaccountId  string    `json:"subaccountId"` // uuid
		CreatedAt     int64     `json:"createdAt"`
	} `json:"data"`
}

type OrderStream struct {
	Data []struct {
		Id                   string      `json:"id"` // uuid
		ClientOrderID        string      `json:"clientOrderId,omitempty"`
		AvailableQuantity    string      `json:"availableQuantity"`
		Quantity             string      `json:"quantity"`
		Side                 OrderSide   `json:"side"`         // 0 BUY, 1 SELL
		ProductId            string      `json:"productId"`    // uuid
		SubaccountId         string      `json:"subaccountId"` // uuid
		Status               string      `json:"status"`       // TODO: eunm status
		ReduceOnly           bool        `json:"reduceOnly"`
		Close                bool        `json:"close"`
		UpdatedAt            int64       `json:"updatedAt"` // epoch
		CreatedAt            int64       `json:"createdAt"` // epoch
		Sender               string      `json:"sender"`
		Price                string      `json:"price"`
		Filled               string      `json:"filled"`
		StopPrice            string      `json:"stopPrice,omitempty"`
		StopType             string      `json:"stopType,omitempty"`
		StopPriceType        string      `json:"stopPriceType,omitempty"`
		TimeInForce          TimeInForce `json:"timeInForce"`
		ExpiresAt            int64       `json:"expiresAt,omitempty"`
		PostOnly             bool        `json:"postOnly"`
		GroupID              string      `json:"groupId,omitempty"` // UUID
		GroupContingencyType int         `json:"groupContingencyType,omitempty"`
	} `json:"data"`
}
