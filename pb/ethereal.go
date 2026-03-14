package pb

import (
	"encoding/json"
	"errors"
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var _NO_INTENT_ERROR = errors.New("No marshal intent")

type eventData interface{}

type SymbolEvent struct {
	eventData
	T string `json:"type"`
	S string `json:"symbol"`
}

type SubaccountEvent struct {
	eventData
	T string `json:"type"`
	S string `json:"subaccountId"`
}

type Intent string

const (
	Sub   Intent = "subscribe"
	Unsub Intent = "unsubscribe"
)

type SubscriptionIntent[T eventData] struct {
	I Intent    `json:"event"`
	D eventData `json:"data"`
}

// intent abstraction   |   TODO: Subscription Intent protos
func sub(e eventData) ([]byte, error) {
	return json.Marshal(&SubscriptionIntent[eventData]{I: Sub, D: e})
}

func unsub(e eventData) ([]byte, error) {
	return json.Marshal(&SubscriptionIntent[eventData]{I: Unsub, D: e})
}

type Event[T proto.Message] interface {
	EventName() string
	EventStruct() (Event[T], error)
	MarshalIntent(to string, i Intent) ([]byte, error)
	UnmarshalToCallback(b json.RawMessage, cb func(T)) error
}

/////////////////
// ENUM EVENTS //
/////////////////

// server -> client lookup

func EventEnum(e string) EventType {
	switch e {
	case "L2Book":
		return EventType_EVENT_TYPE_L2_BOOK
	case "MarketPrice":
		return EventType_EVENT_TYPE_MARKET_PRICE
	case "SubaccountLiquidation":
		return EventType_EVENT_TYPE_SUBACCOUNT_LIQUIDATION
	case "OrderFill":
		return EventType_EVENT_TYPE_ORDER_FILL
	case "OrderUpdate":
		return EventType_EVENT_TYPE_ORDER_UPDATE
	case "TradeFill":
		return EventType_EVENT_TYPE_TRADE_FILL
	case "TokenTransfer":
		return EventType_EVENT_TYPE_TRANSFER
	default:
		return EventType_EVENT_TYPE_UNSPECIFIED
	}
}

// client -> server lookup

func (e EventType) EventName() string {
	switch e {
	case EventType_EVENT_TYPE_L2_BOOK:
		return "L2Book"
	case EventType_EVENT_TYPE_MARKET_PRICE:
		return "MarketPrice"
	case EventType_EVENT_TYPE_SUBACCOUNT_LIQUIDATION:
		return "SubaccountLiquidation"
	case EventType_EVENT_TYPE_ORDER_FILL:
		return "OrderFill"
	case EventType_EVENT_TYPE_ORDER_UPDATE:
		return "OrderUpdate"
	case EventType_EVENT_TYPE_TRADE_FILL:
		return "TradeFill"
	case EventType_EVENT_TYPE_TRANSFER:
		return "TokenTransfer"
	default:
		return e.String()
	}
}

func (e EventType) EventStruct() (Event[proto.Message], error) {
	switch e {
	case EventType_EVENT_TYPE_L2_BOOK:
		return new(L2Book), nil
	case EventType_EVENT_TYPE_MARKET_PRICE:
		return new(MarketPrice), nil
	case EventType_EVENT_TYPE_SUBACCOUNT_LIQUIDATION:
		return new(SubaccountLiquidation), nil
	case EventType_EVENT_TYPE_ORDER_FILL:
		return new(OrderFill), nil
	case EventType_EVENT_TYPE_ORDER_UPDATE:
		return new(OrderUpdate), nil
	case EventType_EVENT_TYPE_TRADE_FILL:
		return new(TradeFill), nil
	case EventType_EVENT_TYPE_TRANSFER:
		return new(Transfer), nil
	default:
		return nil, _NO_INTENT_ERROR
	}
}

func (e EventType) MarshalIntent(to string, i Intent) ([]byte, error) {
	fmt.Println(i, to)
	if s, err := e.EventStruct(); err == nil {
		return s.MarshalIntent(to, i)
	} else {
		fmt.Println(err)
	}

	return nil, _NO_INTENT_ERROR
}

func (e EventType) UnmarshalToCallback(b json.RawMessage, cb func(proto.Message)) error {
	if s, err := e.EventStruct(); err == nil {
		return s.UnmarshalToCallback(b, cb)
	}
	return _NO_INTENT_ERROR
}

/////////////////////
// PROTOBUF EVENTS //
/////////////////////

func (*L2Book) MarshalIntent(to string, i Intent) ([]byte, error) {
	e := &SymbolEvent{
		T: "L2Book",
		S: to,
	}
	switch i {
	case Sub:
		return sub(e)
	case Unsub:
		return unsub(e)
	}
	return nil, _NO_INTENT_ERROR
}

func (*MarketPrice) MarshalIntent(to string, i Intent) ([]byte, error) {
	e := &SymbolEvent{
		T: "MarketPrice",
		S: to,
	}
	switch i {
	case Sub:
		return sub(e)
	case Unsub:
		return unsub(e)
	}
	return nil, _NO_INTENT_ERROR
}

func (*SubaccountLiquidation) MarshalIntent(to string, i Intent) ([]byte, error) {
	e := &SubaccountEvent{
		T: "SubaccountLiquidation",
		S: to,
	}
	switch i {
	case Sub:
		return sub(e)
	case Unsub:
		return unsub(e)
	}
	return nil, _NO_INTENT_ERROR
}

func (*OrderFill) MarshalIntent(to string, i Intent) ([]byte, error) {
	e := &SubaccountEvent{
		T: "OrderFill",
		S: to,
	}
	switch i {
	case Sub:
		return sub(e)
	case Unsub:
		return unsub(e)
	}
	return nil, _NO_INTENT_ERROR
}

func (*OrderUpdate) MarshalIntent(to string, i Intent) ([]byte, error) {
	e := &SubaccountEvent{
		T: "OrderUpdate",
		S: to,
	}
	switch i {
	case Sub:
		return sub(e)
	case Unsub:
		return unsub(e)
	}
	return nil, _NO_INTENT_ERROR
}

func (*TradeFill) MarshalIntent(to string, i Intent) ([]byte, error) {
	e := &SymbolEvent{
		T: "TradeFill",
		S: to,
	}
	switch i {
	case Sub:
		return sub(e)
	case Unsub:
		return unsub(e)
	}
	return nil, _NO_INTENT_ERROR
}

func (*Transfer) MarshalIntent(to string, i Intent) ([]byte, error) {
	e := &SymbolEvent{
		T: "TokenTransfer",
		S: to,
	}
	switch i {
	case Sub:
		return sub(e)
	case Unsub:
		return unsub(e)
	}
	return nil, _NO_INTENT_ERROR
}

func (l *L2Book) UnmarshalToCallback(b json.RawMessage, cb func(proto.Message)) (err error) {
	if err = protojson.Unmarshal(b, l); err == nil {
		cb(l)
	}
	return
}

func (m *MarketPrice) UnmarshalToCallback(b json.RawMessage, cb func(proto.Message)) (err error) {
	if err = protojson.Unmarshal(b, m); err == nil {
		cb(m)
	}
	return
}

func (s *SubaccountLiquidation) UnmarshalToCallback(b json.RawMessage, cb func(proto.Message)) (err error) {
	if err = protojson.Unmarshal(b, s); err == nil {
		cb(s)
	}
	return
}

func (of *OrderFill) UnmarshalToCallback(b json.RawMessage, cb func(proto.Message)) (err error) {
	if err = protojson.Unmarshal(b, of); err == nil {
		cb(of)
	}
	return
}

func (ou *OrderUpdate) UnmarshalToCallback(b json.RawMessage, cb func(proto.Message)) (err error) {
	if err = protojson.Unmarshal(b, ou); err == nil {
		cb(ou)
	}
	return
}

func (tf *TradeFill) UnmarshalToCallback(b json.RawMessage, cb func(proto.Message)) (err error) {
	if err = protojson.Unmarshal(b, tf); err == nil {
		cb(tf)
	}
	return
}

func (t *Transfer) UnmarshalToCallback(b json.RawMessage, cb func(proto.Message)) (err error) {
	if err = protojson.Unmarshal(b, t); err == nil {
		cb(t)
	}
	return
}

func (*L2Book) EventName() string {
	return "L2Book"
}
func (*MarketPrice) EventName() string {
	return "MarketPrice"
}
func (*SubaccountLiquidation) EventName() string {
	return "SubaccountLiquidation"
}
func (*OrderFill) EventName() string {
	return "OrderFill"
}
func (*OrderUpdate) EventName() string {
	return "OrderUpdate"
}
func (*TradeFill) EventName() string {
	return "TradeFill"
}
func (*Transfer) EventName() string {
	return "Transfer"
}

func (p *L2Book) EventStruct() (Event[proto.Message], error) {
	return p, nil
}
func (p *MarketPrice) EventStruct() (Event[proto.Message], error) {
	return p, nil
}
func (p *SubaccountLiquidation) EventStruct() (Event[proto.Message], error) {
	return p, nil
}
func (p *OrderFill) EventStruct() (Event[proto.Message], error) {
	return p, nil
}
func (p *OrderUpdate) EventStruct() (Event[proto.Message], error) {
	return p, nil
}
func (p *TradeFill) EventStruct() (Event[proto.Message], error) {
	return p, nil
}
func (p *Transfer) EventStruct() (Event[proto.Message], error) {
	return p, nil
}
