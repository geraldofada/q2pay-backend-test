package port

import "github.com/geraldofada/q2pay-backend-test/src/core"

type PayeeUseCase interface {
	PayeeLogin(email string, password string) (core.Payee, core.Token, error)
	PayeeSignup(name string, email string, doc string, password string) (core.Payee, error)
}

type AuthUseCase interface {
	Authorize(token core.Token) (bool, error)
}
