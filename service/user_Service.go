package service

import (
	"fmt"
	"gin_chat/common/response"
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
	data, err := models.GetUserList()
	if err != nil {
		response.FailWithDetailed(data, err.Error(), c)
	}
	response.Ok(c)
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
		response.FailWithMessage(err.Error(), c)
		return
	}

	if user.Password != confirm_password {
		response.FailWithMessage("两次密码不一致", c)
		return
	}

	// 如果用户存在，则返回错误
	if data, err := models.FindUserByName(user.Username); err != nil {
		if data != nil {
			response.FailWithMessage("用户名已存在", c)
			return
		}
	}

	// 校验
	validate := validator.New()
	if err := validate.Struct(&user); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 创建的时候生成一个随机数，用于加密密码
	salt := fmt.Sprintf("%06d", rand.Intn(10000))
	user.Salt = salt
	user.Password = utils.EncryptMD5(user.Password, user.Salt)

	// 创建失败
	if err := models.CreateUser(&user); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("注册成功", c)

}

// GetUserList
// @Summary 获取用户列表
// @Tag 获取用户列表
// @Success 200 {string} json{"code","data"}
// @Router /user/getUserList [get]

// TODO:登陆前后的角色是不同的，登陆后可以发一个token
func Login(c *gin.Context) {
	var user models.User_Basic
	if err := c.ShouldBind(&user); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := models.FindUserByNameAndPassword(user.Username, user.Password)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 登陆颁发
	token, err := utils.GenerateToken(user.Username, user.Password)
	println("token>>>>>>>>", token)

	response.OkWithDetailed(data, "登陆成功", c)
}

// UpdateUser
// @Summary 更新用户
// @Tag 更新用户
// @Success 200 {string} json{"code","data"}
// @Router /user/updateUser [put]
func UpdateUserInfo(c *gin.Context) {
	var user models.User_Basic
	if err := c.ShouldBind(&user); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	validate := validator.New()
	if err := validate.Var(user.Username, "omitempty,min=2,max=100"); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	} else if err := validate.Var(user.Phone, "omitempty,len=3"); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	} else if err := validate.Var(user.Email, "omitempty,email"); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := models.UpdateUserInfo(user.Username, user.Password, user.Phone, user.Email); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("更新成功", c)

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
		response.FailWithMessage(err.Error(), c)
		return
	}

	validate := validator.New()

	if err := validate.Var(newpassword, "omitempty,min=3,max=100"); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	//
	u, err := models.FindUserByName(user.Username)
	if err != nil || u != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if !utils.DecryptMD5(user.Salt, newpassword, u.Password) {
		response.FailWithMessage("密码输入错误", c)
		return
	}

	if err := models.UpdateUserPasswd(user.Username, newpassword); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("更新成功", c)

}

// DeleteUser
// @Summary 删除用户
// @Tag 删除用户
// @Success 200 {string} json{"code","data"}
// @Router //user/deleteUser [delete]
func DeleteUser(c *gin.Context) {
	var user models.User_Basic
	if err := c.ShouldBind(&user); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	validate := validator.New()
	if err := validate.Var(user.Username, "omitempty,min=2,max=100"); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := models.DeleteUser(user.Username); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("删除成功", c)

}

func FindFrend(c *gin.Context) {
	var contact *models.Contact
	ownerid, err := strconv.Atoi(c.Query("ownerid"))
	frendid, err := strconv.Atoi(c.Query("frendid"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if contact, err = models.GetFrend(ownerid, frendid); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(contact, "查找成功", c)
}

func FindFrends(c *gin.Context) {
	users := make([]*models.User_Basic, 0)
	userid, err := strconv.Atoi(c.Query("userid"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if users, err = models.GetFrends(userid); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(users, "查找成功", c)

}

func AddFrend(c *gin.Context) {
	owner, err := strconv.Atoi(c.Query("ownerid"))
	frend, err := strconv.Atoi(c.Query("frendid"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = models.AddFrend(owner, frend)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("添加成功", c)
}

func CreateGroup(c *gin.Context) {
	owner, err := strconv.Atoi(c.Query("ownerid"))
	groupName := c.Query("groupname")
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	group, err := models.CreateGroup(owner, groupName)
	if err != nil {
		response.FailWithDetailed(group, err.Error(), c)
		return
	}

	response.OkWithDetailed(group, "创建群组成功", c)
}

func FindGroup(c *gin.Context) {
	groupName := c.Query("groupname")
	group, err := models.FindGroupByName(groupName)
	if err != nil {
		response.FailWithDetailed(group, err.Error(), c)
		return
	}
	response.OkWithDetailed(group, "查找群组成功", c)
}

func AddGroup(c *gin.Context) {

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
