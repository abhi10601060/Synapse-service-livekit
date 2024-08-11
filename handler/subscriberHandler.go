package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"synapse/db"
	"synapse/util"
	"github.com/gin-gonic/gin"
)

func SubscribeStreamer(c *gin.Context){
	tokenStr := c.Request.Header.Get("Authentication-Token")

	subId := util.GetUserNameFromToken(tokenStr)
	if subId == "" {
		c.JSON(401, gin.H{
			"message": "In correct Header Token",
		})
		c.Abort()
		return
	}

	jsonMap := make(map[string]json.RawMessage)
	c.ShouldBind(&jsonMap)
	streamerId := strings.TrimSpace(string(jsonMap["streamerId"]))
	log.Println("AddSubscriber : streamer name fetched from body is : ", streamerId)
	streamerId = streamerId[1 : len(streamerId)-1]

	res := db.AddSubscriber(streamerId, subId)
	if  !res {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "unable to add the subscriber in db",
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successfully added subscriber in db",
	})
}

func UnsubscribeStreamer(c *gin.Context){
	tokenStr := c.Request.Header.Get("Authentication-Token")

	subId := util.GetUserNameFromToken(tokenStr)
	if subId == "" {
		c.JSON(401, gin.H{
			"message": "In correct Header Token",
		})
		c.Abort()
		return
	}

	jsonMap := make(map[string]json.RawMessage)
	c.ShouldBind(&jsonMap)
	streamerId := strings.TrimSpace(string(jsonMap["streamerId"]))
	log.Println("AddSubscriber : streamer name fetched from body is : ", streamerId)
	streamerId = streamerId[1 : len(streamerId)-1]

	res := db.RemoveSubscriber(streamerId, subId)
	if  !res {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "unable to remove the subscriber from db",
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successfully removed subscriber from db",
	})
}