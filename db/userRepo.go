package db

import (
	"log"
	"synapse/model"

	"gorm.io/gorm"
)

func AddUserDetails(userDetails *model.UserDetails) bool {
	res := synapseDb.Save(userDetails)
	if res.Error != nil {
		log.Println("AddUserDetails : error in storing user details : ", res.Error.Error())
		return false
	}
	return true
}

func GetUserDetails(userId string) (model.UserDetails, bool){
	userDetail := model.UserDetails{Id: userId}

	res := synapseDb.First(&userDetail)
	if res.Error != nil {
		log.Println("Get User Details : error in finding user: ", res.Error.Error())
		return userDetail, false  
	}
	return userDetail, true
}

func UpdateProfilePictureUrl(id, url string) bool {
	user := model.UserDetails{Id: id}
	res := synapseDb.Model(&user).Update("profile_picture_url", url)
	if res.Error != nil {
		log.Println("UpdateProfilePictureUrl : error in updating profile picture: ", res.Error.Error())
		return false
	}
	return true
}

func UpdateBio(id, bio string) bool{
	user := model.UserDetails{Id : id}
	res := synapseDb.Model(&user).Update("bio" , bio)
	if res.Error != nil {
		log.Println("UpdateBio : error in updating bio: ", res.Error.Error())
		return false
	}
	return true
}

func AddStreamCount(id string) bool{
	user := model.UserDetails{Id : id}
	res := synapseDb.Model(&user).Update("total_streams" , gorm.Expr("total_streams + ?", 1))
	if res.Error != nil {
		log.Println("AddStreamCount : error in adding stream count: ", res.Error.Error())
		return false
	}
	return true
}

func AddSubScriberCount(id string) bool{
	user := model.UserDetails{Id : id}
	res := synapseDb.Model(&user).Update("total_subs" , gorm.Expr("total_subs + ?", 1))
	if res.Error != nil {
		log.Println("AddSubScriberCount : error in adding sub count: ", res.Error.Error())
		return false
	}
	return true
}

func RemoveSubScriberCount(id string) bool{
	user := model.UserDetails{Id : id}
	res := synapseDb.Model(&user).Update("total_subs" , gorm.Expr("total_subs - ?", 1))
	if res.Error != nil {
		log.Println("RemoveSubScriberCount : error in removing sub count: ", res.Error.Error())
		return false
	}
	return true
}
