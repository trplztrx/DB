package main

import (
	"lab06/app"
	"lab06/config"
	"log"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	app.RunApp(cfg)
}