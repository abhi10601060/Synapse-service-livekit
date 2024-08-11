package db

import (
	"log"
	"synapse/model"
	"synapse/util"

	"gorm.io/gorm/clause"
)

func AddSubscriber(streamerId, subscriberId string) bool {
	subModel := model.Subscriber{
		StreamerId:       streamerId,
		SubscriberId:     subscriberId,
		SubscriptionDate: util.GetCurrentTime(),
	}

	res := synapseDb.Clauses(clause.OnConflict{
		Where:     clause.Where{Exprs: []clause.Expression{clause.Eq{Column: "streamer_id", Value: streamerId}, clause.Eq{Column: "subscriber_id", Value: subscriberId}}},
		DoNothing: true,
	}).Create(&subModel)
	if res.Error != nil {
		log.Println("AddSubscriber : error in adding sub to db : ", res.Error.Error())
		return false
	}
	return true
}

func RemoveSubscriber(streamerId, subscriberId string) bool {
	subModel := model.Subscriber{
		StreamerId:   streamerId,
		SubscriberId: subscriberId,
	}

	res := synapseDb.Where(&subModel).Delete(&subModel)
	if res.Error != nil {
		log.Println("RemoveSubscriber : error in removing : ", res.Error.Error())
		return false
	}
	return true
}
