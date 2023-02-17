package main

import (
	"flag"
	"log"
	"os"

	"github.com/geraldofada/q2pay-backend-test/src/adapter/repository"
	"github.com/geraldofada/q2pay-backend-test/src/adapter/rest"
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
		repo.Conn.AutoMigrate(&core.Payee{})
		os.Exit(0)
	}

	app := app.New(repo)

	rest := rest.New()
	rest.SetupPayeeRoutes(app)

	log.Fatal(rest.Fiber.Listen(":8080"))
}
