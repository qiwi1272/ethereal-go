package restClient

import "github.com/qiwi1272/ethereal-go"

type Product struct {
	ID                     string              `json:"id"`
	Ticker                 string              `json:"ticker"`
	DisplayTicker          string              `json:"displayTicker"`
	EngineType             ethereal.EngineType `json:"engineType"`
	OnchainID              int64               `json:"onchainId"`
	LotSize                string              `json:"lotSize"`
	TickSize               string              `json:"tickSize"`
	MakerFee               string              `json:"makerFee"`
	TakerFee               string              `json:"takerFee"`
	MaxQuantity            string              `json:"maxQuantity"`
	MinQuantity            string              `json:"minQuantity"`
	Volume24h              string              `json:"volume24h"`
	FundingRate1h          string              `json:"fundingRate1h"`
	MaxOpenInterestUsd     string              `json:"maxOpenInterestUsd"`
	MaxPositionNotionalUsd string              `json:"maxPositionNotionalUsd"`
}

type Order struct {
	Subaccount           string               `json:"subaccount"`
	Sender               string               `json:"sender"`
	Nonce                string               `json:"nonce"` // string of nanoseconds
	Type                 ethereal.OrderType   `json:"type"`  // LIMIT or MARKET
	Quantity             string               `json:"quantity"`
	Side                 ethereal.OrderSide   `json:"side"` // 0 BUY, 1 SELL
	OnchainID            int64                `json:"onchainId"`
	EngineType           ethereal.EngineType  `json:"engineType"`
	ClientOrderID        string               `json:"clientOrderId,omitempty"`
	ReduceOnly           bool                 `json:"reduceOnly"`
	Close                bool                 `json:"close,omitempty"`
	StopPrice            int64                `json:"stopPrice,omitempty"`
	StopType             int64                `json:"stopType,omitempty"`
	SignedAt             int64                `json:"signedAt"` // seconds since epoch
	ExpiresAt            int64                `json:"expiresAt,omitempty"`
	GroupID              string               `json:"groupId,omitempty"` // UUID
	GroupContingencyType int                  `json:"groupContingencyType,omitempty"`
	Price                string               `json:"price"`
	TimeInForce          ethereal.TimeInForce `json:"timeInForce"`
	PostOnly             bool                 `json:"postOnly"`
}

type OrderCreated struct {
	Id     string `json:"id"`
	Cloid  string `json:"clientOrderId"`
	Filled string `json:"filled"`
	Result string `json:"result"`
}

type CancelOrder struct {
	Sender     string   `json:"sender"`
	Subaccount string   `json:"subaccount"`
	Nonce      string   `json:"nonce"`
	OrderIDs   []string `json:"orderIds"`
	Cloids     []string `json:"clientOrderIds"`
}

type OrderCancelled struct {
	Id     string `json:"id"`
	Cloid  string `json:"clientOrderId"`
	Result string `json:"result"`
}

type Position struct {
	Id                    string             `json:"id"`
	Cost                  string             `json:"cost"`
	Size                  string             `json:"size"`
	FundingUsd            string             `json:"fundingUsd"`
	FundingAccruedUsd     string             `json:"fundingAccruedUsd"`
	FeesAccruedUsd        string             `json:"feesAccruedUsd"`
	RealizedPnl           string             `json:"realizedPnl"`
	TotalIncreaseNotional string             `json:"totalIncreaseNotional"`
	TotalIncreaseQuantity string             `json:"totalIncreaseQuantity"`
	TotalDecreaseNotional string             `json:"totalDecreaseNotional"`
	TotalDecreaseQuantity string             `json:"totalDecreaseQuantity"`
	Side                  ethereal.OrderSide `json:"side"`
	ProductId             string             `json:"productId"`
	UpdatedAt             uint64             `json:"updatedAt"`
	CreatedAt             uint64             `json:"createdAt"`
	IsLiquidated          bool               `json:"isLiquidated"`
	LiquidationPrice      string             `json:"liquidationPrice"`
}

type AccountBalance struct {
	SubaccountId string `json:"subaccountId"`
	TokenId      string `json:"tokenId"`
	TokenAddress string `json:"tokenAddress"`
	TokenName    string `json:"tokenName"`
	Amount       string `json:"amount"`
	Available    string `json:"available"`
	TotalUsed    string `json:"totalUsed"`
	UpdatedAt    uint64 `json:"updatedAt"`
}
