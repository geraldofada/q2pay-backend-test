package repository

import (
	"errors"
	"fmt"
	"os"

	"github.com/geraldofada/q2pay-backend-test/src/core"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	Conn *gorm.DB
}

func New() (*Repository, error) {
	env := os.Getenv("ENV")

	var conn string

	// NOTE: para produção adicionar as condições nesse if
	if env == "dev" {
		hostname := os.Getenv("DB_DEV_HOST")
		host_port := os.Getenv("DB_DEV_PORT")
		username := os.Getenv("DB_DEV_USER")
		password := os.Getenv("DB_DEV_PASS")
		databasename := os.Getenv("DB_DEV_DATABASE")

		conn =
			fmt.Sprintf(
				"sslmode=disable host=%s port=%s user=%s password=%s dbname=%s",
				hostname,
				host_port,
				username,
				password,
				databasename)
	}

	connection, err := gorm.Open(postgres.Open(conn), &gorm.Config{})

	return &Repository{Conn: connection}, err
}

func (r *Repository) CreateAccount(account *core.Account) error {
	result := r.Conn.Create(account)

	if result.Error != nil {
		if result.Error.(*pgconn.PgError).Code == "23505" {
			return core.AccountDuplicateError{}
		}
		return result.Error
	}
	return nil
}

func (r *Repository) GetAccountByEmail(email string) (core.Account, error) {
	foundAcc := core.Account{}

	result := r.Conn.Where("email = ?", email).First(&foundAcc)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return core.Account{}, core.AccountNotFoundError{}
		}

		return core.Account{}, result.Error
	}

	return foundAcc, nil
}

func (r *Repository) GetAccountByDoc(doc string) (core.Account, error) {
	foundAcc := core.Account{}

	result := r.Conn.Where("doc = ?", doc).First(&foundAcc)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return core.Account{}, core.AccountNotFoundError{}
		}

		return core.Account{}, result.Error
	}

	return foundAcc, nil
}

func (r *Repository) GetAccountById(id uint) (core.Account, error) {
	foundAcc := core.Account{}

	result := r.Conn.First(&foundAcc, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return core.Account{}, core.AccountNotFoundError{}
		}

		return core.Account{}, result.Error
	}

	return foundAcc, nil
}

// NOTE: essa função não lida com regra de negócio e apenas lida com erros que provém
// do banco de dados. Ou seja, é necessário checar quem pode enviar dinheiro, quem possui dinheiro
// para enviar etc etc. Atualmente isso ocorre no App que implementa o useCase de tranferência de
// dinheiro
//
// NOTE 2: é necessário apenas salvar os dados novamente pois a função existente no core
// "Account.TransferMoney" lida com modificar os valores da struct Account
func (r *Repository) SaveMoneyTransferBetweenAccounts(src *core.Account, target *core.Account) (bool, error) {
	err := r.Conn.Transaction(func(tx *gorm.DB) error {
		result := tx.Save(src)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return core.AccountNotFoundError{}
			}
			return result.Error
		}

		result = tx.Save(target)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return core.AccountNotFoundError{}
			}
			return result.Error
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, core.AccountNotFoundError{}
		}
		return false, err
	}

	return true, nil
}
