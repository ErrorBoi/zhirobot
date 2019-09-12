package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"zhirobot/bot"
	"zhirobot/db"
)

// Config contains project configuration
type Config struct {
	BotToken   string `json:"BotToken"`
	DBUser     string `json:"DBUser"`
	DBPassword string `json:"DBPassword"`
	DBAddress  string `json:"DBAddress"`
	DebugMode  string `json:"DebugMode"`
}

func main() {
	// Store Bot Configs in Config struct
	config := Config{}
	decode(&config)

	db.InitDB("postgres://" + config.DBUser + ":" + config.DBPassword + "@" + config.DBAddress + "/zhirobot")

	zBot, err := bot.InitBot(config.BotToken)
	if err != nil {
		log.Panic(err)
	}

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
