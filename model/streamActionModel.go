package model

type Subscriber struct {
	Uid              int    `gorm:"primarykey;autoIncrement:true;unique"`
	StreamerId       string `json:"streamerId" binding:"required" gorm:"size:191;not null"`
	SubscriberId     string `json:"subscriberId" binding:"required" gorm:"not null"`
	SubscriptionDate string `json:"subscriptionDate" gorm:"not null"`
}
