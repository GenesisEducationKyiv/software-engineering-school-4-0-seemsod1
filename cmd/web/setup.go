package main

import (
	"github.com/joho/godotenv"
	"github.com/seemsod1/api-project/internal/config"
	"github.com/seemsod1/api-project/internal/driver"
	"github.com/seemsod1/api-project/internal/handlers"
	"log"
	"os"
)

func setup(app *config.AppConfig) error {
	env, err := loadEnv()
	if err != nil {
		return err
	}

	app.Env = env

	conn, err := driver.ConnectSQL(app.Env)
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}

	repo := handlers.NewRepo(app, conn)
	handlers.NewHandlers(repo)

	return nil
}

func loadEnv() (*config.EnvVariables, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")
	dbPort := os.Getenv("DB_PORT")

	return &config.EnvVariables{
		DBHost:     dbHost,
		DBUser:     dbUser,
		DBPassword: dbPass,
		DBName:     dbName,
		DBSSLMode:  dbSSLMode,
		DBPort:     dbPort,
	}, nil
}
