package model

type CreateUserInput struct{
	Bio string `json:"bio"`
	ProfilePicture string `json:"profilePicture"`
}

type UserDetails struct{
	Id string `json:"id" gorm:"primarykey"`
	Bio string `json:"bio"`
	ProfilePictureUrl string `json:"profilePictureUrl"`
	TotalSubs int `json:"totalSubs"`
	TotalStreams int `json:"totalStreams"`
	CreatedOn string `json:"createdOn" gorm:"not null"`
}

