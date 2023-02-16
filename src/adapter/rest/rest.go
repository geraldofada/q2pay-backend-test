package rest

import (
	"github.com/geraldofada/q2pay-backend-test/src/port"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Rest struct {
	fiber *fiber.App
}

type fiberReturnCtx func(*fiber.Ctx) error

func New() Rest {
	rest := fiber.New()
	rest.Use(cors.New())

	return Rest{
		fiber: rest,
	}
}

func (r *Rest) SetupPayeeRoutes(payee port.PayeeUseCase) {
	r.setupPayeeRoutes(payee)
}
