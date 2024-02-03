package main

import (
	bot "github.com/NoTestsNoScrumMasters/Ben-Bot-Go/pkg/bot"
	log "github.com/sirupsen/logrus"
)

func main() {
	bot.BotToken = "" //todo
	bot.Run()
	log.Info("ShapiroHelperBot started")
}
