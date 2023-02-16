package core

const (
	BRL Currency = "BRL"
	USD Currency = "USD"
)

type Currency string

type Money struct {
	Amount int64
	Currency Currency
}
