package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ExecuteRequest struct {
	Body interface{} `json:"body"`
}

type StandardLanguageSpecificResponse struct {
	Status  int         `json:"status"`
	Body    interface{} `json:"body"`
	Headers interface{} `json:"headers"`
}

type ExecuteResponse struct {
	Response StandardLanguageSpecificResponse `json:"response"`
	Time     int64                            `json:"time"`
}

func execute(c *gin.Context) {

	start := time.Now()

	resp, err := http.Post("http://localhost:3000/", "application/json", c.Request.Body)

	elapsed := time.Since(start).Milliseconds()

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var response StandardLanguageSpecificResponse

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, ExecuteResponse{
		Response: response,
		Time:     elapsed,
	})
}
