package rest

import (
	"github.com/geraldofada/q2pay-backend-test/src/port"
	"github.com/gofiber/fiber/v2"
)

func injectSignupPayee(app port.PayeeUseCase) fiberReturnCtx {
	return func(c *fiber.Ctx) error {
		return signupPayee(c, app)
	}
}

func injectLoginPayee(app port.PayeeUseCase) fiberReturnCtx {
	return func(c *fiber.Ctx) error {
		return loginPayee(c, app)
	}
}

func (r Rest) setupPayeeRoutes(app port.PayeeUseCase) {
	payee := r.Fiber.Group("/payee")

	payee.Post("/", injectSignupPayee(app))
	payee.Post("/login", injectLoginPayee(app))
}
