package db

import (
	"log"
	"synapse/model"

	"gorm.io/gorm"
)

func LikeStream(streamId string) bool {
	stream := model.Stream{Id: streamId}

	res := synapseDb.Model(&stream).Update("likes", gorm.Expr("likes + ?", 1))
	if res.Error != nil {
		log.Println("LikeStream : error in increasing like count in db : ", res.Error.Error())
		return false
	}
	return true
}

func RemoveLikeOfStream(streamId string) bool {
	stream := model.Stream{Id: streamId}

	res := synapseDb.Model(&stream).Update("likes", gorm.Expr("likes - ?", 1))
	if res.Error != nil {
		log.Println("RemoveLikeOfStream : error in decreasing like count in db : ", res.Error.Error())
		return false
	}
	return true
}

func DislikeStream(streamId string) bool {
	stream := model.Stream{Id: streamId}

	res := synapseDb.Model(&stream).Update("dislikes", gorm.Expr("dislikes + ?", 1))
	if res.Error != nil {
		log.Println("DislikeStream : error in increasing dislike count in db : ", res.Error.Error())
		return false
	}
	return true
}

func RemoveDisLikeOfStream(streamId string) bool {
	stream := model.Stream{Id: streamId}

	res := synapseDb.Model(&stream).Update("dislikes", gorm.Expr("dislikes - ?", 1))
	if res.Error != nil {
		log.Println("RemoveDisLikeOfStream : error in decreasing dislike count in db : ", res.Error.Error())
		return false
	}
	return true
}