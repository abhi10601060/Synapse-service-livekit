package model

type User struct {
	Id       string `json:"id" binding:"required" gorm:"primarykey"`
	Password string `json:"password" binding:"required" gorm:"not null"`
}

type CreateUserInput struct {
	Bio            string `json:"bio"`
	ProfilePicture string `json:"profilePicture"`
}

type UserDetails struct {
	Id                string `json:"id" gorm:"primarykey"`
	Bio               string `json:"bio"`
	ProfilePictureUrl string `json:"profilePictureUrl"`
	TotalSubs         int    `json:"totalSubs"`
	TotalStreams      int    `json:"totalStreams"`
	CreatedOn         string `json:"createdOn" gorm:"not null"`
	Status            string `json:"status"`
}
