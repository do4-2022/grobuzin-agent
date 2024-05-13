package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

type Status int 

const (
	Ready Status = iota
	Running
	Down //TBD
)

var status Status = Down 

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
	r.GET("/readiness", func(c *gin.Context) { // defines if the server is ready to accept requests
		switch status {
			case Ready:
				c.JSON(200, gin.H{"status": "ok"})
			case Running:
				c.JSON(200, gin.H{"status": "running"})
			default:
				c.JSON(500, gin.H{"status": "down"})
		}
	})
	r.GET("/liveness", func(c *gin.Context) { // defines if the server CAN accept requests
		c.Data(204, gin.MIMEHTML, nil)
	})
	err = r.Run() // listen and serve on

	if err != nil {
		panic(err)
	}

}
