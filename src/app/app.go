package app

import (
	"errors"

	"github.com/geraldofada/q2pay-backend-test/src/core"
	"github.com/geraldofada/q2pay-backend-test/src/port"
)

type App struct {
	accountRepo port.AccountRepository
}

func New(accountRepo port.AccountRepository) *App {
	return &App{accountRepo: accountRepo}
}

func (app *App) AccountLogin(email string, password string) (core.Account, core.Token, error) {
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

func (app *App) AccountSignup(name string, email string, password string, doc string) (core.Account, error) {
	newAccount, err := core.NewAccount(name, email, password, doc)
	if err != nil {
		// app.log.Fatal("Signup account creation", "error", err)
		panic(err)
	}

	err = app.accountRepo.CreateAccount(newAccount)
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

func (app *App) Authorize(token core.Token) (bool, error) {
	authorized, err := token.Authorize()
	if err != nil {
		if errors.Is(err, core.TokenMissingError{}) {
			// app.log.Info("Authorize failed, missing token", "auth")
			return false, err
		}
		if errors.Is(err, core.TokenInvalidError{}) {
			// app.log.Info("Authorize failed, invalid token", "auth")
			return false, err
		}
		// app.log.Fatal("Authorize failed", "error", err)
		panic(err)
	}

	// if authorized {
	// 	app.log.Info("An user was authorized", "auth")
	// } else {
	// 	app.log.Info("An user tried to get authorization", "auth")
	// }
	return authorized, nil
}
