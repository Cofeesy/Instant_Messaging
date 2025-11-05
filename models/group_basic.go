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
func CreateGroup() {

}

// 加入群

// 通过群号查找群

// 通过群名查找群

// 加载群列表
