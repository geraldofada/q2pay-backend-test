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
	Type     string `validate:"required,oneof='SELLER' 'COMMON'"`
	Password string `validate:"required,min=3,max=64"`
}

type loginAccountInput struct {
	Email    string `validate:"required,email,min=6,max=32"`
	Password string `validate:"required,min=3,max=64"`
}

type transferMoneyAccountInput struct {
	SourceEmailOrDoc string `validate:"required" json:"source_email_or_doc"`
	TargetEmailOrDoc string `validate:"required" json:"target_email_or_doc"`
	Amount           string `validate:"required,min=3,max=64"`
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

	account, err := app.AccountSignup(input.Name, input.Email, input.Password, input.Doc, core.AccountType(input.Type))
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
		"token":   token,
	})
}

func transferMoneyAccount(c *fiber.Ctx, app port.AccountUseCase) error {
	var input transferMoneyAccountInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errValidator := validateInput(input)
	if errValidator != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidator)
	}

	_, err := app.AccountTransferMoney(input.Amount, input.SourceEmailOrDoc, input.TargetEmailOrDoc)
	if err != nil {
		if errors.Is(err, core.AccountNotFoundError{}) {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"message": "Source or target account not found",
			})
		}

		if errors.Is(err, core.MoneyParseError{}) {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"message": "Amount should be like 'BRL 100,00'",
			})
		}

		if errors.Is(err, core.AccountNotEnoughBalanceError{}) {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"message": "Source account does not have enough money",
			})
		}

		if errors.Is(err, core.AccountSellerCannotTransferError{}) {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"message": "Seller account cannot transfer money",
			})
		}

		if errors.Is(err, core.MoneyMismatchCurrencyError{}) {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"message": "Currency mismatch",
			})
		}
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"source":     input.SourceEmailOrDoc,
		"target":     input.TargetEmailOrDoc,
		"transfered": input.Amount,
	})
}
