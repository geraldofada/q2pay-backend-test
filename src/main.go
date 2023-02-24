package main

import (
	"flag"
	"log"
	"os"

	"github.com/geraldofada/q2pay-backend-test/src/adapter/repository"
	"github.com/geraldofada/q2pay-backend-test/src/adapter/rest"
	"github.com/geraldofada/q2pay-backend-test/src/adapter/service"
	"github.com/geraldofada/q2pay-backend-test/src/app"
	"github.com/geraldofada/q2pay-backend-test/src/core"
	"github.com/joho/godotenv"
)

func main() {
	var migrate bool

	flag.BoolVar(&migrate, "migrate", false, "Run migrations using the .env database")
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		panic("Error loading .env")
	}

	repo, err := repository.New()
	if err != nil {
		panic("Error connecting to database")
	}

	if migrate {
		repo.Conn.AutoMigrate(&core.Account{})
		os.Exit(0)
	}

	appAccount := app.NewAppAccount(repo, service.Service{})
	appAuth := app.NewAppAuth()

	rest := rest.New()
	rest.SetupAccountRoutes(appAccount, appAuth)

	log.Fatal(rest.Fiber.Listen(":8080"))
}
