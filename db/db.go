package db

import (
	"database/sql"
	"fmt"
	"time"

	h "zhirobot/helpers"

	// Go PostgreSQL package
	_ "github.com/lib/pq"
)

// DB is a struct to unite *sql.DB and smth else :D
type DB struct {
	DB *sql.DB
}

// Zdb is a working Database
var Zdb DB

// NewDB creates and returns Database
func NewDB(dataSourceName string) error {
	var err error

	Zdb.DB, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		return err
	}
	if err = Zdb.DB.Ping(); err != nil {
		return err
	}
	return nil
}

// NewDatabase creates database with chosen name
func NewDatabase(dbName string) {
	statement := fmt.Sprintf(`SELECT 'CREATE DATABASE %s'
	WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '%s')`, dbName, dbName)
	_, err := Zdb.DB.Exec(statement)
	h.PanicIfErr(err)
}

// NewUserTable creates Table useracc
func NewUserTable() {
	stmt, err := Zdb.DB.Prepare(`CREATE TABLE IF NOT EXISTS useracc(
		id serial PRIMARY KEY,
		tg_id INTEGER UNIQUE NOT NULL,
		created_on TIMESTAMP NOT NULL,
		height INTEGER,
		age INTEGER
	 );`)
	h.PanicIfErr(err)
	_, err = stmt.Exec()
	h.PanicIfErr(err)
}

// NewWeightTable creates Table userweight
func NewWeightTable() {
	stmt, err := Zdb.DB.Prepare(`CREATE TABLE IF NOT EXISTS userweight(
		id serial PRIMARY KEY,
		user_id INTEGER REFERENCES useracc(id),
		weigh_date TIMESTAMP NOT NULL,
		weight_value DECIMAL
	 );`)
	h.PanicIfErr(err)
	_, err = stmt.Exec()
	h.PanicIfErr(err)
}

// CreateUser adds row to useracc Table
func CreateUser(id int) {
	// stmt := fmt.Sprintf(`INSERT INTO useracc (tg_id, created_on)
	// VALUES ($1, $2::timestamp)`, id, time.Now().Format(time.RFC3339))

	// _, err := Zdb.DB.Exec(
	// 	`IF NOT EXISTS (SELECT FROM useracc WHERE tg_id = '$1') THEN
	// 	INSERT INTO useracc (tg_id, created_on)
	// 	VALUES ($1, $2::timestamp)
	// 	END IF;`, id, time.Now().Format(time.RFC3339))
	_, err := Zdb.DB.Exec(`
	INSERT INTO useracc (tg_id, created_on)
	VALUES ($1, $2::timestamp)
	ON CONFLICT (tg_id) DO NOTHING`, id, time.Now().Format(time.RFC3339))
	h.PanicIfErr(err)
}

// InsertWeight inserts user weight into userweight table
func InsertWeight() {
	// _, err := Zdb.DB.Exec(`
	// INSERT INTO useracc (tg_id, created_on)
	// VALUES ($1, $2::timestamp)
	// ON CONFLICT (tg_id) DO UPDATE SET created_on=EXCLUDED.created_on`, id, time.Now().Format(time.RFC3339))
	// h.PanicIfErr(err)
}
