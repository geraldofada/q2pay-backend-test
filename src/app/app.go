package app

import (
	"errors"

	"github.com/geraldofada/q2pay-backend-test/src/core"
	"github.com/geraldofada/q2pay-backend-test/src/port"
)

type App struct {
	payeeRepo port.PayeeRepository
}

func New(payeeRepo port.PayeeRepository) *App {
	return &App{payeeRepo: payeeRepo}
}

func (app *App) PayeeLogin(email string, password string) (core.Payee, core.Token, error) {
	payee, err := app.payeeRepo.GetPayeeByEmail(email)
	if err != nil {
		if errors.Is(err, core.PayeeNotFoundError{}) {
			// app.log.Info("Login failed, not found", "payee", payee)
			return payee, "", err
		}

		// app.log.Fatal("Login account", "error", err)
		panic(err)
	}

	token, err := payee.Login(password)
	if err != nil {
		if errors.Is(err, core.PayeeInvalidPasswordError{}) {
			// app.log.Info("Login failed, invalid password", "payee", payee)
			return payee, token, err
		}

		// app.log.Fatal("Login account", "error", err)
		panic(err)
	}

	// app.log.Info("Login ocurred", "payee", payee)
	return payee, token, nil
}

func (app *App) PayeeSignup(name string, email string, password string, doc string) (core.Payee, error) {
	newPayee, err := core.NewPayee(name, email, password, doc)
	if err != nil {
		// app.log.Fatal("Signup payee creation", "error", err)
		panic(err)
	}

	err = app.payeeRepo.CreatePayee(newPayee)
	if err != nil {
		if errors.Is(err, core.PayeeDuplicateError{}) {
			// app.log.Info("Signup payee failed, email or document already exists", "payee", newPayee)
			return newPayee, err
		}
		// app.log.Fatal("Signup account creation", "error", err)
		panic(err)
	}

	// app.log.Info("Signup account created", "account", newAcc)
	return newPayee, nil
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
