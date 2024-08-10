package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"synapse/awshelper"
	"synapse/db"
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
	log.Println("room name fetched from body is : ", streamInput.Title)

	roomId := userName + "-" + uuid.NewString()
	log.Println("created room Id is  : ", roomId)

	thumbnailPath := userName + "/" + roomId
	thumbnailUrl, awsErr := awshelper.S3ServiceObject.UploadBase64(streamInput.Thumbnail, thumbnailPath)
	if awsErr != nil {
		log.Println("error during uploading thumbnail : " , err)
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"message" : "Error during uploading thumbnail.",
			})
		c.Abort()
		return
	}

	stream := model.Stream{
		Id: roomId,
		UserId: userName,
		Title: streamInput.Title,
		Desc: streamInput.Desc,
		Tags: streamInput.Tags,
		Status: "Live",
		ThumbnailUrl: thumbnailUrl,
		CreatedOn: time.Now().Format("01-02-2006 15:04:05"),
	}
	dbRes := db.AddStream(&stream)
	if !dbRes {
		log.Println("error in storing stream")
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"message" : "Error during stream storing to db",
			})
		c.Abort()
		return
	}

	room, err := roomClient.CreateRoom(context.Background(), &livekit.CreateRoomRequest{
		Name: roomId,
		EmptyTimeout: 1*60,
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
	roomId := strings.TrimSpace(string(jsonMap["room"]))
	log.Println("room name fetched from body is : ", roomId)
	roomId = roomId[1 : len(roomId)-1]
	
	var userNameSize = len(userName)
	log.Println("userName from Id : ", roomId[0:userNameSize])
	if userName != roomId[0:userNameSize]{
		log.Println("unauthorized streamer closing stream...")
		c.JSON(http.StatusNotAcceptable,
			gin.H{
				"message": "unauthorized stream owner closing stream",
			})
		c.Abort()
		return
	}

	res := db.ChangeStreamStatusToEnded(roomId)
	if !res {
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"message": "unable to change stream status into db",
			})
		c.Abort()
		return
	}

	_, err := roomClient.DeleteRoom(context.Background(), &livekit.DeleteRoomRequest{
		Room: roomId,
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
			"message": "stream closed successfully : " + roomId,
		})
}

func JoinRoomAsStreamer(c *gin.Context, roomId, userName string) {
	at := auth.NewAccessToken(API_KEY, SECRET)
	grant := &auth.VideoGrant{
		Room:      roomId,
		RoomAdmin: true,
		RoomJoin:  true,
	}
	at.AddGrant(grant).SetIdentity(userName).SetValidFor(24 * time.Hour)
	streamAccessToken, err := at.ToJWT()
	if err != nil {
		log.Println("Error in Jwt token gneration for room "+roomId+" : ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error during stream token generation",
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":   streamAccessToken,
		"streamId" : roomId,
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
	streams, isSuccessfull := db.GetAllActiveStreams()
	if !isSuccessfull {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message" : "Could not get streams from db",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message" : "live streams fetched successfully",
		"streams": streams,
	})
}
