package driver

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/seemsod1/api-project/pkg/logger"
	"go.uber.org/zap"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const retryNumber = 10

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
	l  *logger.Logger
}

func NewGORMDriver(l *logger.Logger) *GORMDriver {
	return &GORMDriver{
		l: l,
	}
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

	d, er := gd.openDB(cfg)
	if er != nil {
		return nil, er
	}

	dbConn := &GORMDriver{}

	dbConn.DB = d

	return dbConn, nil
}

// openDB opens a database connection
func (gd *GORMDriver) openDB(cfg *DatabaseConfig) (*gorm.DB, error) {
	var counts int64
	var db *gorm.DB
	var err error
	for {
		db, err = gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
		if err != nil {
			gd.l.Error("Postgres not yet ready...", zap.Int64("attempt", counts), zap.Error(err))
			counts++
		} else {
			gd.l.Debug("connected to Postgres!")
			return db, nil
		}

		if counts > retryNumber {
			gd.l.Error("maximum retry attempts exceeded", zap.Error(err))
			return nil, err
		}

		gd.l.Debug("retrying to connect to Postgres...")
		time.Sleep(2 * time.Second)
	}
}
