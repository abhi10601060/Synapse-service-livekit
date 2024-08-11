package model

type StartStreamInput struct {
	Title     string `json:"title" binding:"required"`
	Desc      string `json:"desc"`
	Tags      string `json:"tags"`
	Thumbnail string `json:"thumbnail"`
	ToSave    bool   `json:"tosave omitempty" binding:"required"`
}

type Stream struct {
	Id             string `json:"id" gorm:"primarykey"`
	UserId         string `json:"userId" gorm:"not null"`
	Title          string `json:"title" gorm:"not null"`
	Desc           string `json:"desc"`
	Tags           string `json:"tags"`
	Status         string `json:"status" gorm:"not null"`
	SavedStreamUrl string `json:"savedStreamUrl"`
	Likes          int    `json:"likes"`
	Dislikes       int    `json:"dislikes"`
	ThumbnailUrl   string `json:"thumbNailUrl"`
	CreatedOn      string `json:"createdOn" gorm:"not null"`
	EndedOn        string `json:"endedOn"`
}
