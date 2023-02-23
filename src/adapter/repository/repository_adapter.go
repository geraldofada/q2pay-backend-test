package repository

import (
	"fmt"
	"os"

	"github.com/geraldofada/q2pay-backend-test/src/core"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	Conn *gorm.DB
}

func New() (*Repository, error) {
	env := os.Getenv("ENV")

	var conn string

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

func (r *Repository) CreateAccount(account core.Account) error {
	result := r.Conn.Create(&account)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *Repository) GetAccountByEmail(email string) (core.Account, error) {
	return core.Account{}, nil
}

func (r *Repository) GetAccountByDoc(email string) (core.Account, error) { return core.Account{}, nil }
