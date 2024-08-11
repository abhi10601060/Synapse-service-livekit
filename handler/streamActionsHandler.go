package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"synapse/db"

	"github.com/gin-gonic/gin"
)

func LikeStream(c *gin.Context){
	jsonMap := make(map[string]json.RawMessage)
	c.ShouldBind(&jsonMap)
	streamId := strings.TrimSpace(string(jsonMap["streamId"]))
	log.Println("LikeStream : stream name fetched from body is : ", streamId)
	streamId = streamId[1 : len(streamId)-1]

	res := db.LikeStream(streamId)
	if !res{
		c.JSON(http.StatusInternalServerError, gin.H{
			"message" : "error in liking Stream",
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message" : "successfully liked stream: " + streamId,
	})
}

func RemoveLikeOfStream(c *gin.Context){
	jsonMap := make(map[string]json.RawMessage)
	c.ShouldBind(&jsonMap)
	streamId := strings.TrimSpace(string(jsonMap["streamId"]))
	log.Println("RemoveLikeStream : stream name fetched from body is : ", streamId)
	streamId = streamId[1 : len(streamId)-1]

	res := db.RemoveLikeOfStream(streamId)
	if !res{
		c.JSON(http.StatusInternalServerError, gin.H{
			"message" : "error in removing like of Stream",
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message" : "successfully removed like of stream: " + streamId,
	})
}


func DislikeStream(c *gin.Context){
	jsonMap := make(map[string]json.RawMessage)
	c.ShouldBind(&jsonMap)
	streamId := strings.TrimSpace(string(jsonMap["streamId"]))
	log.Println("LikeStream : stream name fetched from body is : ", streamId)
	streamId = streamId[1 : len(streamId)-1]

	res := db.DislikeStream(streamId)
	if !res{
		c.JSON(http.StatusInternalServerError, gin.H{
			"message" : "error in disliking Stream",
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message" : "successfully disliked stream: " + streamId,
	})
}

func RemoveDislikeOfStream(c *gin.Context){
	jsonMap := make(map[string]json.RawMessage)
	c.ShouldBind(&jsonMap)
	streamId := strings.TrimSpace(string(jsonMap["streamId"]))
	log.Println("RemoveDislikeOfStream : stream name fetched from body is : ", streamId)
	streamId = streamId[1 : len(streamId)-1]

	res := db.RemoveDisLikeOfStream(streamId)
	if !res{
		c.JSON(http.StatusInternalServerError, gin.H{
			"message" : "error in removing dislike of Stream",
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message" : "successfully removed dislike of stream: " + streamId,
	})
}
