package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var db *sql.DB

func InitDB(config *DatabaseConfig) {
	var err error
	switch config.Driver {
	case "mysql":
		db, err = sql.Open("mysql", config.DataSourceName)
	default:
		log.Fatalf("Unsupported database driver: %s", config.Driver)
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to database")
}

func GetDB() *sql.DB {
	return db
}
