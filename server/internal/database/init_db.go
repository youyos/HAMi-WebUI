package database

import (
	"database/sql"
	"github.com/go-kratos/kratos/v2/log"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type DatabaseConfig struct {
	Driver         string `yaml:"driver"`
	DataSourceName string `yaml:"dataSourceName"`
}

type Config struct {
	Database DatabaseConfig `yaml:"database"`
}

func InitDB(config *DatabaseConfig) {
	var err error
	switch config.Driver {
	case "mysql":
		db, err = sql.Open("mysql", config.DataSourceName)
	default:
		log.Errorf("Unsupported database driver: %s", config.Driver)
	}

	if err != nil {
		log.Errorf("Failed to connect to database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Info("Connected to database")
}

func GetDB() *sql.DB {
	return db
}
