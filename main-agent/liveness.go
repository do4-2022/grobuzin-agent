package main

import (
	"github.com/gin-gonic/gin"
)

func LivenessProbe(c *gin.Context) {
	if isRunning {
		c.JSON(503, gin.H{
			"status": "Pending",
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "Ready",
	})
}