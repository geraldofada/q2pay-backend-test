package app

import (
	"errors"

	"github.com/geraldofada/q2pay-backend-test/src/core"
)

type AppAuth struct{}

func NewAppAuth() AppAuth {
	return AppAuth{}
}

func (app AppAuth) Authorize(token core.Token) (bool, uint, error) {
	authorized, accId, err := token.Authorize()
	if err != nil {
		if errors.Is(err, core.TokenMissingError{}) {
			// app.log.Info("Authorize failed, missing token", "auth")
			return false, 0, err
		}
		if errors.Is(err, core.TokenInvalidError{}) {
			// app.log.Info("Authorize failed, invalid token", "auth")
			return false, 0, err
		}
		// app.log.Fatal("Authorize failed", "error", err)
		panic(err)
	}

	// if authorized {
	// 	app.log.Info("An user was authorized", "auth")
	// } else {
	// 	app.log.Info("An user tried to get authorization", "auth")
	// }
	return authorized, accId, nil
}
