package db

import (
	"database/sql"
)

var db *sql.DB

// InitDB inits database
func InitDB(dataSourceName string) (*sql.DB, error) {
	var err error
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
