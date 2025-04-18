/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"log"
	"os"

	"github.com/sumeshmurali/mandarin/cmd/app/cmd"
	"github.com/sumeshmurali/mandarin/internal/config"
	"github.com/sumeshmurali/mandarin/internal/orchestrator"
)

func main() {
	if os.Getenv("ENV_DOCKER") == "true" {
		// If running in docker, set the config file to /mandarin/config.yaml
		path := "/mandarin/config.yaml"
		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Println("Config file not found at path:", path)
			// If the file does not exist, exit with an error
			os.Exit(1)
		}
		c, err := config.ParseConfiguration(path)
		if err != nil {
			// If there is an error parsing the config file, exit with an error
			log.Println("Error parsing config file:", err)
			os.Exit(1)
		}
		orchestrator.Run(c)
	} else {
		cmd.Execute()
	}
}
