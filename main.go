package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	BotToken string `json:"BotToken"`
}

func main() {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	config := Config{}
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(config.BotToken)
}
