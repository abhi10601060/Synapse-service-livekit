package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"synapse/awshelper"
	"synapse/model"
	"synapse/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/livekit/protocol/auth"
	"github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go/v2"
)

const (
	HOST    = "https://synapse-wj7x7eni.livekit.cloud"
	API_KEY = "APIFQLP5QAf8Qwh"
	SECRET  = "SZOeMU0GetYpTFWXkVxGNLwSqGDHsJNf35S0GDeJF8uB"
)

var (
	roomClient lksdk.RoomServiceClient
)

func init() {
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

	streamInput := model.StartStreamInput{}
	err := c.BindJSON(&streamInput)
	if err != nil {
		log.Println("bad request in start stream : ", err)
		c.Abort()
		return
	}
	log.Println("received StreamInput: ", streamInput)
	roomTitle := streamInput.Title
	log.Println("room name fetched from body is : ", roomTitle)

	roomId := userName + "-" + uuid.NewString()
	log.Println("created room Id is  : ", roomId)

	awsErr := awshelper.S3ServiceObject.UploadBase64(streamInput.Thumbnail, roomTitle)
	if awsErr != nil {
		log.Println("error during uploading thumbnail : " , err)
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"message" : "Error during uploading thumbnail.",
			})
		c.Abort()
		return
	}

	room, err := roomClient.CreateRoom(context.Background(), &livekit.CreateRoomRequest{
		Name: roomId,
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

	JoinRoomAsStreamer(c, room.Name, userName)
}

func CloseRoom(c *gin.Context) {
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
	roomName = roomName[1 : len(roomName)-1]

	_, err := roomClient.DeleteRoom(context.Background(), &livekit.DeleteRoomRequest{
		Room: roomName,
	})

	if err != nil {
		log.Println("error in cliosing stream : ", err)
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"message": "unable to close stream due to error",
			})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK,
		gin.H{
			"message": "stream closed successfully : " + roomName,
		})
}

func JoinRoomAsStreamer(c *gin.Context, roomName, userName string) {
	at := auth.NewAccessToken(API_KEY, SECRET)
	grant := &auth.VideoGrant{
		Room:      roomName,
		RoomAdmin: true,
		RoomJoin:  true,
	}
	at.AddGrant(grant).SetIdentity(userName).SetValidFor(24 * time.Hour)
	streamAccessToken, err := at.ToJWT()
	if err != nil {
		log.Println("Error in Jwt token gneration for room "+roomName+" : ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error during stream token generation",
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":   streamAccessToken,
		"message": "Room created successfully with streamer",
	})
}

func JoinRoomAsViewer(c *gin.Context) {
	tokenStr := c.Request.Header.Get("Authentication-Token")

	userName := util.GetUserNameFromToken(tokenStr)

	pid, isExist := c.Params.Get("pid")
	if !isExist {
		log.Println("Pid is empty for joining")
		c.Abort()
		return
	}
	pid = strings.TrimSpace(pid)

	at := auth.NewAccessToken(API_KEY, SECRET)
	grant := &auth.VideoGrant{
		Room:       pid,
		RoomCreate: false,
		RoomAdmin:  false,
		RoomJoin:   true,
	}
	at.AddGrant(grant).SetIdentity(userName).SetValidFor(24 * time.Hour)

	streamAccessToken, err := at.ToJWT()
	if err != nil {
		log.Println("Error in Jwt token gneration for room "+pid+" : ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error during stream token generation",
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":   streamAccessToken,
		"message": "stream join token generated successfully for : " + pid,
	})
}

func GetAllActiveStreams(c *gin.Context) {
	rooms, err := roomClient.ListRooms(context.Background(), &livekit.ListRoomsRequest{})
	if err != nil {
		log.Println("Error in fetching list of rooms : ", err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"rooms": rooms.Rooms,
	})
}
