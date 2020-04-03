package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/ErrorBoi/zhirobot/bot"
	"github.com/ErrorBoi/zhirobot/db"
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
	resp.Write([]byte("Hi there! I'm Zhirobot!"))
}

func main() {
	Zdb, err := db.NewDB(os.Getenv("DATABASE_URL"))
	defer Zdb.DB.Close()
	if err != nil {
		panic(err)
	}
	Zdb.NewDatabase("testdb")
	Zdb.NewUserTable()
	Zdb.NewWeightTable()

	// Init Bot and listen to updates
	ZBot, err := bot.InitBot(os.Getenv("BOT_TOKEN"), Zdb)
	if err != nil {
		panic(err)
	}

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
