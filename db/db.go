package db

import (
	"database/sql"
	"fmt"
	"log"
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
		created_on VARCHAR (50) NOT NULL,
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
		id serial,
		user_id INTEGER REFERENCES useracc(id),
		weigh_date VARCHAR (50) NOT NULL,
		weight_value DECIMAL,
		PRIMARY KEY(user_id, weigh_date)
	 );`)
	h.PanicIfErr(err)
	_, err = stmt.Exec()
	h.PanicIfErr(err)
}

// CreateUser adds row to useracc Table
func CreateUser(tgID int) {
	_, err := Zdb.DB.Exec(`
	INSERT INTO useracc (tg_id, created_on)
	VALUES ($1, $2)
	ON CONFLICT (tg_id) DO NOTHING`, tgID, time.Now().Format("02/01/2006"))
	h.PanicIfErr(err)
}

// SetUserWeight inserts/updates user weight
func SetUserWeight(tgID int, w float64) {
	ID := UserExists(tgID)
	if ID != 0 {
		_, err := Zdb.DB.Exec(`
		INSERT INTO userweight (user_id, weigh_date, weight_value)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, weigh_date) DO UPDATE 
		SET weight_value=EXCLUDED.weight_value`, ID, time.Now().Format("02/01/2006"), w)
		h.PanicIfErr(err)
	} else {
		log.Println("No user")
	}
}

// UserExists checks if user with given tgID existsand returns it's ID or 0
func UserExists(tgID int) int {
	ID := 0
	row := Zdb.DB.QueryRow(`
	SELECT id FROM useracc
	WHERE tg_id = $1`, tgID)
	err := row.Scan(&ID)
	h.PanicIfErr(err)
	return ID
}

// InsertWeight inserts user weight into userweight table
func InsertWeight() {
	// _, err := Zdb.DB.Exec(`
	// INSERT INTO useracc (tg_id, created_on)
	// VALUES ($1, $2::timestamp)
	// ON CONFLICT (tg_id) DO UPDATE SET created_on=EXCLUDED.created_on`, id, time.Now().Format(time.RFC3339))
	// h.PanicIfErr(err)
}
