package model

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	UserId     uint   `json:"userid"`
	TargetId   uint   `json:"targetid"`
	Type       int    `json:"type"`  // 1私聊 2群聊 3心跳
	Media      int    `json:"media"` // 1文字 2暂时没用 3语音 4图片，表情包
	Content    string `json:"content"`
	CreateTime uint64 `json:"createTime"`
	ReadTime   uint64 `json:"readTime"`
	Pic        string `json:"pic"`
	Url        string `json:"url"`
	Desc       string `json:"desc"`
	Amount     int    `json:"amount"`
}

func (message *Message) TableName() string {
	return "message"
}
