package models

import (
	"errors"
	"fmt"
	"gin_chat/models/system"
	"gin_chat/utils"
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
	// fmt.Println("BeforeCreate hook is being called!") // 添加日志
	u.UUID = uuid.New()
	fmt.Print(u.UUID)
	return
}

func GetUserList() ([]*User_Basic, error) {
	data := make([]*User_Basic, 10)
	// 记得这里是传地址
	// gorm操作并不熟悉
	if err := db.Find(&data).Error; err != nil {
		return nil, err
	}
	for _, v := range data {
		fmt.Println(v)
	}
	return data, nil
}

// 为什么需要这个,登陆的时候需要查找，或者其他的操作也可能需要查找
func FindUserByName(name string) (*User_Basic, error) {
	var user User_Basic
	err := db.Where("username = ?", name).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 没找到用户
			return nil, errors.New("用户不存在")
		}
		// 其他数据库错误
		return nil, err
	}
	return &user, nil
}

func FindUserByID(id uint) (*User_Basic, error) {
	var user User_Basic
	err := db.Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 没找到用户
			return nil, errors.New("用户不存在")
		}
		// 其他数据库错误
		return nil, err
	}
	return &user, nil
}

// FIXME:这里应该有问题
func FindUserByNameAndPassword(name, password string) (*User_Basic, error) {
	var user User_Basic
	err := db.Where("username = ?", name).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 没找到用户
			return nil, errors.New("用户不存在")
		}
		// 其他数据库错误
		return nil, err
	}

	if !utils.DecryptMD5(user.Salt, user.Password, password) {
		return nil, errors.New("密码输入错误")
	}
	return &user, nil
}

// 创建这个用户
func CreateUser(user_register *system.User_Register) error {
	var user User_Basic
	user.Username = user_register.Name
	user.Password = user_register.Password
	user.Salt = user_register.Salt
	// fmt.Println(user.Username)
	// fmt.Println(user.Password)
	result := db.Create(&user)
	return result.Error
}

// 更新用户信息名字和电话，邮箱
func UpdateUserInfo(updateuserinfo *system.UpdateUserInfo) error {
	var user User_Basic
	// 检查用户名是否已经被使用
	r := db.Where("username = ?", updateuserinfo.Username).First(&user).RowsAffected
	if r > 0 {
		return errors.New("该用户名已被使用")
	}

	// 根据id查找user
	err := db.Where("id = ?", updateuserinfo.ID).First(&user).Error
	if err != nil {
		return err
	}
	
	result := db.Model(&user).Updates(map[string]interface{}{"UserName": updateuserinfo.Username, "Phone": updateuserinfo.Phone, "Email": updateuserinfo.Email, "Icon": updateuserinfo.Icon})
	// 这个错误由db记录
	return result.Error
}

func UpdateUserPasswd(name, password string) error {
	var user User_Basic
	if err := db.Where("username = ?", name).First(&user).Error; err != nil {
		return err
	}
	password = utils.EncryptMD5(password, user.Salt)
	result := db.Model(&user).Updates(map[string]interface{}{"Password": password})
	// 这个错误由db记录
	return result.Error
}

// 删除用户
// 目前的逻辑是逻辑删除，也就是加上一个删除时间，但数据仍在数据库存在
// 是否是逻辑删除（Soft Delete），取决于 模型 User_Basic 是否包含了 GORM 的软删除字段。
// 即默认的gorm.Model包含DeletedAt字段
// 如果想要实际的物理删除
// 可以使用 Unscoped 方法，例如：db.Unscoped().Delete(&user)
func DeleteUser(name string) error {
	var user User_Basic
	if err := db.Where("username = ?", name).First(&user).Error; err != nil {
		return err
	}
	result := db.Delete(&user)
	// 这个错误由db记录
	return result.Error
}
