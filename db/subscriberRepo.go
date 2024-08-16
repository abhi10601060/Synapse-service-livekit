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

func GetAllSubscribers(streamerId string) ([]model.UserDetails, error) {
	subs := []model.UserDetails{}

	subQuery := synapseDb.Model(&model.Subscriber{}).Distinct("subscriber_id").Where("streamer_id = ?", streamerId)

	res := synapseDb.Where("id IN (?)", subQuery).Find(&subs)
	if res.Error != nil {
		log.Println("error in getting subs : ", res.Error.Error())
		return nil, res.Error
	}
	return subs, nil
}

func GetAllSubscribedStreamers(userId string) ([]model.UserDetails, error) {
	subs := []model.UserDetails{}

	subQuery := synapseDb.Model(&model.Subscriber{}).Distinct("streamer_id").Where("subscriber_id = ?", userId)

	res := synapseDb.Where("id IN (?)", subQuery).Find(&subs)
	if res.Error != nil {
		log.Println("error in getting subscribed streamers : ", res.Error.Error())
		return nil, res.Error
	}
	return subs, nil
}

func IsSubscribed(streamerId, userId string) (bool, error){
	var count int64

	res := synapseDb.Model(&model.Subscriber{}).Where("streamer_id = ?", streamerId).Where("subscriber_id = ?", userId).Count(&count)
	if res.Error != nil {
		log.Println("IsSubscribed: error in getting is subcribed: ", res.Error.Error())
		return false, res.Error
	}
	log.Println("count in isSubscribed func : " , count)
	return count > 0 , nil
}