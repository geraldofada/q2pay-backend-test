package rest

import (
	"github.com/geraldofada/ledger-backend/src/port"
	"github.com/gofiber/fiber/v2"
)

func injectCreatePayee(app port.AppAccountPort) fiberReturnCtx {
	return func(c *fiber.Ctx) error {
		return createPayee(c, app)
	}
}

func injectLoginPayee(app port.AppAccountPort) fiberReturnCtx {
	return func(c *fiber.Ctx) error {
		return loginPayee(c, app)
	}
}

func (r *Rest) setupPayeeRoutes(app port.AppAccountPort) {
	payee := r.fiber.Group("/payee")

	payee.Post("/", injectCreatePayee(app))
	payee.Post("/login", injectLoginPayee(app))
}
