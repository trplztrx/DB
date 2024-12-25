package main

import (
	"log"
	"os/exec"
	"time"
)

func main() {
	ticker := time.NewTicker(5 * time.Minute)
	// ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	log.Println("Starting data generation scheduler...")
	if err := runScript("/app/venv/bin/python3", "users_gen.py"); err != nil {
		log.Printf("Error: %v", err)
	} else {
		log.Println("Success: Data generated successfully.")
	}
	for {
		select {
		case <-ticker.C:
			log.Println("Triggering data generation script...")
			if err := runScript("python3", "users_gen.py"); err != nil {
				log.Printf("Error: %v", err)
			} else {
				log.Println("Success: Data generated successfully.")
			}
		}
	}
}

func runScript(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Output: %s", string(output))
		return err
	}
	log.Printf("Output: %s", string(output))
	return nil
}
