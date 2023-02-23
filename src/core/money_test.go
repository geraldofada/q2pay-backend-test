package core

import (
	"errors"
	"testing"
)


func TestMoney_Format(t *testing.T) {
	money1 := Money{ Amount: 0, Currency: BRL }
	money2 := Money{ Amount: -1, Currency: BRL }
	money3 := Money{ Amount: 9, Currency: BRL }
	money4 := Money{ Amount: 100, Currency: BRL }
	money5 := Money{ Amount: 99999, Currency: BRL }

	t.Log("It should return BRL 0,00")
	if money1.Format() != "BRL 0,00" {
		t.Error("Wrong format for Amount 0")
	}

	t.Log("It should return BRL 0,00")
	if money2.Format() != "BRL 0,00" {
		t.Error("Wrong format for negative Amount")
	}

	t.Log("It should return BRL 0,09")
	if money3.Format() != "BRL 0,09" {
		t.Error("Wrong format for Amount lower than 10")
	}

	t.Log("It should return BRL 1,00")
	if money4.Format() != "BRL 1,00" {
		t.Error("Wrong format for Amount higher than 10")
	}

	t.Log("It should return BRL 999,99")
	if money5.Format() != "BRL 999,99" {
		t.Error("Wrong format for Amount higher than 10")
	}
}

func TestParseStringToMoney(t *testing.T) {
	money1 := Money{ Amount: 0, Currency: BRL }
	money3 := Money{ Amount: 9, Currency: BRL }
	money4 := Money{ Amount: 100, Currency: BRL }
	money5 := Money{ Amount: 99999, Currency: BRL }

	t.Log("It should return Amount 0, Currency BRL")
	m, err := ParseStringToMoney("BRL 0,00")
	if err != nil {
		t.Error("Unexpected error while parsing string to money")
  }
	if money1 != m {
		t.Error("Parse string to money got wrong struct")
	}

	t.Log("It should return Amount 9, Currency BRL")
	m, err = ParseStringToMoney("BRL 0,09")
	if err != nil {
		t.Error("Unexpected error while parsing string to money")
  }
	if money3 != m {
		t.Error("Parse string to money got wrong struct")
	}

	t.Log("It should return Amount 100, Currency BRL")
	m, err = ParseStringToMoney("BRL 1,00")
	if err != nil {
		t.Error("Unexpected error while parsing string to money")
  }
	if money4 != m {
		t.Error("Parse string to money got wrong struct")
	}

	t.Log("It should return Amount 99999, Currency BRL")
	m, err = ParseStringToMoney("BRL 999,99")
	if err != nil {
		t.Error("Unexpected error while parsing string to money")
  }
	if money5 != m {
		t.Error("Parse string to money got wrong struct")
	}

	t.Log("It should return MoneyParseError with negative value")
	_, err = ParseStringToMoney("BRL -999,99")
	if !errors.Is(err, MoneyParseError{}) {
		t.Error("Unexpected error while parsing string to money")
  }

	t.Log("It should return MoneyParseError with invalid Currency")
	_, err = ParseStringToMoney("AAAA 999,99")
	if !errors.Is(err, MoneyParseError{}) {
		t.Error("Unexpected error while parsing string to money")
  }

	t.Log("It should return MoneyParseError with invalid string format")
	_, err = ParseStringToMoney("BRL999,99")
	if !errors.Is(err, MoneyParseError{}) {
		t.Error("Unexpected error while parsing string to money")
  }

	t.Log("It should return MoneyParseError with invalid string format")
	_, err = ParseStringToMoney("BRL99999")
	if !errors.Is(err, MoneyParseError{}) {
		t.Error("Unexpected error while parsing string to money")
  }

	t.Log("It should return MoneyParseError with invalid string format")
	_, err = ParseStringToMoney("BRL 99999")
	if !errors.Is(err, MoneyParseError{}) {
		t.Error("Unexpected error while parsing string to money")
  }
}

func TestMoney_withdraw(t *testing.T) {
	money1 := Money{ Amount: 0, Currency: BRL }
	money2 := Money{ Amount: -1, Currency: BRL }
	money3 := Money{ Amount: 9, Currency: BRL }
	money4 := Money{ Amount: 100, Currency: BRL }
	money5 := Money{ Amount: 99, Currency: USD }

	t.Log("It shouldreturn MoneyNotEnoughToWithdrawError if does not have enough money")
	_, err := money1.withdraw(money4)
	if err == nil {
		t.Error("Unexpected error while withdrawing money")
	}
	if !errors.Is(err, MoneyNotEnoughToWithdrawError{}) {
		t.Error("Unexpected error while withdrawing money")
  }

	t.Log("It should return MoneyNotEnoughToWithdrawError if withdrawing negative amount")
	_, err = money1.withdraw(money2)
	if err == nil {
		t.Error("Unexpected error while withdrawing money")
	}
	if !errors.Is(err, MoneyNotEnoughToWithdrawError{}) {
		t.Error("Unexpected error while withdrawing money")
  }

	t.Log("It should return MoneyMismatchCurrencyError if withdrawing from different currencies")
	_, err = money4.withdraw(money5)
	if err == nil {
		t.Error("Unexpected error while withdrawing money")
	}
	if !errors.Is(err, MoneyMismatchCurrencyError{}) {
		t.Error("Unexpected error while withdrawing money")
  }

	t.Log("It should modify the original Money on success")
	ok, _ := money4.withdraw(money3)
	if !ok {
		t.Error("Unexpected error while withdrawing money")
	}
	if money4.Amount != (100 - 9) {
		t.Error("Wrong withdraw value, expected 91")
	}
}

func TestMoney_deposit(t *testing.T) {
	money1 := Money{ Amount: 0, Currency: BRL }
	money2 := Money{ Amount: -1, Currency: BRL }
	money3 := Money{ Amount: 9, Currency: BRL }
	money4 := Money{ Amount: 100, Currency: BRL }
	money5 := Money{ Amount: 99, Currency: USD }

	t.Log("It should return MoneyNotEnoughToWithdrawError if depositing negative amount")
	_, err := money1.deposit(money2)
	if err == nil {
		t.Error("Unexpected error while withdrawing money")
	}
	if !errors.Is(err, MoneyNotEnoughToWithdrawError{}) {
		t.Error("Unexpected error while withdrawing money")
  }

	t.Log("It should return MoneyMismatchCurrencyError if withdrawing from different currencies")
	_, err = money4.deposit(money5)
	if err == nil {
		t.Error("Unexpected error while withdrawing money")
	}
	if !errors.Is(err, MoneyMismatchCurrencyError{}) {
		t.Error("Unexpected error while withdrawing money")
  }

	t.Log("It should modify the original Money on success")
	ok, _ := money4.deposit(money3)
	if !ok {
		t.Error("Unexpected error while withdrawing money")
	}
	if money4.Amount != (100 + 9) {
		t.Error("Wrong withdraw value, expected 109")
	}
}
