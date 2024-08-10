package db

import (
	"log"
	"synapse/model"
	"time"
)

func AddStream(stream *model.Stream) bool{
	res := synapseDb.Save(stream)

	if res.Error != nil {
		log.Println("AddStreamer : error in storing stream : ", res.Error.Error())
		return false
	}
	log.Println("AddStreamer : user added succesfully")
	return true
}

func ChangeStreamStatusToEnded(roomId string) bool{
	stream := model.Stream{
		Id : roomId,
	}
	res := synapseDb.Model(&stream).Updates(model.Stream{Status: "ended" , EndedOn: time.Now().Format("01-02-2006 15:04:05")})
	if res.Error != nil {
		log.Println("error in changing stream status to ended : ", res.Error.Error())
		return false
	}
	return true
}

func GetAllActiveStreams() ([]model.Stream, bool){
	var streams []model.Stream
	res := synapseDb.Where(&model.Stream{Status: "live"}).Find(&streams)
	if res.Error != nil {
		log.Println("GetAllActiveStreams : error in listing all active streams : ", res.Error.Error())
		return nil, false
	}
	log.Println("GetAllActiveStreams : listed active streams successFully : ", streams)
	return streams, true
}