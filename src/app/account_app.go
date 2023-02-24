package app

import (
	"errors"

	"github.com/geraldofada/q2pay-backend-test/src/core"
	"github.com/geraldofada/q2pay-backend-test/src/port"
)

type AppAccount struct {
	accountRepo port.AccountRepository
}

func NewAppAccount(accountRepo port.AccountRepository) AppAccount {
	return AppAccount{accountRepo: accountRepo}
}

func (app AppAccount) AccountLogin(email string, password string) (core.Account, core.Token, error) {
	account, err := app.accountRepo.GetAccountByEmail(email)
	if err != nil {
		if errors.Is(err, core.AccountNotFoundError{}) {
			// app.log.Info("Login failed, not found", "account", account)
			return account, "", err
		}

		// app.log.Fatal("Login account", "error", err)
		panic(err)
	}

	token, err := account.Login(password)
	if err != nil {
		if errors.Is(err, core.AccountInvalidPasswordError{}) {
			// app.log.Info("Login failed, invalid password", "account", account)
			return account, token, err
		}

		// app.log.Fatal("Login account", "error", err)
		panic(err)
	}

	// app.log.Info("Login ocurred", "account", account)
	return account, token, nil
}

func (app AppAccount) AccountSignup(name string, email string, password string, doc string, accType core.AccountType) (core.Account, error) {
	newAccount, err := core.NewAccount(name, email, password, doc, accType)
	if err != nil {
		// app.log.Fatal("Signup account creation", "error", err)
		panic(err)
	}

	err = app.accountRepo.CreateAccount(&newAccount)
	if err != nil {
		if errors.Is(err, core.AccountDuplicateError{}) {
			// app.log.Info("Signup account failed, email or document already exists", "account", newAccount)
			return newAccount, err
		}
		// app.log.Fatal("Signup account creation", "error", err)
		panic(err)
	}

	// app.log.Info("Signup account created", "account", newAcc)
	return newAccount, nil
}

func (app AppAccount) AccountTransferMoney(amount string, srcEmailOrDoc string, targetEmailOrDoc string) (bool, error) {
	srcAcc, err := app.accountRepo.GetAccountByEmail(srcEmailOrDoc)

	checkForDoc := false

	if err != nil {
		if errors.Is(err, core.AccountNotFoundError{}) {
			checkForDoc = true
		} else {
			// app.log.Fatal("Signup account creation", "error", err)
			panic(err)
		}
	}

	if checkForDoc {
		srcAcc, err = app.accountRepo.GetAccountByDoc(srcEmailOrDoc)
		if err != nil {
			if errors.Is(err, core.AccountNotFoundError{}) {
				return false, err
			}
			// app.log.Fatal("Signup account creation", "error", err)
			panic(err)
		}
	}

	targetAcc, err := app.accountRepo.GetAccountByEmail(targetEmailOrDoc)

	checkForDoc = false

	if err != nil {
		if errors.Is(err, core.AccountNotFoundError{}) {
			checkForDoc = true
		} else {
			// app.log.Fatal("Signup account creation", "error", err)
			panic(err)
		}
	}

	if checkForDoc {
		targetAcc, err = app.accountRepo.GetAccountByDoc(targetEmailOrDoc)
		if err != nil {
			if errors.Is(err, core.AccountNotFoundError{}) {
				return false, err
			}
			// app.log.Fatal("Signup account creation", "error", err)
			panic(err)
		}
	}

	moneyToTransfer, err := core.ParseStringToMoney(amount)
	if err != nil {
		if errors.Is(err, core.MoneyParseError{}) {
			return false, err
		}
		panic(err)
	}

	_, err = srcAcc.TransferMoney(moneyToTransfer, &targetAcc)
	if err != nil {
		if errors.Is(err, core.AccountNotEnoughBalanceError{}) {
			return false, err
		}

		if errors.Is(err, core.AccountSellerCannotTransferError{}) {
			return false, err
		}

		if errors.Is(err, core.MoneyMismatchCurrencyError{}) {
			return false, err
		}

		panic(err)
	}

	_, err = app.accountRepo.SaveMoneyTransferBetweenAccounts(&srcAcc, &targetAcc)
	if err != nil {
		if errors.Is(err, core.AccountNotFoundError{}) {
			return false, err
		}
		// app.log.Fatal("Signup account creation", "error", err)
		panic(err)
	}

	return true, nil
}

func (app AppAccount) AccountGetById(id uint) (core.Account, error) {
	srcAcc, err := app.accountRepo.GetAccountById(id)

	if err != nil {
		if errors.Is(err, core.AccountNotFoundError{}) {
			return core.Account{}, err
		}
		// app.log.Fatal("Signup account creation", "error", err)
		panic(err)
	}

	return srcAcc, nil
}
