package driver

import (
	"fmt"
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
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable port=5432",
		env.DBHost, env.DBUser, env.DBName, env.DBPassword)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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
