package models

import (
	"errors"
	"gin_chat/models/system"

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

// 创建群
// 需要创建关系表
func CreateGroup(sysgroup system.CreatGroup) (*Group, error) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建群
	group := Group{
		OwnerId: sysgroup.OwnerId,
		// 群号应该自动生成
		GroupNumber: 11,
		GroupName:   sysgroup.GroupName,
		Description: sysgroup.Memo,
		Img:         sysgroup.Icon,
	}

	// 查找该群是否存在
	_, err := FindGroupByName(sysgroup.GroupName)
	if err == nil {
		return nil, errors.New("该群名已经被使用")
	}
	if err := db.Create(&group).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 创建关系,这里链接群id
	contact := Contact{
		OwnerId:  sysgroup.OwnerId,
		TargetId: group.ID,
		// 关系2是群关系
		Relation: 2,
	}
	if err := db.Create(&contact).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	return &group, tx.Commit().Error
}

// 加入群
// owner_id:userid
func AddGroup(addGroup *system.AddGroup) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 搜索群
	group, err := FindGroupByName(addGroup.GroupName)
	if err != nil {
		return errors.New("群不存在")
	}

	// 搜索关系
	var cont Contact
	row := tx.Where("owner_id=? AND target_id=? AND relation=?", addGroup.UserId, group.ID, 2).First(&cont).RowsAffected
	if row != 0 {
		return errors.New("你已经加过该群")
	}
	// 创建关系
	contact := Contact{
		OwnerId:  addGroup.UserId,
		TargetId: group.ID,
		Relation: 2,
	}
	if err := db.Create(&contact).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 加入群
	// group_new:=Group{
	// 	OwnerId: ,
	// }
	// if err := db.Create(&contact).Error; err != nil {
	// 	tx.Rollback()
	// 	return err
	// }

	return tx.Commit().Error
}

// 通过群号查找群

// 通过群名查找群
func FindGroupByName(groupName string) (*Group, error) {
	var group Group
	if err := db.Where("group_name = ?", groupName).First(&group).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

// 加载群列表
// 返回当前用户的群以及所在的群列表
// 这里都是查找，不用事务
func FindGroupsByUserID(ownerid uint) ([]*Group, error) {
	// 查找关系
	contacts := make([]*Contact, 0)
	err := db.Where("owner_id = ? AND relation=?", ownerid, 2).Find(&contacts).Error
	if err != nil {
		return nil, err
	}

	if len(contacts) == 0 {
		return nil, errors.New("你未加入任何群")
	}

	groups := make([]*Group, 0)
	groupid := make([]uint, 0)
	for _, v := range contacts {
		groupid = append(groupid, v.TargetId)
	}

	// 查找group
	if err := db.Where("id IN ?", groupid).Find(&groups).Error; err != nil {
		return nil, err
	}

	return groups, nil

}
