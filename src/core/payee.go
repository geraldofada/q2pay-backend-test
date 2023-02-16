package core

import (
	"encoding/base64"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type Payee struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Doc      string `json:"doc"`
	Salt     string `json:"-"`
	Password string `json:"-"`
}

type PayeeDuplicateError struct{}
type PayeeNotFoundError struct{}
type PayeeInvalidPasswordError struct{}
type PayeeMissingToken struct{}
type PayeeInvalidToken struct{}

func (e PayeeDuplicateError) Error() string {
	return "payee duplicate"
}
func (e PayeeNotFoundError) Error() string {
	return "payee not found"
}
func (e PayeeInvalidPasswordError) Error() string {
	return "payee invalid password"
}
func (e PayeeMissingToken) Error() string {
	return "payee missing token"
}
func (e PayeeInvalidToken) Error() string {
	return "account invalid token"
}

func NewPayee(name string, email string, password string, doc string) (Payee, error) {
	pepper := os.Getenv("PASS_SECRET")
	salt, err := generateSalt(32)
	if err != nil {
		return Payee{}, err
	}
	salt64 := base64.StdEncoding.EncodeToString(salt)

	secret := pepper + ":" + salt64

	hashedPassword := base64.StdEncoding.EncodeToString(hashPass(password, []byte(secret), 32))

	return Payee{
		Name:     name,
		Email:    email,
		Doc:      doc,
		Salt:     salt64,
		Password: hashedPassword,
	}, nil
}

func (a *Payee) Login(password string) (Token, error) {
	pepper := os.Getenv("PASS_SECRET")
	secret := pepper + ":" + a.Salt

	hashedPassword := base64.StdEncoding.EncodeToString(hashPass(password, []byte(secret), 32))

	if a.Password != hashedPassword {
		return "", PayeeInvalidPasswordError{}
	}

	token, err := generateJwt()
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *Payee) Authorize(token Token) (bool, error) {
	if token == "" {
		return false, PayeeInvalidToken{}
	}

	validatedToken, err := validateJwt(token)
	if err != nil {
		return false, err
	}

	if !validatedToken.Valid {
		return false, PayeeInvalidToken{}
	}

	claims, ok := validatedToken.Claims.(*jwt.StandardClaims)
	if !ok {
		return false, AccountInvalidToken{}
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return false, AccountInvalidToken{}
	}

	return true, nil
}
