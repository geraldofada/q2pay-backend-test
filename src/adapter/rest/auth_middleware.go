package rest

import (
	"errors"
	"strings"

	"github.com/geraldofada/q2pay-backend-test/src/core"
	"github.com/geraldofada/q2pay-backend-test/src/port"
	"github.com/gofiber/fiber/v2"
)

func authorize(auth port.AuthUseCase) fiberReturnCtx {
	return func(c *fiber.Ctx) error {
		bearer := c.Get("Authorization", "")

		if bearer == "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Missing token",
			})
		}

		_, token, found := strings.Cut(bearer, "Bearer ")
		if !found {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Invalid token",
			})
		}

		authorized, accId, err := auth.Authorize(core.Token(token))
		if err != nil {
			if errors.Is(err, core.TokenMissingError{}) {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"message": "Invalid token",
				})
			}
			if errors.Is(err, core.TokenInvalidError{}) {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"message": "Invalid token",
				})
			}
		}

		if !authorized {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Not enough privileges",
			})
		}

		c.Locals("loggedId", accId)
		return c.Next()
	}
}
