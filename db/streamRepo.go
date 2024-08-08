package db

import (
	"log"
	"synapse/model"
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