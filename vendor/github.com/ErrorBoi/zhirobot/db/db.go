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

// Stat is a struct that contains weight date and value
type Stat struct {
	WeighDate   string
	WeightValue float64
}

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
	if ID == 0 {
		CreateUser(tgID)
		SetUserWeight(tgID, w)
	} else {
		_, err := Zdb.DB.Exec(`
		INSERT INTO userweight (user_id, weigh_date, weight_value)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, weigh_date) DO UPDATE 
		SET weight_value=EXCLUDED.weight_value`, ID, time.Now().Format("02/01/2006"), w)
		h.PanicIfErr(err)
	}
}

// UserExists checks if user with given tgID exists and returns it's ID or 0
func UserExists(tgID int) int {
	ID := 0

	stmt, err := Zdb.DB.Exec(`
	SELECT id FROM useracc
	WHERE tg_id = $1`, tgID)
	h.PanicIfErr(err)
	rowsAffected, err := stmt.RowsAffected()
	h.PanicIfErr(err)

	if rowsAffected != 0 {
		row := Zdb.DB.QueryRow(`
		SELECT id FROM useracc
		WHERE tg_id = $1`, tgID)
		err = row.Scan(&ID)
		h.PanicIfErr(err)
	}

	return ID
}

// GetUserWeight returns weight stats for chosen user
func GetUserWeight(tgID int) []*Stat {
	ID := UserExists(tgID)
	if ID == 0 {
		CreateUser(tgID)
		GetUserWeight(tgID)
	}
	rows, err := Zdb.DB.Query(`
	SELECT weigh_date, weight_value from userweight
	WHERE user_id = $1`, ID)
	h.PanicIfErr(err)
	defer rows.Close()

	stats := make([]*Stat, 0)
	for rows.Next() {
		stat := new(Stat)
		err := rows.Scan(&stat.WeighDate, &stat.WeightValue)
		h.PanicIfErr(err)
		stats = append(stats, stat)
	}
	err = rows.Err()
	h.PanicIfErr(err)
	return stats
}
