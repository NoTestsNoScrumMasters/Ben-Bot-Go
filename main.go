package main

import (
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	bot "github.com/NoTestsNoScrumMasters/Ben-Bot-Go/tree/master/pkg/bot"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	token := os.Getenv("TOKEN")
	bot.Run(token, "1017181081721639002")
	log.Info("ShapiroHelperBot started by ")

}
