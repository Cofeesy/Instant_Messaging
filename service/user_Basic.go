package service

import (
	"errors"
	"ZustChat/model"
	"ZustChat/model/request"
	"ZustChat/utils"

	"gorm.io/gorm"
)

// 如果传进来的值中只包含一部份，像测试中我只会传name和password
// 而LoginInTime 就没有，如果是time.Time类型的话，系统会设置为默认值
// 而默认值是0001-01-01 00:00:00，这样mysql会报错
// 因此这里使用*time.Time指针类型

func GetUserList() ([]*model.User_Basic, error) {
	data := make([]*model.User_Basic, 10)
	if err := utils.DB.Find(&data).Error; err != nil {
		return nil, err
	}
	// 移除调试打印，如需调试可使用日志
	return data, nil
}

func FindUserByName(name string) (*model.User_Basic, error) {
	var user model.User_Basic
	err := utils.DB.Where("username = ?", name).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

func FindUserByID(id uint) (*model.User_Basic, error) {
	var user model.User_Basic
	err := utils.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

func FindUserByNameAndPassword(name, password string) (*model.User_Basic, error) {
	var user model.User_Basic
	err := utils.DB.Where("username = ?", name).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	if !utils.DecryptMD5(user.Salt, user.Password, password) {
		return nil, errors.New("密码输入错误")
	}
	return &user, nil
}

// 创建用户
func CreateUser(user_register *request.User_Register) error {
	var user model.User_Basic
	user.Username = user_register.Name
	user.Password = user_register.Password
	user.Salt = user_register.Salt
	result := utils.DB.Create(&user)
	return result.Error
}

// 更新用户信息名字和电话，邮箱
func UpdateUserInfo(updateuserinfo *request.UpdateUserInfo) error {
	var user model.User_Basic
	// 检查用户名是否已经被使用
	r := utils.DB.Where("username = ?", updateuserinfo.Username).First(&user).RowsAffected
	if r > 0 {
		return errors.New("该用户名已被使用")
	}

	// 根据id查找user
	err := utils.DB.Where("id = ?", updateuserinfo.ID).First(&user).Error
	if err != nil {
		return err
	}

	result := utils.DB.Model(&user).Updates(map[string]interface{}{"UserName": updateuserinfo.Username, "Phone": updateuserinfo.Phone, "Email": updateuserinfo.Email, "Icon": updateuserinfo.Icon})
	return result.Error
}

func UpdateUserPasswd(name, password string) error {
	var user model.User_Basic
	if err := utils.DB.Where("username = ?", name).First(&user).Error; err != nil {
		return err
	}
	password = utils.EncryptMD5(password, user.Salt)
	result := utils.DB.Model(&user).Updates(map[string]interface{}{"Password": password})
	return result.Error
}

// 删除用户
// 目前的逻辑是逻辑删除，也就是加上一个删除时间，但数据仍在数据库存在
// 是否是逻辑删除（Soft Delete），取决于 模型 User_Basic 是否包含了 GORM 的软删除字段。
// 即默认的gorm.Model包含DeletedAt字段
// 如果想要实际的物理删除
// 可以使用 Unscoped 方法，例如：utils.DB.Unscoped().Delete(&user)
func DeleteUser(name string) error {
	var user model.User_Basic
	if err := utils.DB.Where("username = ?", name).First(&user).Error; err != nil {
		return err
	}
	result := utils.DB.Delete(&user)
	return result.Error
}
