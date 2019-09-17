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

	err = db.NewDB(psqlInfo)
	defer db.Zdb.DB.Close()
	h.PanicIfErr(err)

	// Init Bot and listen to updates
	ZBot, err := bot.InitBot(config.BotToken)
	h.PanicIfErr(err)

	ZBot.SetDebugMode(strconv.ParseBool(config.DebugMode))

	ZBot.InitUpdates()
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
