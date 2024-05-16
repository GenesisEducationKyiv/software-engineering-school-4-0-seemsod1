package driver

import (
	"fmt"
	"github.com/seemsod1/api-project/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type DB struct {
	SQL *gorm.DB
}

var dbConn = &DB{}

// openDb opens a database connection
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

	if err = testDB(d); err != nil {
		return nil, err
	}

	return dbConn, nil
}

func testDB(d *gorm.DB) error {
	err := d.Exec("SELECT 1").Error
	if err != nil {
		return err
	}
	return nil
}
