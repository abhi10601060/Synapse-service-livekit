package main

import (
	"log"
	"net/http"
	"synapse/handler"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.GET("/ping", pong)

	stream := r.Group("/stream")
	{
		stream.POST("start", handler.CreateRoom)
		stream.POST("stop", handler.CloseRoom)
	}

	log.Fatal(r.Run(":8010"))
}

func pong(c *gin.Context){
	c.JSON(http.StatusOK, 
		gin.H{
			"message" : "Synapse is alive",
		})
}