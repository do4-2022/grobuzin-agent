package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		configPath = "/etc/grobuzin/agent-config.json"
	}

	config, err := ReadConfig(configPath)

	if err != nil {
		log.Println("Error reading config file, using default config")
		config.Engine = "nodejs"
	}

	stopChan := make(chan int, 1)

	// start the engine
	go LaunchEngine(config, stopChan)

	r := gin.Default()
	r.POST("/execute", execute)
	err = r.Run() // listen and serve on 8080

	stopChan <- 1
	if err != nil {
		panic(err)
	}

}
