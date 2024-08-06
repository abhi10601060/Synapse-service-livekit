package model

type StartStreamInput struct {
	Title     string `json:"title" binding:"required"`
	Desc      string  `json:"desc"`
	Tags      string  `json:tags`
	Thumbnail string  `json:thumbnail`
	ToSave    bool    `json:tosave omitempty binding:"required"`
}
