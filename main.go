package main

import (
	"os"

	bot "github.com/NoTestsNoScrumMasters/Ben-Bot-Go/pkg/bot"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {
	godotenv.Load(".env")
	token := os.Getenv("TOKEN")
	bot.Run(token)
	log.Info("ShapiroHelperBot started")
}
