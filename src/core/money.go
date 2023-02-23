package core

import (
	"fmt"
)

const (
	BRL Currency = "BRL"
	USD Currency = "USD"
)

type Currency string

type Money struct {
	Amount   int64
	Currency Currency
}

func (m Money) Format() string {
	reals := m.Amount / 100
	cents := m.Amount % 100

	if m.Amount < 0 {
		return fmt.Sprintf("%s 0,00", m.Currency)
	}

	if cents < 10 {
		return fmt.Sprintf("%s %d,0%d", m.Currency, reals, cents)
	} else {
		return fmt.Sprintf("%s %d,%d", m.Currency, reals, cents)
	}
}

// func ParseStringToMoney(money string) Money {
// }
