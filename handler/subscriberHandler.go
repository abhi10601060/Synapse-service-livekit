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

func SubscribeStreamer(c *gin.Context) {
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
	if !res {
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

func UnsubscribeStreamer(c *gin.Context) {
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
	if !res {
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

func GetAllSubscribers(c *gin.Context) {
	tokenStr := c.Request.Header.Get("Authentication-Token")
	userId := util.GetUserNameFromToken(tokenStr)
	if userId == "" {
		c.JSON(401, gin.H{
			"message": "In correct Header Token",
		})
		c.Abort()
		return
	}

	subs, err := db.GetAllSubscribers(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to fetch subscribers for: " + userId,
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Successfully fetched subscribers for: " + userId,
		"subscribers": subs,
	})
}

func GetAllSubscriptions(c *gin.Context) {
	tokenStr := c.Request.Header.Get("Authentication-Token")
	userId := util.GetUserNameFromToken(tokenStr)
	if userId == "" {
		c.JSON(401, gin.H{
			"message": "In correct Header Token",
		})
		c.Abort()
		return
	}

	subscriptions, err := db.GetAllSubscribedStreamers(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to fetch subscribed streamers for: " + userId,
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Successfully fetched subscribed streamers for: " + userId,
		"subscribers": subscriptions,
	})
}

func IsSubscribed(c *gin.Context) {
	tokenStr := c.Request.Header.Get("Authentication-Token")
	userId := util.GetUserNameFromToken(tokenStr)
	if userId == "" {
		c.JSON(401, gin.H{
			"message": "In correct Header Token",
		})
		c.Abort()
		return
	}

	streamerId := c.Param("streamerId")

	res, err := db.IsSubscribed(streamerId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to check is subscribed streamer for: " + userId,
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully fetched is subscribed streamer for: " + userId,
		"status":  res,
	})
}
