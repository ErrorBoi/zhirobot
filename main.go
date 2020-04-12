package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"go.uber.org/zap"

	"github.com/ErrorBoi/zhirobot/bot"
	"github.com/ErrorBoi/zhirobot/db"
)

// MainHandler responds to http request
func MainHandler(resp http.ResponseWriter, _ *http.Request) {
	resp.Write([]byte("Hi there! I'm Zhirobot!"))
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Errorf("init logger error: %w", err))
	}
	defer logger.Sync()

	sugar := logger.Sugar()

	Zdb, err := db.NewDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		sugar.Fatalf("Database connection error: %w", err)
	}
	defer Zdb.DB.Close()

	//Zdb.NewDatabase("testdb")
	Zdb.NewUserTable()
	Zdb.NewWeightTable()

	// Init Bot and listen to updates
	ZBot, err := bot.InitBot(os.Getenv("BOT_TOKEN"), Zdb, sugar)
	if err != nil {
		sugar.Infow("Bot init error",
			"botToken", os.Getenv("BOT_TOKEN"))
		sugar.Fatalf("Bot init error: %w", err)
	}

	debugMode, err := strconv.ParseBool(os.Getenv("DEBUG_MODE"))
	if err != nil {
		sugar.Errorf("Parse bool error: %w", err)
	}
	ZBot.BotAPI.Debug = debugMode

	http.HandleFunc("/", MainHandler)
	go http.ListenAndServe(":"+os.Getenv("PORT"), nil)

	ZBot.InitUpdates(os.Getenv("BOT_TOKEN"))
}
