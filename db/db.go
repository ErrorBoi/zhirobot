package db

import (
	"database/sql"
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

type User struct {
	ID        int
	TgID      int
	CreatedOn string
	Notify    bool
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
//func (db *DB) NewDatabase(dbName string) {
//	statement := fmt.Sprintf(`SELECT 'CREATE DATABASE %s'
//	WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '%s')`, dbName, dbName)
//	_, err := db.DB.Exec(statement)
//	if err != nil {
//		panic(err)
//	}
//}

// NewUserTable creates Table useracc
func (db *DB) NewUserTable() error {
	_, err := db.DB.Exec(`CREATE TABLE IF NOT EXISTS useracc(
		id serial PRIMARY KEY,
		tg_id INTEGER UNIQUE NOT NULL,
		created_on VARCHAR (50) NOT NULL,
		height INTEGER,
		age INTEGER,
		imt DECIMAL,
		notify BOOL DEFAULT FALSE
	 );`)
	return err
}

// NewWeightTable creates Table userweight
func (db *DB) NewWeightTable() error {
	_, err := db.DB.Exec(`CREATE TABLE IF NOT EXISTS userweight(
		id serial,
		user_id INTEGER REFERENCES useracc(id),
		weigh_date VARCHAR (50) NOT NULL,
		weight_value DECIMAL,
		PRIMARY KEY(user_id, weigh_date)
	 );`)
	return err
}

// CreateUser adds row to useracc Table if it doesn't already exist
func (db *DB) CreateUser(tgID int) error {
	_, err := db.DB.Exec(`
	INSERT INTO useracc (tg_id, created_on, notify)
	VALUES ($1, $2, false)
	ON CONFLICT (tg_id) DO NOTHING`, tgID, time.Now().Format("2006-01-02"))
	return err
}

// GetUser returns userID from user with given tgID
func (db *DB) GetUserID(tgID int) int {
	var userID int
	row := db.DB.QueryRow(`SELECT id FROM useracc
	WHERE tg_id = $1`, tgID)
	row.Scan(&userID)

	return userID
}

// SetUserWeight inserts/updates user weight, then returns difference between last 2 weights
func (db *DB) SetUserWeight(tgID int, weight float64) (*float64, error) {
	err := db.CreateUser(tgID)
	if err != nil {
		return nil, err
	}

	userID := db.GetUserID(tgID)
	_, err = db.DB.Exec(`
		INSERT INTO userweight (user_id, weigh_date, weight_value)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, weigh_date) DO UPDATE 
		SET weight_value=EXCLUDED.weight_value`, userID, time.Now().Format("2006-01-02"), weight)
	if err != nil {
		return nil, err
	}

	rows, err := db.DB.Query(`
	select weight_value from userweight
	where user_id = $1
	order by weigh_date desc
	limit 2`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		weightValues = make([]float64, 0)
		weightValue  float64
	)

	for rows.Next() {
		err := rows.Scan(&weightValue)
		if err != nil {
			return nil, err
		}
		weightValues = append(weightValues, weightValue)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	// if weightValues contains <2 values, user has 1 weight measure, so difference cannot be found
	if len(weightValues) < 2 {
		return &weight, nil
	}

	height, err := db.GetUserHeight(tgID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == nil && height != 0 {
		db.SetUserBMI(tgID)
	}

	weightDiff := weightValues[0] - weightValues[1]
	return &weightDiff, nil
}

// SetUserHeight updates user height
func (db *DB) SetUserHeight(tgID int, height int) error {
	err := db.CreateUser(tgID)
	if err != nil {
		return err
	}
	_, err = db.DB.Exec(`
		UPDATE useracc
		SET height = $1
		WHERE tg_id = $2
	`, height, tgID)
	if err != nil {
		return err
	}

	err = db.SetUserBMI(tgID)
	if err != nil {
		return err
	}

	return nil
}

// GetUserWeight returns weight stats for chosen user
func (db *DB) GetUserWeight(tgID, page int) ([]*Stat, *bool, error) {
	err := db.CreateUser(tgID)
	if err != nil {
		return nil, nil, err
	}

	userID := db.GetUserID(tgID)
	rows, err := db.DB.Query(`
	SELECT weigh_date, weight_value FROM userweight
	WHERE user_id = $1
	ORDER BY weigh_date DESC
	LIMIT 10
	OFFSET $2*10`, userID, page)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	stats := make([]*Stat, 0)
	for rows.Next() {
		stat := new(Stat)
		err := rows.Scan(&stat.WeighDate, &stat.WeightValue)
		if err != nil {
			return nil, nil, err
		}
		stats = append(stats, stat)
	}
	err = rows.Err()
	if err != nil {
		return nil, nil, err
	}

	var count int
	row := db.DB.QueryRow(`
	SELECT count(*) FROM userweight
	WHERE user_id = $1`, userID)

	err = row.Scan(&count)
	if err != nil {
		return nil, nil, err
	}

	last := true
	if count > (page+1)*10 {
		last = false
	}

	return stats, &last, nil
}

// GetUserHeight returns weight stats for chosen user
func (db *DB) GetUserHeight(tgID int) (int, error) {
	var height int
	row := db.DB.QueryRow(`SELECT height FROM useracc
	WHERE tg_id = $1`, tgID)

	if err := row.Scan(&height); err != nil {
		return 0, err
	}

	return height, nil
}

// SetUserBMI updates user BMI
func (db *DB) SetUserBMI(tgID int) error {
	err := db.CreateUser(tgID)
	if err != nil {
		return err
	}

	//get height
	height, err := db.GetUserHeight(tgID)
	if err != nil {
		return err
	}

	//get last weight
	weightValues, _, err := db.GetUserWeight(tgID, 0)
	if err != nil {
		return err
	}
	weight := weightValues[0].WeightValue

	//calc BMI
	var (
		bmiValue       float64
		heightInMeters float64
	)

	heightInMeters = float64(height) / 100
	bmiValue = weight / (heightInMeters * heightInMeters)

	//set BMI
	_, err = db.DB.Exec(`
		UPDATE useracc
		SET imt = $1
		WHERE tg_id = $2
	`, bmiValue, tgID)
	if err != nil {
		return err
	}
	return nil
}

// GetUserBMI gets user BMI
func (db *DB) GetUserBMI(tgID int) (float64, error) {
	var bmi float64
	row := db.DB.QueryRow(`SELECT imt FROM useracc
	WHERE tg_id = $1`, tgID)

	if err := row.Scan(&bmi); err != nil {
		return 0, err
	}

	return bmi, nil
}

func (db *DB) GetUsers() ([]User, error) {
	rows, err := db.DB.Query(`select id, tg_id, created_on, notify from useracc;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.ID, &user.TgID, &user.CreatedOn, &user.Notify)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (db *DB) SetNotify(tgID int, value bool) error {
	_, err := db.DB.Exec(`
	UPDATE useracc
	SET notify = $1
	WHERE tg_id = $2`, value, tgID)
	return err
}
