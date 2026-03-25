package model

// USDJPYRate represents one daily USDJPY snapshot loaded from JSON files.
type USDJPYRate struct {
	Date  string  `json:"date"`
	Label string  `json:"label,omitempty"`
	Pair  string  `json:"pair"`
	Bid   float64 `json:"bid"`
	Ask   float64 `json:"ask"`
	Open  float64 `json:"open"`
	High  float64 `json:"high"`
	Low   float64 `json:"low"`
	Diff  float64 `json:"diff"`
	Close float64 `json:"close"`
}

// USDJPYRatesResponse is the API payload for the USDJPY chart.
type USDJPYRatesResponse struct {
	Pair      string       `json:"pair"`
	Timeframe string       `json:"timeframe"`
	Rates     []USDJPYRate `json:"rates"`
}
