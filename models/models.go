package models

type CurrencyUpdate struct {
	Currencies []CurrencyValue `json:"currencies"`
}

type CurrencyValue struct {
	Name   string  `json:"name"`
	Value  float64 `json:"value"`
	Symbol string  `json:"symbol"`
}
