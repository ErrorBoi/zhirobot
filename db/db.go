package db

import (
	"database/sql"
	"fmt"

	// Go PostgreSQL package
	h "zhirobot/helpers"

	_ "github.com/lib/pq"
)

type DB struct {
	DB *sql.DB
}

// NewDB creates and returns Database
func NewDB(dataSourceName string) (*DB, error) {
	var db DB
	var err error

	db.DB, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.DB.Ping(); err != nil {
		return nil, err
	}
	return &db, nil
}

func (db DB) NewDatabase(dbName string) {
	statement := fmt.Sprintf(`SELECT 'CREATE DATABASE %s'
	WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '%s')`, dbName, dbName)
	_, err := db.DB.Exec(statement)
	h.PanicIfErr(err)
}
