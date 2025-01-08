package models

import "time"

type CurrencyRate struct {
	RateID   int64     `db:"rate_id"`
	Currency string    `db:"currency"`
	Date     time.Time `db:"date"`
	Rate     float64   `db:"rate"`
}

type CurrencyStat struct {
	Currency string  `db:"currency"`
	MinRate  float64 `db:"min_rate"`
	MaxRate  float64 `db:"max_rate"`
	AvgRate  float64 `db:"avg_rate"`
}

type CurrencyRateResponse struct {
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
}

type CurrencyStatsResponse struct {
	Base         string          `json:"base"`
	RatesAnalyze map[string]Stat `json:"rates_analyze"`
}

type Stat struct {
	MinRate float64 `json:"min"`
	MaxRate float64 `json:"max"`
	AvgRate float64 `json:"avg"`
}
