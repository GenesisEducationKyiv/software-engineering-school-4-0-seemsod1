package driver

import (
	"errors"
	"log"

	"github.com/kelseyhightower/envconfig"
	"github.com/seemsod1/api-project/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DatabaseConfig is a struct that holds the database configuration
type DatabaseConfig struct {
	DSN string
}

func NewDatabaseConfig() *DatabaseConfig {
	var dbConfig DatabaseConfig
	if err := envconfig.Process("DB", &dbConfig); err != nil {
		log.Fatal(err)
	}
	return &dbConfig
}

func (dc *DatabaseConfig) Validate() bool {
	return dc.DSN != ""
}

// GORMDriver is a struct that embeds the gorm.DB
type GORMDriver struct {
	SQL *gorm.DB
}

func NewGORMDriver() *GORMDriver {
	return &GORMDriver{}
}

// ConnectSQL connects to the database
func (gd *GORMDriver) ConnectSQL() (*GORMDriver, error) {
	cfg := NewDatabaseConfig()

	if !cfg.Validate() {
		return nil, errors.New("invalid database configuration")
	}

	d := gd.openDB(cfg)

	dbConn := &GORMDriver{}

	dbConn.SQL = d

	return dbConn, nil
}

// openDB opens a database connection
func (gd *GORMDriver) openDB(cfg *DatabaseConfig) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		log.Fatal("error connecting to database: ", err)
	}

	return db
}

func (gd *GORMDriver) RunMigrations() error {
	err := gd.SQL.AutoMigrate(&models.Subscriber{})
	if err != nil {
		return err
	}
	return nil
}
