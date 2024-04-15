package main

import (
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
		panic(err)
	}

	stopChan := make(chan int)

	// start the engine
	go LaunchEngine(config, stopChan)

	r := gin.Default()
	r.POST("/execute", execute)
	err = r.Run() // listen and serve on

	if err != nil {
		panic(err)
	}

}
