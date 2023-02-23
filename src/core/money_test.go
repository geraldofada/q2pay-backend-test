package core

import (
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

