package wss_client

import "encoding/json"

type RequestType string

const WssSubscribeEvent RequestType = "subscribe"
const WssUnsubscribeEvent RequestType = "unsubscribe"

type WssMessage[T eventData] struct {
	Event RequestType `json:"event"` // subscribe
	Data  T
}

type EventType string

const BookEventType EventType = "L2Book"
const MarketPriceEventType EventType = "MarketPrice"

type eventData interface{}

type WssSymbolData struct {
	Type   EventType `json:"type"`
	Symbol string    `json:"symbol"`
	eventData
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
	MarketPrice EventType `json:"e"`
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
