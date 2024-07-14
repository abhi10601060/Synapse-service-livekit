package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"synapse/util"

	"github.com/gin-gonic/gin"
	"github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go/v2"
)

const(
	HOST = "https://synapse-wj7x7eni.livekit.cloud"
	API_KEY = "APIFQLP5QAf8Qwh"
	SECRET = "SZOeMU0GetYpTFWXkVxGNLwSqGDHsJNf35S0GDeJF8uB"
)

var(
	roomClient lksdk.RoomServiceClient
)

func init(){
	roomClient = *lksdk.NewRoomServiceClient(HOST, API_KEY, SECRET)
}

func CreateRoom(c *gin.Context) {
	tokenStr := c.Request.Header.Get("Authentication-Token")
	 
	userName := util.GetUserNameFromToken(tokenStr)
	if userName == "" {
		c.JSON(401, gin.H{
			"message": "In correct Header Token",
		})
		c.Abort()
		return
	}

	jsonMap := make(map[string]json.RawMessage)
	c.ShouldBind(&jsonMap)
	roomName := strings.TrimSpace(string(jsonMap["room"]))
	log.Println("room name fetched from body is : ", roomName)
	roomName = roomName[1: len(roomName) -1]

	room, err := roomClient.CreateRoom(context.Background(), &livekit.CreateRoomRequest{
		Name: roomName,
		EmptyTimeout: 10*60,
		MaxParticipants: 80,
	})
	if err != nil {
		log.Println("error during room creation : " , err)
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"message" : "Error during room creation try again.",
			})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK,
	gin.H{
		"message" : "Room created successFully : " + room.Name, 
	})
}

func CloseRoom(c *gin.Context){
	tokenStr := c.Request.Header.Get("Authentication-Token")
	 
	userName := util.GetUserNameFromToken(tokenStr)
	if userName == "" {
		c.JSON(401, gin.H{
			"message": "In correct Header Token",
		})
		c.Abort()
		return
	}

	jsonMap := make(map[string]json.RawMessage)
	c.ShouldBind(&jsonMap)
	roomName := strings.TrimSpace(string(jsonMap["room"]))
	log.Println("room name fetched from body is : ", roomName)
	roomName = roomName[1: len(roomName) -1]

	_, err := roomClient.DeleteRoom(context.Background(), &livekit.DeleteRoomRequest{
		Room : roomName,
	})

	if err != nil {
		log.Println("error in cliosing stream : ", err)
		c.JSON(http.StatusInternalServerError,
		gin.H{
			"message" : "unable to close stream due to error",
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK,
		gin.H{
			"message" : "stream closed successfully : " + roomName,
		})
}