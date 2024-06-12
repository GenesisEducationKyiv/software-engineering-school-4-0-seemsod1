package driver

import (
	"log"

	"github.com/seemsod1/api-project/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is a struct that embeds the gorm.DB
type DB struct {
	SQL *gorm.DB
}

var dbConn = &DB{}

// openDB opens a database connection
func openDB(env *config.EnvVariables) *gorm.DB {
	db, err := gorm.Open(postgres.Open(env.DSN), &gorm.Config{})
	if err != nil {
		log.Fatal("error connecting to database: ", err)
	}

	return db
}

// ConnectSQL connects to the database
func ConnectSQL(env *config.EnvVariables) (*DB, error) {
	d := openDB(env)

	dbConn.SQL = d

	return dbConn, nil
}
