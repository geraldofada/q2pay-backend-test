package port

import "github.com/geraldofada/q2pay-backend-test/src/core"

type PayeeRepository interface {
	CreatePayee(payee core.Payee) error
	GetPayeeByEmail(email string) (core.Payee, error)
	GetPayeeByDoc(doc string) (core.Payee, error)
}
