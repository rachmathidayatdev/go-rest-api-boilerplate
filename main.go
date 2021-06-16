package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-rest-api-boilerplate/application"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

// @title GO REST API DOCUMENTATION
// @version 1.0
// @description This is a documentation of API.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:9000
// @BasePath /
// @schemes http
func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	app := application.SetupApp()

	clientApp := cli.NewApp()
	clientApp.Name = "go-grpc-starter"
	clientApp.Version = "0.0.1"
	clientApp.Commands = []cli.Command{
		{
			Name:        "book",
			Description: "start book",
			Action: func(c *cli.Context) error {
				daemon := app.NewBookDaemon()
				return application.AppRunner(daemon)
			},
			Subcommands: []cli.Command{
				{
					Name:        "db-migrate",
					Usage:       "start book migration",
					Description: "start book migration",
					Action: func(c *cli.Context) error {
						source := os.Getenv("DB_MIGRATION_BOOK_SRC")
						log.Printf("Migrate database from source '%s'\n", source)

						migration, err := migrate.New(
							source,
							fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
								os.Getenv("DB_USERNAME"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), 5432, os.Getenv("DB_NAME")))

						if err != nil {
							fmt.Println(err)
						}

						defer migration.Close()

						if err := migration.Up(); err != nil {
							fmt.Println(err)
						}

						return nil
					},
				},
				{
					Name:        "db-rollback",
					Usage:       "start book rollback",
					Description: "start book rollback",
					Action: func(c *cli.Context) error {
						source := os.Getenv("DB_MIGRATION_BOOK_SRC")
						log.Printf("Rollback database from source '%s'\n", source)

						migration, err := migrate.New(
							source,
							fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
								os.Getenv("DB_USERNAME"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), 5432, os.Getenv("DB_NAME")))

						if err != nil {
							fmt.Println(err)
						}

						defer migration.Close()

						if err := migration.Down(); err != nil {
							fmt.Println(err)
						}

						return nil
					},
				},
			},
		},
	}

	clientApp.Run(os.Args)
}
