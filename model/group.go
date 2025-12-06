package model

import (
	"gorm.io/gorm"
)

type Group struct {
	gorm.Model
	// 群主ID
	OwnerId uint `json:"ownerid"`
	// 群号
	GroupNumber int `json:"groupnumber"`
	// 群名
	GroupName string `json:"groupname"`
	// 群描述
	Description string `json:"description"`
	// 群头像
	Img string `json:"img"`
}

func (group *Group) TableName() string {
	return "group"
}

