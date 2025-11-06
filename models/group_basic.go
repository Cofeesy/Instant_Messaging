package models

type Group struct {
	// 群主
	OwnerId int `json:"ownerid"`
	// 群号
	GroupNumber int `json:"groupnumber"`
	// 群名
	GroupName string `json:"groupname"`
	// 群成员
}

func (group *Group) TableName() string {
	return "group"
}

// 创建群
func CreateGroup(ownerId int, groupName string) (*Group, error) {
	group := Group{
		OwnerId:   ownerId,
		GroupName: groupName,
	}
	if err := db.Create(&group).Error; err != nil {
		return nil, err
	}

	return &group, nil
}

// 加入群

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
