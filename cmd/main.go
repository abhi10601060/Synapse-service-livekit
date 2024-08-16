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
		stream.GET("/join/:pid", handler.JoinRoomAsViewer)
		stream.GET("all", handler.GetAllActiveStreams)
		stream.POST("/like", handler.LikeStream)
		stream.POST("/remove-like", handler.RemoveLikeOfStream)
		stream.POST("/dislike", handler.DislikeStream)
		stream.POST("/remove-dislike", handler.RemoveDislikeOfStream)
	}

	user := r.Group("/user")
	{
		user.POST("/create", handler.CreateUser)
		user.POST("/update/profile-pic", handler.UpdateProfilePicture)
		user.POST("/update/bio", handler.UpdateBio)
		user.POST("/subscribe", handler.SubscribeStreamer)
		user.POST("/unsubscribe", handler.UnsubscribeStreamer)
		user.GET("/subscribers", handler.GetAllSubscribers)
		user.GET("/subscriptions", handler.GetAllSubscriptions)
		user.GET("/is-subscribed/:streamerId", handler.IsSubscribed)
	}

	log.Fatal(r.Run(":8010"))
}

func pong(c *gin.Context) {
	c.JSON(http.StatusOK,
		gin.H{
			"message": "Synapse is alive",
		})
}
