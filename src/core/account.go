package core

import (
	"encoding/base64"
	"os"

	"gorm.io/gorm"
)

type AccountType string

const (
	COMMON AccountType = "COMMON"
	SELLER AccountType = "SELLER"
)

type Account struct {
	gorm.Model
	Name     string      `json:"name"`
	Email    string      `json:"email"`
	Doc      string      `json:"doc"`
	Balance  Money       `json:"balance" gorm:"type:string"`
	Type     AccountType `json:"type" gorm:"type:string"`
	Salt     string      `json:"-"`
	Password string      `json:"-"`
}

type AccountDuplicateError struct{}
type AccountNotFoundError struct{}
type AccountInvalidPasswordError struct{}
type AccountInvalidTypeError struct{}

func (e AccountDuplicateError) Error() string {
	return "account duplicate"
}
func (e AccountNotFoundError) Error() string {
	return "account not found"
}
func (e AccountInvalidPasswordError) Error() string {
	return "account invalid password"
}

func (e AccountInvalidTypeError) Error() string {
	return "account invalid type"
}

// IMPORTANT: remember to add new types in here
func validateAccountType(toValidate string) bool {
	switch toValidate {
	case string(COMMON), string(SELLER):
		return true
	}

	return false
}

func NewAccount(name string, email string, password string, doc string, accType string) (Account, error) {
	if !validateAccountType(accType) {
		return Account{}, AccountInvalidTypeError{}
	}

	pepper := os.Getenv("PASS_SECRET")
	salt, err := generateSalt(32)
	if err != nil {
		return Account{}, err
	}
	salt64 := base64.StdEncoding.EncodeToString(salt)

	secret := pepper + ":" + salt64

	hashedPassword := base64.StdEncoding.EncodeToString(hashPass(password, []byte(secret), 32))

	return Account{
		Name:     name,
		Email:    email,
		Doc:      doc,
		Balance:  Money{Amount: 0, Currency: BRL},
		Salt:     salt64,
		Password: hashedPassword,
	}, nil
}

func (a *Account) Login(password string) (Token, error) {
	pepper := os.Getenv("PASS_SECRET")
	secret := pepper + ":" + a.Salt

	hashedPassword := base64.StdEncoding.EncodeToString(hashPass(password, []byte(secret), 32))

	if a.Password != hashedPassword {
		return "", AccountInvalidPasswordError{}
	}

	token, err := generateJwt()
	if err != nil {
		return "", err
	}

	return token, nil
}
