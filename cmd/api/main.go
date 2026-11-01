package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// fmt.Println(" ------ TODO API SERVER ------")
	// fmt.Println("Project is up")
	var router *gin.Engine = gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "TODO API is runnig!",
			"status":  "sucess",
		})
	})

	router.Run(":3000")
}
