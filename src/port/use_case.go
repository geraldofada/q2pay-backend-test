package port

import "github.com/geraldofada/q2pay-backend-test/src/core"

type AccountUseCase interface {
	AccountLogin(email string, password string) (core.Account, core.Token, error)
	AccountSignup(name string, email string, doc string, password string, accType core.AccountType) (core.Account, error)
	AccountTransferMoney(amount string, srcEmailOrDoc string, targetEmailOrDoc string) (bool, error)
}

type AuthUseCase interface {
	Authorize(token core.Token) (bool, error)
}
