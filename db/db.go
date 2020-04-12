package db

import (
	"database/sql"
	"fmt"
	"time"

	// Go PostgreSQL package
	_ "github.com/lib/pq"
)

type DB struct {
	DB *sql.DB
}

// Stat is a struct that contains weight date and value
type Stat struct {
	WeighDate   string
	WeightValue float64
}

// NewDB creates and returns Database
func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	return &DB{DB: db}, nil
}

// NewDatabase creates database with chosen name
func (db *DB) NewDatabase(dbName string) {
	statement := fmt.Sprintf(`SELECT 'CREATE DATABASE %s'
	WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '%s')`, dbName, dbName)
	_, err := db.DB.Exec(statement)
	if err != nil {
		panic(err)
	}
}

// NewUserTable creates Table useracc
func (db *DB) NewUserTable() {
	stmt, err := db.DB.Prepare(`CREATE TABLE IF NOT EXISTS useracc(
		id serial PRIMARY KEY,
		tg_id INTEGER UNIQUE NOT NULL,
		created_on VARCHAR (50) NOT NULL,
		height INTEGER,
		age INTEGER
	 );`)
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		panic(err)
	}
}

// NewWeightTable creates Table userweight
func (db *DB) NewWeightTable() {
	stmt, err := db.DB.Prepare(`CREATE TABLE IF NOT EXISTS userweight(
		id serial,
		user_id INTEGER REFERENCES useracc(id),
		weigh_date VARCHAR (50) NOT NULL,
		weight_value DECIMAL,
		PRIMARY KEY(user_id, weigh_date)
	 );`)
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		panic(err)
	}
}

// CreateUser adds row to useracc Table
func (db *DB) CreateUser(tgID int) {
	_, err := db.DB.Exec(`
	INSERT INTO useracc (tg_id, created_on)
	VALUES ($1, $2)
	ON CONFLICT (tg_id) DO NOTHING`, tgID, time.Now().Format("2006-01-02"))
	if err != nil {
		panic(err)
	}
}

// GetUser returns userID from user with given tgID
func (db *DB) GetUserID(tgID int) int {
	var userID int

	row := db.DB.QueryRow(`SELECT id FROM useracc
	WHERE tg_id = $1`, tgID)
	row.Scan(&userID)

	return userID
}

// SetUserWeight inserts/updates user weight, then returns new and old weight values
func (db *DB) SetUserWeight(tgID int, weight float64) (float64, float64) {
	//TODO: "create if not exists" should be one sql statement
	exists := db.UserExists(tgID)
	if !exists {
		db.CreateUser(tgID)
	}

	userID := db.GetUserID(tgID)
	_, err := db.DB.Exec(`
		INSERT INTO userweight (user_id, weigh_date, weight_value)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, weigh_date) DO UPDATE 
		SET weight_value=EXCLUDED.weight_value`, userID, time.Now().Format("2006-01-02"), weight)
	if err != nil {
		panic(err)
	}

	rows, err := db.DB.Query(`
	select weight_value from userweight
	where user_id = $1
	order by weigh_date desc
	limit 2;`, userID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var (
		weightValues = make([]float64, 2)
		weightValue  float64
	)

	for rows.Next() {
		err := rows.Scan(&weight)
		if err != nil {
			panic(err)
		}
		weightValues = append(weightValues, weightValue)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return weightValues[0], weightValues[1]
}

// UserExists checks if user with given tgID exists and returns it's ID or 0
func (db *DB) UserExists(tgID int) bool {
	var exists bool

	row := db.DB.QueryRow(`select exists(select 1 from useracc where tg_id=$1)`, tgID)

	err := row.Scan(&exists)
	if err != nil {
		panic(err)
	}

	return exists
}

// GetUserWeight returns weight stats for chosen user
func (db *DB) GetUserWeight(tgID int) []*Stat {
	exists := db.UserExists(tgID)
	if !exists {
		db.CreateUser(tgID)
	}

	userID := db.GetUserID(tgID)
	rows, err := db.DB.Query(`
	SELECT weigh_date, weight_value from userweight
	WHERE user_id = $1`, userID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	stats := make([]*Stat, 0)
	for rows.Next() {
		stat := new(Stat)
		err := rows.Scan(&stat.WeighDate, &stat.WeightValue)
		if err != nil {
			panic(err)
		}
		stats = append(stats, stat)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return stats
}
