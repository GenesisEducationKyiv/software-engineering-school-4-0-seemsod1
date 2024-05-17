package driver

import (
	"fmt"
	"github.com/seemsod1/api-project/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

// DB is a struct that embeds the gorm.DB
type DB struct {
	SQL *gorm.DB
}

var dbConn = &DB{}

// openDB opens a database connection
func openDB(env *config.EnvVariables) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=%s port=%s",
		env.DBHost, env.DBUser, env.DBName, env.DBPassword, env.DBSSLMode, env.DBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("error connecting to database: ", err)
	}

	return db, nil
}

// ConnectSQL connects to the database
func ConnectSQL(env *config.EnvVariables) (*DB, error) {
	d, err := openDB(env)
	if err != nil {
		return nil, err
	}

	dbConn.SQL = d

	return dbConn, nil
}
