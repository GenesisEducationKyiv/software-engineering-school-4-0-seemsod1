package main

import (
	"log"

	"github.com/seemsod1/api-project/internal/storage/dbrepo"

	"github.com/joho/godotenv"
	"github.com/seemsod1/api-project/internal/config"
	"github.com/seemsod1/api-project/internal/driver"
	"github.com/seemsod1/api-project/internal/handlers"
	"github.com/seemsod1/api-project/internal/notifier"
)

// setup sets up the application
func setup(_ *config.AppConfig) error {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	dr := driver.NewGORMDriver()

	db, err := dr.ConnectSQL()
	if err != nil {
		log.Println("Cannot connect to database! Dying...")
		return err
	}

	if err = db.RunMigrations(); err != nil {
		log.Println("Cannot run schemas migration! Dying...")
		return err
	}

	dbRepository := dbrepo.NewGormRepo(db.DB)
	repo := handlers.NewRepo(dbRepository)

	handlers.NewHandlers(repo)

	notification := notifier.NewEmailNotifier(dbRepository)
	if err = notification.Start(); err != nil {
		log.Println("Cannot start mail sender! Dying...")
		return err
	}

	return nil
}
