package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"zhirobot/bot"
	"zhirobot/db"
	h "zhirobot/helpers"
)

// Config contains project configuration
type Config struct {
	BotToken   string `json:"BotToken"`
	DBUser     string `json:"DBUser"`
	DBPassword string `json:"DBPassword"`
	DBHost     string `json:"DBHost"`
	DBPort     string `json:"DBPort"`
	DebugMode  string `json:"DebugMode"`
}

func main() {
	// Store Bot Configs in Config struct
	config := Config{}
	decode(&config)

	// Create DB and Test Table inside of it
	port, err := strconv.Atoi(config.DBPort)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=testdb sslmode=disable",
		config.DBHost, port, config.DBUser, config.DBPassword)

	zDB, err := db.NewDB(psqlInfo)
	h.PanicIfErr(err)
	defer zDB.DB.Close()
	zDB.NewDatabase("testdb")
	// _, err = db.Exec("USE testdb")
	// h.PanicIfErr(err)
	stmt, err := zDB.DB.Prepare(`CREATE TABLE useracc(
		user_id serial PRIMARY KEY,
		username VARCHAR (50) UNIQUE NOT NULL,
		password VARCHAR (50) NOT NULL,
		email VARCHAR (355) UNIQUE NOT NULL,
		created_on TIMESTAMP NOT NULL,
		last_login TIMESTAMP
	 );`)
	h.PanicIfErr(err)
	_, err = stmt.Exec()
	h.PanicIfErr(err)

	// Init Bot and listen to updates
	zBot, err := bot.InitBot(config.BotToken)
	h.PanicIfErr(err)

	zBot.SetDebugMode(strconv.ParseBool(config.DebugMode))

	zBot.InitUpdates()
}

func decode(c *Config) {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(c)
	if err != nil {
		fmt.Println("error:", err)
	}
}
