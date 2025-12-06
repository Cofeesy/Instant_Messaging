package model

import (
	"gorm.io/gorm"
)

// 好友/群关系
// 添加好友的时候，ownerid和targetid都是userid
// 添加群的时候ownerid是userid，targetid是groupid
type Contact struct {
	gorm.Model
	// 谁的关系，是一个userid
	OwnerId uint `json:"ownerid"`
	// 和谁的关系,是一个userid吗？
	TargetId uint `json:"targetid"`
	// 关系是什么，好友？群关系？
	Relation int `json:"relation"` //1:好友 2:群聊
}

func (contact *Contact) TableName() string {
	return "contact"
}
