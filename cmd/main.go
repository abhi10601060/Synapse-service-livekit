package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.GET("/ping", pong)

	log.Fatal(r.Run(":8010"))
}

func pong(c *gin.Context){
	c.JSON(http.StatusOK, 
		gin.H{
			"message" : "Synapse is alive",
		})
}