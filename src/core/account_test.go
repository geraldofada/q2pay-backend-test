package core

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../.env"); err != nil {
		panic("Error loading .env")
	}

	os.Exit(m.Run())
}

func TestNewAccount(t *testing.T) {
	account1, err := NewAccount("T", "t@teste.com", "123", "11122233344", "COMMON")
	account2, _ := NewAccount("T", "t@teste.com", "123", "11122233344", "SELLER")

	if err != nil {
		t.Error("Salt failed to generate")
	}

	t.Log("It should return AccountInvalidType with an invalid account type")
	_, err = NewAccount("T", "t@teste.com", "123", "11122233344", "invalid")
	if err == nil || !errors.Is(err, AccountInvalidType{}) {
		t.Error("Expected AccountInvalidType error to return")
	}

	t.Log("It should always have a new hash for the Password")
	if account1.Password == account2.Password {
		t.Error("Equal password hash found")
	}

	t.Log("It should always have a new salt")
	if account1.Salt == account2.Salt {
		t.Error("Equal salt found")
	}

	t.Log("It should be created with balance zeroed")
	if account1.Balance.Amount != 0 {
		t.Error("Balance not zeroed found")
	}
}

func TestAccount_Login(t *testing.T) {
	account1, err := NewAccount("T", "t@teste.com", "123", "11122233344", "COMMON")

	if err != nil {
		t.Error("Salt failed to generate")
	}

	token, err := account1.Login("1234")

	t.Log("It should return AccountInvalidPasswordError{} in case of wrong password")
	if !errors.Is(err, AccountInvalidPasswordError{}) {
		t.Error("Found wrong error type")
	}

	t.Log("It should return an empty token on AccountInvalidPasswordError{}")
	if token != "" {
		t.Error("Found filled token while erroing")
	}

	token, err = account1.Login("123")
	t.Log("It should an err nil if login was successful")
	if err != nil {
		t.Error("An error returned with a valid login")
	}

	t.Log("It should return a valid jwt when err is nil")
	if token == "" {
		t.Error("Found empty token on a valid login")
	} else {
		_, err := jwt.ParseWithClaims(
			string(token),
			&jwt.StandardClaims{},
			func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("invalid signing method")
				}
				return []byte(os.Getenv("JWT_SECRET")), nil
			},
		)

		if err != nil {
			t.Error("An invalid token returned on a valid login")
		}
	}

	// Adding sleep here because the tokens were being generated too fast
	time.Sleep(1 * time.Second)
	token2, _ := account1.Login("123")
	t.Log("It should always return an unique token on multiples logins")
	if token == token2 {
		t.Error("Found duplicate tokens")
	}
}
