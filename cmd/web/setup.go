package main

import (
	"fmt"

	logger "github.com/seemsod1/api-project/internal/logger"

	"github.com/seemsod1/api-project/internal/config"
	"github.com/seemsod1/api-project/internal/driver"
)

// setup sets up the application
func setup(_ *config.AppConfig, logg *logger.Logger) (*driver.GORMDriver, error) {
	dr := driver.NewGORMDriver()

	logg.Info("Connecting to database...")
	db, err := dr.ConnectSQL()
	if err != nil {
		logg.Error("Cannot connect to database! Dying...")
		return nil, fmt.Errorf("connecting to database: %w", err)
	}

	logg.Info("Running migrations...")
	if err = db.RunMigrations(); err != nil {
		logg.Error("Cannot run schemas migration! Dying...")
		return nil, fmt.Errorf("running schemas migration: %w", err)
	}
	return db, nil
}
