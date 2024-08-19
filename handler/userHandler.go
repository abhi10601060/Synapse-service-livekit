package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"synapse/awshelper"
	"synapse/db"
	"synapse/model"
	"synapse/util"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context){
	tokenStr := c.Request.Header.Get("Authentication-Token")
	userId := util.GetUserNameFromToken(tokenStr)
	if userId == "" {
		c.JSON(401, gin.H{
			"message": "In correct Header Token",
		})
		c.Abort()
		return
	}

	incomingUserDetail := model.CreateUserInput{}
	c.ShouldBind(&incomingUserDetail)

	profilePictureUrl := ""
	if incomingUserDetail.ProfilePicture != "" {
		var err error
		profilePictureUrl, err = awshelper.S3ServiceObject.UploadBase64(incomingUserDetail.ProfilePicture, userId + "/" + "profile-picture")
		if err != nil {
			log.Println("error in uploading profilePicture : ", err)
		}
	}

	userDetail := model.UserDetails{
		Id: userId,
		Bio: incomingUserDetail.Bio,
		ProfilePictureUrl: profilePictureUrl,
		CreatedOn: time.Now().Format("01-02-2006 15:04:05"),
	}
	res := db.AddUserDetails(&userDetail)
	if !res{
		c.JSON(http.StatusInternalServerError, gin.H{
			"message" : "unable to add user details to db",
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message" : "streamer added succesfully to the db",
	})
}

func GetOwnProfileDetail(c *gin.Context){
	tokenStr := c.Request.Header.Get("Authentication-Token")
	userId := util.GetUserNameFromToken(tokenStr)
	if userId == "" {
		c.JSON(401, gin.H{
			"message": "In correct Header Token",
		})
		c.Abort()
		return
	}
	
	userDetails, isExist := db.GetUserDetails(userId)
	if !isExist{
		c.JSON(http.StatusOK, gin.H{
			"message" : "user does not exist",
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message" : "user found",
		"details" : userDetails,
	})
}

func GetOthersProfileDetail(c *gin.Context){
	userId := c.Param("userId") 
	
	userDetails, isExist := db.GetUserDetails(userId)
	if !isExist{
		c.JSON(http.StatusOK, gin.H{
			"message" : "user does not exist",
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message" : "user found",
		"details" : userDetails,
	})
}

func UpdateProfilePicture(c *gin.Context){
	tokenStr := c.Request.Header.Get("Authentication-Token")
	userId := util.GetUserNameFromToken(tokenStr)
	if userId == "" {
		c.JSON(401, gin.H{
			"message": "In correct Header Token",
		})
		c.Abort()
		return
	}

	jsonMap := make(map[string]json.RawMessage)
	err := c.ShouldBind(&jsonMap)
	if err != nil {
		log.Println("UpdateProfilePicture : error in binding profile picture : ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message" : "unable to bind json",
		})
		c.Abort()
		return
	}
	
	profilePicture := string(jsonMap["profilePicture"])
	profilePicture = profilePicture[1 : len(profilePicture)-1]

	profilePictureUrl, err := awshelper.S3ServiceObject.UploadBase64(profilePicture, userId + "/" + "profile-picture")
	if err != nil {
		log.Println("error in uploading profilePicture : ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message" : "unable to upload image to s3",
		})
		c.Abort()
		return
	}

	res := db.UpdateProfilePictureUrl(userId, profilePictureUrl)
	if !res{
		c.JSON(http.StatusInternalServerError, gin.H{
			"message" : "unable to add profilePictureUrl to db",
		})
		c.Abort()
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message" : "Profile picture added succesfully",
	})
}

func UpdateBio(c *gin.Context){
	tokenStr := c.Request.Header.Get("Authentication-Token")
	userId := util.GetUserNameFromToken(tokenStr)
	if userId == "" {
		c.JSON(401, gin.H{
			"message": "In correct Header Token",
		})
		c.Abort()
		return
	}

	jsonMap := make(map[string]json.RawMessage)
	c.ShouldBind(&jsonMap)

	bio := string(jsonMap["bio"])

	res := db.UpdateBio(userId, bio)
	if !res{
		c.JSON(http.StatusInternalServerError, gin.H{
			"message" : "unable to add profilePictureUrl to db",
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message" : "Bio added succesfully",
	})
}