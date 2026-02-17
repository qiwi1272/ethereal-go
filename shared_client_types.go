package ethereal

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

type Subaccount struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Account string `json:"account"`
}
