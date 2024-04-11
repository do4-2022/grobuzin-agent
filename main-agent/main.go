package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.POST("/execute", execute)
	err := r.Run() // listen and serve on

	if err != nil {
		panic(err)
	}

}
