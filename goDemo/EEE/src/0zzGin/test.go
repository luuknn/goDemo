package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Hello tomorrow")
	})

	router.GET("/sjl", func(c *gin.Context) {
		c.String(200, "Hello sjl")
	})

	router.POST("/submit", func(c *gin.Context) {
		c.String(401, "not authorized")
	})

	router.PUT("/error", func(c *gin.Context) {
		c.String(500, "and error hapenned :(")
	})
	router.Run(":8088")

}
