package port

import "github.com/geraldofada/q2pay-backend-test/src/core"

type AccountRepository interface {
	CreateAccount(account *core.Account) error
	GetAccountByEmail(email string) (core.Account, error)
	GetAccountByDoc(doc string) (core.Account, error)
	SaveMoneyTransferBetweenAccounts(src *core.Account, target *core.Account) (bool, error)
}
