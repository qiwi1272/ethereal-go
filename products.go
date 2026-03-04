package ethereal

type Product struct {
	ID                     string          `json:"id"`
	Ticker                 string          `json:"ticker"`
	DisplayTicker          string          `json:"displayTicker"`
	EngineType             OrderEngineType `json:"engineType"`
	OnchainID              int64           `json:"onchainId"`
	LotSize                string          `json:"lotSize"`
	TickSize               string          `json:"tickSize"`
	MakerFee               string          `json:"makerFee"`
	TakerFee               string          `json:"takerFee"`
	MaxQuantity            string          `json:"maxQuantity"`
	MinQuantity            string          `json:"minQuantity"`
	Volume24h              string          `json:"volume24h"`
	FundingRate1h          string          `json:"fundingRate1h"`
	MaxOpenInterestUsd     string          `json:"maxOpenInterestUsd"`
	MaxPositionNotionalUsd string          `json:"maxPositionNotionalUsd"`
}
