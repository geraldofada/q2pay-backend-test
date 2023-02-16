package core

import (
	"crypto/rand"
	"os"
	"strconv"
	"time"
	"errors"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/argon2"
)

type Token string

func generateSalt(size int) ([]byte, error) {
	salt := make([]byte, size)

	_, err := rand.Read(salt[:])

	if err != nil {
		return nil, err
	}

	return salt, nil
}

func hashPass(password string, salt []byte, size uint32) []byte {
	time, _ := strconv.ParseUint(os.Getenv("PASS_TIME"), 10, 32)
	mem, _ := strconv.ParseUint(os.Getenv("PASS_MEM"), 10, 32)
	threads, _ := strconv.ParseUint(os.Getenv("PASS_THREADS"), 10, 8)
	return argon2.IDKey([]byte(password), salt, uint32(time), uint32(mem), uint8(threads), size)
}

func generateJwt() (Token, error) {
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    "q2-backend-test",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return "", err
	}
	return Token(tokenString), nil
}

func validateJwt(token Token) (*jwt.Token, error) {
	tokenParsed, err := jwt.ParseWithClaims(
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
		return nil, err
	}

	return tokenParsed, nil
}
