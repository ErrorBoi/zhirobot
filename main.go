package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/ErrorBoi/zhirobot/bot"
	"github.com/ErrorBoi/zhirobot/db"
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

// MainHandler responds to http request
func MainHandler(resp http.ResponseWriter, _ *http.Request) {
	resp.Write([]byte("Hi there! I'm DndSpellsBot!"))
}

func main() {

	// Store Bot Configs in Config struct
	// config := Config{}
	// decode(&config)

	// CreateConfig()

	// Create DB and Tables
	// port, err := strconv.Atoi(config.DBPort)
	// psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=testdb sslmode=disable",
	// 	config.DBHost, port, config.DBUser, config.DBPassword)

	err := db.NewDB(os.Getenv("DATABASE_URL"))
	defer db.Zdb.DB.Close()
	h.PanicIfErr(err)
	db.NewDatabase("testdb")
	db.NewUserTable()
	db.NewWeightTable()

	// Init Bot and listen to updates
	ZBot, err := bot.InitBot(os.Getenv("BOT_TOKEN"))
	h.PanicIfErr(err)

	ZBot.SetDebugMode(strconv.ParseBool(os.Getenv("DEBUG_MODE")))

	http.HandleFunc("/", MainHandler)
	go http.ListenAndServe(":"+os.Getenv("PORT"), nil)

	ZBot.InitUpdates(os.Getenv("BOT_TOKEN"))
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
