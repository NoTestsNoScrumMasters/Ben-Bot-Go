package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("TOKEN environment variable not set. Make sure it’s defined in GitHub Actions (env or secrets).")
	}

	guild := os.Getenv("GUILD")
	if guild == "" {
		log.Println("GUILD environment variable not set. Using default or skipping if not required.")
	}
	Run(token)
	log.Info("ShapiroHelperBot started")

}
