package rest

import (
	"github.com/geraldofada/q2pay-backend-test/src/port"
	"github.com/gofiber/fiber/v2"
)

func injectSignupAccount(app port.AccountUseCase) fiberReturnCtx {
	return func(c *fiber.Ctx) error {
		return signupAccount(c, app)
	}
}

func injectLoginAccount(app port.AccountUseCase) fiberReturnCtx {
	return func(c *fiber.Ctx) error {
		return loginAccount(c, app)
	}
}

func injectTransferMoneyAccount(app port.AccountUseCase) fiberReturnCtx {
	return func(c *fiber.Ctx) error {
		return transferMoneyAccount(c, app)
	}
}

func (r Rest) setupAccountRoutes(app port.AccountUseCase) {
	account := r.Fiber.Group("/account")

	account.Post("/", injectSignupAccount(app))
	account.Post("/login", injectLoginAccount(app))

	account.Post("/transfer-money", authorize(appAuth), injectTransferMoneyAccount(appAccount))
}
