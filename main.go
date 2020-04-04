package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/ErrorBoi/zhirobot/bot"
	"github.com/ErrorBoi/zhirobot/db"
)

// MainHandler responds to http request
func MainHandler(resp http.ResponseWriter, _ *http.Request) {
	resp.Write([]byte("Hi there! I'm Zhirobot!"))
}

func main() {
	Zdb, err := db.NewDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer Zdb.DB.Close()

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
