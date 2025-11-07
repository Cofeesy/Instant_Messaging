package models

import (
	"errors"
	"gorm.io/gorm"
)

// 好友/群关系
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
func AddFrend(ownerid, frendid int) error {
	// 这是一个事务操作
	tx := db.Begin()
	//
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	// 查找是否存在
	var contact Contact
	// 这里用Find，避免First的没找到错误
	if tx.Where("owner_id=? AND frend_id=? AND relation=?", ownerid, frendid, 1).Find(&contact).RowsAffected != 0 {
		tx.Rollback()
		return errors.New("已添加过好友")
	}

	// 创建1
	contact1 := Contact{
		OwnerId:  ownerid,
		FrendId:  frendid,
		Relation: 1,
	}
	if err := tx.Create(&contact1).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 创建2
	contact2 := Contact{
		OwnerId:  frendid,
		FrendId:  ownerid,
		Relation: 1,
	}
	if err := tx.Create(&contact2).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// 删除好友
