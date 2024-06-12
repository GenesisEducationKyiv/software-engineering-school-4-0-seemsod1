package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/seemsod1/api-project/internal/config"
	"github.com/seemsod1/api-project/internal/driver"
	"github.com/seemsod1/api-project/internal/handlers"
	mailSender "github.com/seemsod1/api-project/internal/mail-sender"
	"github.com/seemsod1/api-project/internal/models"
)

// setup sets up the application
func setup(app *config.AppConfig) error {
	env := loadEnv()

	app.Env = env

	conn, err := driver.ConnectSQL(app.Env)
	if err != nil {
		log.Println("Cannot connect to database! Dying...")
		return err
	}

	if err = runSchemasMigration(conn); err != nil {
		log.Println("Cannot run schemas migration! Dying...")
		return err
	}

	repo := handlers.NewRepo(app, conn)
	handlers.NewHandlers(repo)

	ms := mailSender.NewMailSender(app, conn)
	if err = ms.Start(); err != nil {
		log.Println("Cannot start mail sender! Dying...")
		return err
	}

	return nil
}

// loadEnv loads the environment variables
func loadEnv() *config.EnvVariables {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	dbURL := os.Getenv("DB_URL")

	return &config.EnvVariables{
		DSN: dbURL,
	}
}

// runSchemasMigration runs the schemas migration
func runSchemasMigration(db *driver.DB) error {
	err := db.SQL.AutoMigrate(&models.Subscriber{})
	if err != nil {
		return err
	}
	return nil
}
