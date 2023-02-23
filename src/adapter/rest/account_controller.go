package rest

import (
	"errors"

	"github.com/geraldofada/q2pay-backend-test/src/core"
	"github.com/geraldofada/q2pay-backend-test/src/port"
	"github.com/gofiber/fiber/v2"
)

type signupAccountInput struct {
	Name     string `validate:"required,min=3,max=32"`
	Email    string `validate:"required,email,min=6,max=32"`
	Doc      string `validate:"required,min=6,max=32"`
	Password string `validate:"required,min=3,max=64"`
}

type loginAccountInput struct {
	Email    string `validate:"required,email,min=6,max=32"`
	Password string `validate:"required,min=3,max=64"`
}

func signupAccount(c *fiber.Ctx, app port.AccountUseCase) error {
	var input signupAccountInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errValidator := validateInput(input)
	if errValidator != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidator)
	}

	account, err := app.AccountSignup(input.Name, input.Email, input.Doc, input.Password)
	if err != nil {
		if errors.Is(err, core.AccountDuplicateError{}) {
			return c.Status(fiber.StatusPreconditionFailed).JSON(fiber.Map{
				"message": "Email or document already exists",
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(account)
}

func loginAccount(c *fiber.Ctx, app port.AccountUseCase) error {
	var input loginAccountInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errValidator := validateInput(input)
	if errValidator != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidator)
	}

	account, token, err := app.AccountLogin(input.Email, input.Password)
	if err != nil {
		if errors.Is(err, core.AccountInvalidPasswordError{}) {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"message": "Wrong e-mail or password",
			})
		}

		if errors.Is(err, core.AccountNotFoundError{}) {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"message": "Wrong e-mail or password",
			})
		}
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"account": account,
		"token": token,
	})
}
