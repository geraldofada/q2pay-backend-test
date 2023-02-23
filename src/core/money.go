package core

import (
	"fmt"
	"strconv"
	"strings"
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

type MoneyParseError struct{}
type MoneyNotEnoughToWithdrawError struct{}
type MoneyMismatchCurrencyError struct{}

func (e MoneyParseError) Error() string {
	return "invalid string to parse as Money"
}

func (e MoneyNotEnoughToWithdrawError) Error() string {
	return "not enough money to withdraw"
}

func (e MoneyMismatchCurrencyError) Error() string {
	return "mismatch currency"
}

// IMPORTANT: remember to add new currencies in here
func validateCurrencyFromString(toValidate string) bool {
	switch toValidate {
	case string(BRL), string(USD):
		return true
	}

	return false
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

func ParseStringToMoney(money string) (Money, error) {
	currency := strings.Split(money, " ")

	if len(currency) == 1 {
		return Money{}, MoneyParseError{}
	}

	if !validateCurrencyFromString(currency[0]) {
		return Money{}, MoneyParseError{}
	}

	amount := strings.Split(currency[1], ",")
	if len(amount) != 2 {
		return Money{}, MoneyParseError{}
	}

	reals, err := strconv.Atoi(amount[0])
	if err != nil {
		return Money{}, MoneyParseError{}
	}

	cents, err := strconv.Atoi(amount[1])
	if err != nil {
		return Money{}, MoneyParseError{}
	}

	amountInt := (reals * 100) + cents

	if amountInt < 0 {
		return Money{}, MoneyParseError{}
	}

	return Money{
		Currency: Currency(currency[0]),
		Amount:   int64(amountInt),
	}, nil
}

func (m *Money) withdraw(toWithdraw Money) (bool, error) {
	if m.Amount < toWithdraw.Amount || toWithdraw.Amount < 0 {
		return false, MoneyNotEnoughToWithdrawError{}
	}

	if m.Currency != toWithdraw.Currency {
		return false, MoneyMismatchCurrencyError{}
	}

	m.Amount -= toWithdraw.Amount

	return true, nil
}

func (m *Money) deposit(toDeposit Money) (bool, error) {
	if toDeposit.Amount < 0 {
		return false, MoneyNotEnoughToWithdrawError{}
	}

	if m.Currency != toDeposit.Currency {
		return false, MoneyMismatchCurrencyError{}
	}

	m.Amount += toDeposit.Amount

	return true, nil
}
