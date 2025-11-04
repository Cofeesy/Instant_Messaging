package service

import (
	"fmt"
	"gin_chat/models"
	"gin_chat/utils"
	"math/rand"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// @Success <code> {<type>} <model or object> "<description>"

// GetUserList
// @Summary 获取用户列表
// @Tag 获取用户列表
// @Success 200 {string} json{"code","data"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := models.GetUserList()
	c.JSON(200, gin.H{
		"message": data,
	})
}

// 测试成功，应该能看到数据库该用户并且该用户有salt值
// Login
// @Summary 创建用户
// @Tag 创建用户
// @Success 200 {string} json{"code","data"}
// @Router /user/createUser [post]
func Register(c *gin.Context) {
	var user models.User_Basic
	confirm_password := c.Query("confirm_password")

	if err := c.ShouldBind(&user); err != nil {
		c.JSON(400, gin.H{
			"error": "111",
		})
		return
	}

	if user.Password != confirm_password {
		c.JSON(400, gin.H{
			"error": "两次输入的密码不一致",
		})
		return
	}

	// 如果用户存在，则返回错误
	if data, err := models.FindUserByName(user.Username); err != nil {
		if data != nil {
			c.JSON(400, gin.H{
				"error": "user already exists",
			})
			return
		}
	}

	// 校验
	validate := validator.New()
	if err := validate.Struct(&user); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 创建的时候生成一个随机数，用于加密密码
	salt := fmt.Sprintf("%06d", rand.Intn(10000))
	user.Salt = salt
	user.Password = utils.EncryptMD5(user.Password, user.Salt)

	// 创建失败
	if err := models.CreateUser(&user); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "user created successfully",
	})

}

// GetUserList
// @Summary 获取用户列表
// @Tag 获取用户列表
// @Success 200 {string} json{"code","data"}
// @Router /user/getUserList [get]
func Login(c *gin.Context) {
	var user models.User_Basic
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	data, err := models.FindUserByNameAndPassword(user.Username, user.Password)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 登陆颁发
	token, err := utils.GenerateToken(user.Username, user.Password)
	println("token>>>>>>>>", token)
	// FIXME:更新用户identity字段,
	//
	c.JSON(200, gin.H{
		"message": data,
	})
}

// UpdateUser
// @Summary 更新用户
// @Tag 更新用户
// @Success 200 {string} json{"code","data"}
// @Router /user/updateUser [put]
func UpdateUserInfo(c *gin.Context) {
	var user models.User_Basic
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	validate := validator.New()
	if err := validate.Var(user.Username, "omitempty,min=2,max=100"); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	} else if err := validate.Var(user.Phone, "omitempty,len=3"); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	} else if err := validate.Var(user.Email, "omitempty,email"); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := models.UpdateUserInfo(user.Username, user.Password, user.Phone, user.Email); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "user updated successfully",
	})

}

// UpdateUser
// @Summary 更新用户
// @Tag 更新用户
// @Param username query string true "用户名"
// @Param password query string true "密码"
// @Success 200 {string} json{"code","data"}
// @Router /user/updateUser [put]
func UpdateUserPasswd(c *gin.Context) {
	var user models.User_Basic
	newpassword := c.Query("newpassword")
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	validate := validator.New()

	if err := validate.Var(newpassword, "omitempty,min=3,max=100"); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	//
	u, err := models.FindUserByName(user.Username)
	if err != nil || u != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !utils.DecryptMD5(user.Salt, newpassword, u.Password) {
		c.JSON(500, gin.H{
			"error": "密码输入错误",
		})
		return
	}

	if err := models.UpdateUserPasswd(user.Username, newpassword); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "user updated successfully",
	})

}

// DeleteUser
// @Summary 删除用户
// @Tag 删除用户
// @Success 200 {string} json{"code","data"}
// @Router //user/deleteUser [delete]
func DeleteUser(c *gin.Context) {
	var user models.User_Basic
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	validate := validator.New()
	if err := validate.Var(user.Username, "omitempty,min=2,max=100"); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := models.DeleteUser(user.Username); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "user deleted successfully",
	})

}

func WsHandler(c *gin.Context) {
	var msg models.Message
	// 这里怎么得到这些信息，用apifox测试好像不行
	FormId := c.Query("FormId")
	TargetId := c.Query("TargetId")
	Content := c.Query("Content")

	msg.FromId, _ = strconv.Atoi(FormId)
	msg.TargetId, _ = strconv.Atoi(TargetId)
	msg.Content = Content
	models.Myws(msg, c)

}
