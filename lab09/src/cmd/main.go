package main

import (
	"lab09/config"
	"lab09/internal/app"
	"log"
)

func main() {
	log.Println("start app")
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	app.RunApp(cfg)
	log.Println("stop app")
}

// SELECT a.city, COUNT(o.id) AS total_deliveries
// FROM orders o
// JOIN addresses a ON o.delivery_address_id = a.id
// GROUP BY a.city
// ORDER BY total_deliveries DESC;
