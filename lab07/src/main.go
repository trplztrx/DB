package main

import (
	"lab07/app"
	"lab07/config"
	"log"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	app.RunApp(cfg)
}