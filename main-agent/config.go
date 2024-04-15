package main

import (
	"encoding/json"
	"os"
	"path"
)

type Configuration struct {
	Engine string `json:"engine"`
}

func ReadConfig(configPath string) (config Configuration, err error) {
	// Read the configuration file
	file, err := os.Open(path.Join(configPath))

	if err != nil {
		return
	}

	// Decode the configuration file
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)

	return
}
