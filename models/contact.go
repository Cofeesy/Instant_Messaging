package models

import "github.com/jinzhu/gorm"

// 关系
type Contact struct {
	gorm.Model
	OwnerId int `json:"ownerid"`
	FrendId int `json:"frendid"`
	Relation int `json:"relation"` //0:好友 1:群聊
}

func (contact *Contact) TableName() string {
	return "contact"
}

// 返回指定的好友
func GetFrend(frendid string) (*User_Basic,error){
	var user User_Basic
	if err:=db.Where("frendid=?", frendid).First(&user).Error;err!=nil{
		return nil,err
	}
	return &user,nil
}

// 返回所有的好友
func GetFrends(ownerid string) []*User_Basic{

	// 查找所有好友的关系表

	// 保存好友的id

	// 通过id查找user
}

// 添加好友

// 删除好友