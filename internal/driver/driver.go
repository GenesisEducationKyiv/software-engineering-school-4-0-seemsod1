package driver

import (
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"
	"github.com/seemsod1/api-project/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DatabaseConfig is a struct that holds the database configuration
type DatabaseConfig struct {
	DSN string `required:"true"`
}

func NewDatabaseConfig() (*DatabaseConfig, error) {
	var dbConfig DatabaseConfig
	if err := envconfig.Process("DB", &dbConfig); err != nil {
		return nil, fmt.Errorf("processing database config: %w", err)
	}

	return &dbConfig, nil
}

func (dc *DatabaseConfig) Validate() bool {
	return dc.DSN != ""
}

// GORMDriver is a struct that embeds the gorm.DB
type GORMDriver struct {
	DB *gorm.DB
}

func NewGORMDriver() *GORMDriver {
	return &GORMDriver{}
}

// ConnectSQL connects to the database
func (gd *GORMDriver) ConnectSQL() (*GORMDriver, error) {
	cfg, err := NewDatabaseConfig()
	if err != nil {
		return nil, fmt.Errorf("creating database configuration: %w", err)
	}

	if !cfg.Validate() {
		return nil, fmt.Errorf("validating database configuration: %t", cfg.Validate())
	}

	d := gd.openDB(cfg)

	dbConn := &GORMDriver{}

	dbConn.DB = d

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
	err := gd.DB.AutoMigrate(&models.Subscriber{})
	if err != nil {
		return fmt.Errorf("auto migrating: %w", err)
	}
	return nil
}
