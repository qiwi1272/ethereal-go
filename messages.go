package ethereal

import (
	"fmt"
	"math/big"
	"strconv"
	"time"

	abi "github.com/ethereum/go-ethereum/signer/core/apitypes"
)

// -------- BEGIN ENUMS -------- //

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

type OrderEngineType int64

const (
	PERPETUAL OrderEngineType = iota
	SPOT
)

type OrderSide int64

const (
	BUY OrderSide = iota
	SELL
)

// -------- BEGIN ENUMS -------- //

// -------- BEGIN ORDER -------- //

type Order struct {
	Subaccount           string          `json:"subaccount"`
	Sender               string          `json:"sender"`
	Nonce                string          `json:"nonce"` // string of nanoseconds
	Type                 OrderType       `json:"type"`  // LIMIT or MARKET
	Quantity             string          `json:"quantity"`
	Side                 OrderSide       `json:"side"` // 0 BUY, 1 SELL
	OnchainID            int64           `json:"onchainId"`
	EngineType           OrderEngineType `json:"engineType"`
	ClientOrderID        string          `json:"clientOrderId,omitempty"`
	ReduceOnly           bool            `json:"reduceOnly"`
	Close                bool            `json:"close,omitempty"`
	StopPrice            int64           `json:"stopPrice,omitempty"`
	StopType             int64           `json:"stopType,omitempty"`
	SignedAt             int64           `json:"signedAt"` // seconds since epoch
	ExpiresAt            int64           `json:"expiresAt,omitempty"`
	GroupID              string          `json:"groupId,omitempty"` // UUID
	GroupContingencyType int             `json:"groupContingencyType,omitempty"`
	Price                string          `json:"price"`
	TimeInForce          TimeInForce     `json:"timeInForce"`
	PostOnly             bool            `json:"postOnly"`
}

// needed for building
type Subaccount struct {
	Id      string `json:"id"`      // bytes32
	Name    string `json:"name"`    // id
	Account string `json:"account"` // EOA
}

type SubaccountHolder interface {
	GetSubaccount() *Subaccount
}

func NewRawOrder(
	orderType OrderType,
	onchainId int64,
	marketType OrderEngineType,
	qty float64,
	px float64,
	reduce bool,
	side OrderSide,
	tif TimeInForce,
) *Order {
	return &Order{
		Type:        orderType,
		Quantity:    fmt.Sprintf("%.9f", qty),
		Side:        side,
		OnchainID:   onchainId,
		EngineType:  marketType,
		ReduceOnly:  reduce,
		Price:       fmt.Sprintf("%.9f", px),
		TimeInForce: tif,
		PostOnly:    false,
	}
}

func (p *Product) NewOrder(
	orderType OrderType,
	qty float64,
	px float64,
	reduce bool,
	side OrderSide,
	tif TimeInForce,
) *Order {
	return &Order{
		Type:        orderType,
		Quantity:    fmt.Sprintf("%.9f", qty),
		Side:        side,
		OnchainID:   p.OnchainID,
		EngineType:  p.EngineType,
		ReduceOnly:  reduce,
		Price:       fmt.Sprintf("%.9f", px),
		TimeInForce: tif,
		PostOnly:    false,
	}
}

func (o *Order) ToMessage() (abi.TypedDataMessage, error) {
	qtyBig, err := Scale1e9(o.Quantity)
	if err != nil {
		return abi.TypedDataMessage{}, err
	}
	priceBig, err := Scale1e9(o.Price)
	if err != nil {
		return abi.TypedDataMessage{}, err
	}

	// even though we expect these values to be uint8 according to their signatures,
	// setting them as native uint8 raises a compiler error. strings or big ints are accepted.
	engine := big.NewInt(int64(o.EngineType))
	side := big.NewInt(int64(o.Side))
	id := big.NewInt(o.OnchainID)
	sigTs := big.NewInt(o.SignedAt)

	return abi.TypedDataMessage{
		"sender":     o.Sender,
		"subaccount": o.Subaccount,
		"quantity":   qtyBig,
		"price":      priceBig,
		"reduceOnly": o.ReduceOnly,
		"side":       side,
		"engineType": engine,
		"productId":  id,
		"nonce":      o.Nonce,
		"signedAt":   sigTs,
	}, nil
}

func (o *Order) Build(cl SubaccountHolder) {
	sub := cl.GetSubaccount()
	o.Sender = sub.Account
	o.Subaccount = sub.Name
	nonce := getOrderNonce()

	o.Nonce = nonce
	o.SignedAt, _ = strconv.ParseInt(nonce[:len(nonce)-9], 10, 64)
}

// -------- END ORDER -------- //

// -------- BEGIN CANCEL -------- //

type OrderCancel struct {
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

func NewCancel(oids ...string) *OrderCancel {
	return &OrderCancel{
		OrderIDs: oids,
	}
}

func NewCancelOrderFromCreated(orders ...OrderCreated) *OrderCancel {
	ids := make([]string, len(orders))
	for i, o := range orders {
		ids[i] = o.Id
	}
	return NewCancel(ids...)
}

func (o *OrderCancel) ToMessage() (abi.TypedDataMessage, error) {
	return abi.TypedDataMessage{
		"sender":     o.Sender,
		"subaccount": o.Subaccount,
		"nonce":      o.Nonce,
		//"orderIds":   co.OrderIDs,
	}, nil
}

func (o *OrderCancel) Build(cl SubaccountHolder) {
	sub := cl.GetSubaccount()
	o.Sender = sub.Account
	o.Subaccount = sub.Name
	o.Nonce = getOrderNonce()
}

// -------- END CANCEL -------- //

// -------- POSITION -------- //

type Position struct {
	Id                    string    `json:"id"`
	Cost                  string    `json:"cost"`
	Size                  string    `json:"size"`
	FundingUsd            string    `json:"fundingUsd"`
	FundingAccruedUsd     string    `json:"fundingAccruedUsd"`
	FeesAccruedUsd        string    `json:"feesAccruedUsd"`
	RealizedPnl           string    `json:"realizedPnl"`
	TotalIncreaseNotional string    `json:"totalIncreaseNotional"`
	TotalIncreaseQuantity string    `json:"totalIncreaseQuantity"`
	TotalDecreaseNotional string    `json:"totalDecreaseNotional"`
	TotalDecreaseQuantity string    `json:"totalDecreaseQuantity"`
	Side                  OrderSide `json:"side"`
	ProductId             string    `json:"productId"`
	UpdatedAt             uint64    `json:"updatedAt"`
	CreatedAt             uint64    `json:"createdAt"`
	IsLiquidated          bool      `json:"isLiquidated"`
	LiquidationPrice      string    `json:"liquidationPrice"`
}

// -------- BALANCE -------- //

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

// -------- HELPERS -------- //

func getOrderNonce() string {
	now := time.Now()
	return strconv.FormatInt(now.UnixNano(), 10)
}
