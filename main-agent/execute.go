package main

import (
	"bytes"
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
	var request ExecuteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	body, err := json.Marshal(request.Body)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	start := time.Now()
	status = Running

	resp, err := http.Post("http://localhost:3000/", "application/json", bytes.NewBuffer(body))

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

	status = Ready

	c.JSON(200, ExecuteResponse{
		Response: response,
		Time:     elapsed,
	})
}
