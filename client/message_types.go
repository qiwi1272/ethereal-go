package ethereal

type Product struct {
	ID                     string     `json:"id"`
	Ticker                 string     `json:"ticker"`
	DisplayTicker          string     `json:"displayTicker"`
	EngineType             EngineType `json:"engineType"`
	OnchainID              int64      `json:"onchainId"`
	LotSize                string     `json:"lotSize"`
	TickSize               string     `json:"tickSize"`
	MakerFee               string     `json:"makerFee"`
	TakerFee               string     `json:"takerFee"`
	MaxQuantity            string     `json:"maxQuantity"`
	MinQuantity            string     `json:"minQuantity"`
	Volume24h              string     `json:"volume24h"`
	FundingRate1h          string     `json:"fundingRate1h"`
	MaxOpenInterestUsd     string     `json:"maxOpenInterestUsd"`
	MaxPositionNotionalUsd string     `json:"maxPositionNotionalUsd"`
}

type TimeInForce string

const (
	TIF_GTD TimeInForce = "GTD"
	TIF_FOK TimeInForce = "FOK"
	TIF_IOC TimeInForce = "IOC"
)

type OrderType string

const (
	ORDER_LIMIT  OrderType = "LIMIT"
	ORDER_MARKET OrderType = "MARKET"
)

type EngineType int64

const (
	PERPETUAL EngineType = iota
	SPOT
)

type OrderSide int64

const (
	BUY OrderSide = iota
	SELL
)

type Order struct {
	Subaccount           string      `json:"subaccount"`
	Sender               string      `json:"sender"`
	Nonce                string      `json:"nonce"` // string of nanoseconds
	Type                 OrderType   `json:"type"`  // LIMIT or MARKET
	Quantity             string      `json:"quantity"`
	Side                 OrderSide   `json:"side"` // 0 BUY, 1 SELL TODO: enum
	OnchainID            int64       `json:"onchainId"`
	EngineType           EngineType  `json:"engineType"` // TODO: enum
	ClientOrderID        string      `json:"clientOrderId,omitempty"`
	ReduceOnly           bool        `json:"reduceOnly"`
	Close                bool        `json:"close,omitempty"`
	StopPrice            int64       `json:"stopPrice,omitempty"`
	StopType             int64       `json:"stopType,omitempty"`
	SignedAt             int64       `json:"signedAt"` // seconds since epoch
	ExpiresAt            int64       `json:"expiresAt,omitempty"`
	GroupID              string      `json:"groupId,omitempty"` // UUID
	GroupContingencyType int         `json:"groupContingencyType,omitempty"`
	Price                string      `json:"price"`
	TimeInForce          TimeInForce `json:"timeInForce"`
	PostOnly             bool        `json:"postOnly"`
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
