package models

import "github.com/jinzhu/gorm"

// 关系
type Contact struct {
	gorm.Model
	OwnerId  int `json:"ownerid"`
	FrendId  int `json:"frendid"`
	Relation int `json:"relation"` //0:好友 1:群聊
}

func (contact *Contact) TableName() string {
	return "contact"
}

// 返回指定的好友信息
func GetFrend(ownerid, frendid int) (*Contact, error) {
	var contact Contact
	if err := db.Where("owner_id=? AND frend_id=?", ownerid, frendid).Find(&contact).Error; err != nil {
		return nil, err
	}
	return &contact, nil
}

// 返回所有的好友信息
// 问题1：这个函数查找了两张表，gorm是怎么推断表的呢？我没有指定呢
func GetFrends(ownerid int) ([]*User_Basic, error) {
	contacts := make([]*Contact, 0)
	// 查找所有好友的关系表
	if err := db.Where("owner_id=?", ownerid).Find(&contacts).Error; err != nil {
		return nil, err
	}
	println(contacts)
	// 保存好友的id
	searchid := make([]int, 0)
	for _, v := range contacts {
		searchid = append(searchid, v.FrendId)
	}
	// 通过id查找user
	users := make([]*User_Basic, 0)
	if err := db.Where("id IN ?", searchid).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil

}

// 添加好友

// 删除好友
