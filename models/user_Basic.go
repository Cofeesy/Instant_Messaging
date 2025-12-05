package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 如果传进来的值中只包含一部份，像测试中我只会传name和password
// 而LoginInTime 就没有，如果是time.Time类型的话，系统会设置为默认值
// 而默认值是0001-01-01 00:00:00，这样mysql会报错
// 因此这里使用*time.Time指针类型

// omitempty是可选
type User_Basic struct {
	gorm.Model
	UUID          uuid.UUID  `json:"uuid" gorm:"type:varchar(36);comment:'对外暴露的用户ID'"`
	Username      string     `json:"username" gorm:"unique;not null" validate:"required"`
	Password      string     `json:"password" validate:"required,min=2,max=20"`
	Phone         string     `json:"phone" validate:"omitempty,len=11"`
	Email         string     `json:"email" validate:"omitempty,email"`
	Icon          string     `json:"icon"` //头像
	ClientIP      string     `json:"client_ip"`
	ClientPort    string     `json:"client_port"`
	Salt          string     `json:"salt"`
	LoginInTime   *time.Time `json:"login_in_time"`
	HeartbeatTime *time.Time `json:"heartbeat_time"`
	LoginOutTime  *time.Time `json:"login_out_time"`
	IsLoginOut    bool       `json:"is_login_out"`
	DeviceInfo    string     `json:"device_info"`
	LoginToken    string     `json:"token"`
}

func (user *User_Basic) TableName() string {
	return "user_basic"
}

// hook
func (u *User_Basic) BeforeCreate(tx *gorm.DB) (err error) {
	u.UUID = uuid.New()
	fmt.Print(u.UUID)
	return
}
