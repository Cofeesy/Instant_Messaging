package models

import (
	// "errors"
	"errors"
	"gin_chat/models/system"

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

// TODO:
// 返回指定的好友信息
func FindFrend(ownerid, targetid uint) (*Contact, error) {
	var contact Contact
	// TODO:先查找关系

	// TODO:再返回具体好友
	err := db.Where("owner_id=? AND target_id=? AND relation=?", ownerid, targetid, 1).First(&contact).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("你尚未添加该好友")
		}
		return nil, err
	}
	// // TODO:这里不是返回关系
	return &contact, nil
}

// 返回所有的好友信息
// 问题1：这个函数查找了两张表，gorm是怎么推断表的呢？我没有指定呢
func FindFriendsByUserID(ownerid uint) ([]*User_Basic, error) {
	contacts := make([]*Contact, 0)
	// 查找所有好友的关系表
	if err := db.Where("owner_id=? AND relation=?", ownerid, 1).Find(&contacts).Error; err != nil {
		return nil, err
	}
	println(contacts)
	// 保存好友的id
	searchid := make([]uint, 0)
	for _, v := range contacts {
		searchid = append(searchid, v.TargetId)
	}

	if len(searchid) == 0 {
		return nil, errors.New("你尚未添加任何好友")
	}

	// 通过id查找user
	users := make([]*User_Basic, 0)
	if err := db.Where("id IN ?", searchid).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil

}

// 添加好友
func AddFrend(addPayload *system.AddFriend) error {
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

	// 这里用Find，避免First的没找到错误
	// if tx.Where("owner_id=? AND frend_id=? AND relation=?", ownerid, targetid, 1).Find(&contact).RowsAffected != 0 {
	// 	tx.Rollback()
	// 	return errors.New("已添加过好友")
	// }

	// 查找是否存在该用户
	var friend *User_Basic
	friend, err := FindUserByName(addPayload.FriendName)
	if err != nil {
		return errors.New("不存在该用户")
	}

	// 查找是否存在关系
	_, err = FindFrend(addPayload.UserId, friend.ID)
	// 不等于nil说明已经找到
	if err == nil {
		return errors.New("已经添加过该好友")
	}

	// 双向创建关系1
	contact1 := Contact{
		OwnerId:  addPayload.UserId,
		TargetId: friend.ID,
		Relation: 1,
	}
	if err := tx.Create(&contact1).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 双向创建关系2
	contact2 := Contact{
		OwnerId:  friend.ID,
		TargetId: addPayload.UserId,
		Relation: 1,
	}
	if err := tx.Create(&contact2).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// 删除好友
