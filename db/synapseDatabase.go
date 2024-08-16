package db

import (
	"log"
	"synapse/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	DB_URL = "avnadmin:AVNS_3zGeELI4YI0-siKTItS@tcp(synapse-db-synapse.g.aivencloud.com:21929)/synapse_db?charset=utf8mb4&parseTime=True&loc=Local"
)

var (
	synapseDb *gorm.DB
)

func init(){
	synapseDb = connectToSynapseDb()

	if synapseDb != nil {
		autoMigrateModels()
	}
}

func connectToSynapseDb() *gorm.DB{
	db, err := gorm.Open(mysql.Open(DB_URL), &gorm.Config{})
	if err != nil {
		log.Println("connectToSynapseDb : error in connecting database : ", err)
		return nil
	}
	log.Println("Synapse database connected Successfully")
	return db
}

func autoMigrateModels(){
	err := synapseDb.AutoMigrate(&model.Stream{}, &model.Subscriber{}, &model.User{}, &model.UserDetails{})
	if err != nil {
		log.Println("error in autoMigrating models to database : ", err)
	}
}
