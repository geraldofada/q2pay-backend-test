package rest

import (
	"errors"

	"github.com/geraldofada/q2pay-backend-test/src/core"
	"github.com/geraldofada/q2pay-backend-test/src/port"
	"github.com/gofiber/fiber/v2"
)

type signupPayeeInput struct {
	Name     string `validate:"required,min=3,max=32"`
	Email    string `validate:"required,email,min=6,max=32"`
	Doc      string `validate:"required,min=6,max=32"`
	Password string `validate:"required,min=3,max=64"`
}

type loginPayeeInput struct {
	Email    string `validate:"required,email,min=6,max=32"`
	Password string `validate:"required,min=3,max=64"`
}

func signupPayee(c *fiber.Ctx, app port.PayeeUseCase) error {
	var input signupPayeeInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errValidator := validateInput(input)
	if errValidator != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidator)
	}

	payee, err := app.PayeeSignup(input.Name, input.Email, input.Doc, input.Password)
	if err != nil {
		if errors.Is(err, core.PayeeDuplicateError{}) {
			return c.Status(fiber.StatusPreconditionFailed).JSON(fiber.Map{
				"message": "Email or document already exists",
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(payee)
}

func loginPayee(c *fiber.Ctx, app port.PayeeUseCase) error {
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

	payee, token, err := app.PayeeLogin(input.Email, input.Password)
	if err != nil {
		if errors.Is(err, core.PayeeInvalidPasswordError{}) {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"message": "Wrong e-mail or password",
			})
		}

		if errors.Is(err, core.PayeeNotFoundError{}) {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"message": "Wrong e-mail or password",
			})
		}
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"payee": payee,
		"token": token,
	})
}
