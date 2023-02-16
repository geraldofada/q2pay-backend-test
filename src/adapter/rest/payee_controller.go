package rest

import (
	"errors"

	account "github.com/geraldofada/ledger-backend/src/core"
	"github.com/geraldofada/ledger-backend/src/port"
	"github.com/gofiber/fiber/v2"
)

type createPayeeInput struct {
	Name     string `validate:"required,min=3,max=32"`
	Email    string `validate:"required,email,min=6,max=32"`
	Password string `validate:"required,min=3,max=64"`
}

type loginPayeeInput struct {
	Email    string `validate:"required,email,min=6,max=32"`
	Password string `validate:"required,min=3,max=64"`
}

func createPayee(c *fiber.Ctx, app port.AppAccountPort) error {
	var input createPayeeInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errValidator := validateInput(input)
	if errValidator != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidator)
	}

	acc, err := app.Signup(input.Name, input.Email, input.Password)
	if err != nil {
		if errors.Is(err, account.AccountDuplicateError{}) {
			return c.Status(fiber.StatusPreconditionFailed).JSON(fiber.Map{
				"message": "Email already exists",
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(acc)
}

func loginPayee(c *fiber.Ctx, app port.AppAccountPort) error {
	var input loginPayeeInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errValidator := validateInput(input)
	if errValidator != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidator)
	}

	acc, err := app.Login(input.Email, input.Password)
	if err != nil {
		if errors.Is(err, account.AccountInvalidPasswordError{}) {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"message": "Wrong e-mail or password",
			})
		}

		if errors.Is(err, account.AccountNotFoundError{}) {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"message": "Wrong e-mail or password",
			})
		}
	}

	return c.Status(fiber.StatusAccepted).JSON(acc)
}
