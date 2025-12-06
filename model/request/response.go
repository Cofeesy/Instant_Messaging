package request

import (
	"github.com/google/uuid"
)

type User_Register struct {
	Name       string `json:"username" gorm:"unique;not null" validate:"required"`
	Password   string `json:"password" validate:"required,min=2,max=20"`
	Repassword string `json:"repassword"`
	Salt       string `json:"salt"`
}

type User_Login struct {
	Name     string `json:"name" gorm:"unique;not null" validate:"required"`
	Password string `json:"password" validate:"required,min=2,max=20"`
	Token    string `json:"token"`
}

type UpdateUserInfo struct {
	ID 		 uint   `json:"userid"`
	Username string `json:"username" gorm:"unique;not null" validate:"required"`
	Phone    string `json:"phone" validate:"omitempty,len=11"`
	Email    string `json:"email" validate:"omitempty,email"`
	Icon     string `json:"icon"`
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
	UserId uint `json:"userid"`
	FriendName string `json:"friendname"`
}

type FindFriend struct {
	UserId uint `json:"userid"`
	FriendId uint `json:"friendid"`
}

type CreatGroup struct {
	OwnerId   uint   `json:"ownerid"`
	GroupName string `json:"name"`
	Icon      string `json:"icon"`
	Memo string `json:"memo"`
}

type LoadGroups struct {
	UserId uint `json:"userid"`
}

type AddGroup struct {
	UserId    uint   `json:"userid"`
	GroupName string `json:"groupname"`
}

type AuthData struct {
	UserId uint   `json:"userid"`
	Token  string `json:"token"`
}

type AuthMessage struct {
	Cmd int `json:"cmd"`
	UserId uint   `json:"userid"`
	Token  string `json:"token"`
}

type SingleRedisPayload struct {
	UserId    uint  `json:"userid"`
	TargetId  uint  `json:"targetid"`
	Start     int64 `json:"start"`
	End       int64 `json:"end"`
	IsReverse bool  `json:"isreverse"`
}

type GroupRedisPayload struct {
	GroupId uint  `json:"groupId"`
	Start   int64 `json:"start"`
	End     int64 `json:"end"`
}

type AiRedisMsgPayload struct {
	UserId  uint  `json:"userid"`
	TargetId  uint  `json:"targetid"`
	Start   int64 `json:"start"`
	End     int64 `json:"end"`
}

type BaseClaims struct { 
    UUID        uuid.UUID
    ID          uint
    Username    string
}