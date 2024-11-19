package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/util"
	_ "github.com/lib/pq"
)

func InitDatabase(config util.Configuration) (*sql.DB, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s",
		config.DB.Username, config.DB.Password, config.DB.Name, config.DB.Host)
	db, err := sql.Open("postgres", connStr)

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	return db, err
}
