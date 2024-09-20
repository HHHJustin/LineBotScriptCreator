package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/v8/linebot"
)

func ConnectLineBot() *linebot.Client {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	bot, err := linebot.New(os.Getenv("LINE_CHANNEL_SECRET"), os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}
	return bot
}
