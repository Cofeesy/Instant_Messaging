package system

// import (
// 	"gin_chat/models"
// )

// 这个json要和前端一致
type User_Register struct {
	Name     string `json:"username" gorm:"unique;not null" validate:"required"`
	Password string `json:"password" validate:"required,min=2,max=20"`
	Repassword string `json:"repassword"`
	Salt     string `json:"salt"`
}

type User_Login struct {
	Name     string `json:"name" gorm:"unique;not null" validate:"required"`
	Password string `json:"password" validate:"required,min=2,max=20"`
	Token    string `json:"token"`
}

type UpdateUserInfo struct {
	Username string `json:"username" gorm:"unique;not null" validate:"required"`
	Phone    string `json:"phone" validate:"omitempty,len=11"`
	Email    string `json:"email" validate:"omitempty,email"`
}

type UpdateUserPasswd struct {
	Password    string `json:"password" validate:"required,min=2,max=20"`
	NewPassword string `json:"newpassword" validate:"required,min=2,max=20"`
}

// type DeleteUserPayload struct {
// }

type FindUser struct {
	UserId uint `json:"userid"`
}

type LoadFriendsPayload struct {
	UserId uint `json:"userid"`
}

type AddFriend struct {
	// 自己的id
	UserId uint `json:"userid"`
	// 要添加的好友的名字
	FriendName string `json:"friendname"`
}

type FindFriend struct {
	// 自己的id
	UserId uint `json:"userid"`
	// 好友的id
	FriendId uint `json:"friendid"`
}

type CreatGroup struct {
	OwnerId   uint   `json:"ownerid"`
	GroupName string `json:"name"`
	Icon      string `json:"icon"`
	// 群描述
	Memo string `json:"memo"`
}

// type FindGroup struct {
// }

type LoadGroups struct {
	UserId uint `json:"userid"`
}

type AddGroup struct {
	UserId    uint   `json:"userid"`
	GroupName string `json:"groupname"`
}

// type ChatPayload struct {
// 	UserId    uint   `json:"userid"`
// 	Token string `json:"token"`
// }


type AuthData struct{
	UserId    uint   `json:"userid"`
	Token string `json:"token"`
}

type AuthMessage struct{
	// 用来验证token
	Cmd int `json:"cmd"`
	// AuthData
	UserId    uint   `json:"userid"`
	Token string `json:"token"`
}
type RedisPayload struct{
	UserId uint   `json:"userid"`
	TargetId uint `json:"targetid"`
	Start int64 `json:"start"`
	End int64 `json:"end"`
	IsReverse bool `json:"isreverse"`
}