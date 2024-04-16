package main

import (
	"github.com/gin-gonic/gin"
)

var isRunning = false 

func main() {
	r := gin.Default()
	r.POST("/execute", execute)
	r.GET("/liveness", LivenessProbe)
	err := r.Run() // listen and serve on

	if err != nil {
		panic(err)
	}

}
