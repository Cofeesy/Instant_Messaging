package service

import (
	"encoding/json"
	"fmt"
	"gin_chat/common/response"
	"gin_chat/models"
	"gin_chat/models/system"
	"gin_chat/utils"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
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
	var user_register system.User_Register
	if err := c.ShouldBindJSON(&user_register); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if user_register.Password != user_register.Identity {
		response.FailWithMessage("两次密码不一致", c)
		return
	}

	// 如果用户存在，则返回错误
	if data, err := models.FindUserByName(user_register.Name); err != nil {
		if data != nil {
			response.FailWithMessage("用户名已存在", c)
			return
		}
	}
	// 校验
	validate := validator.New()
	if err := validate.Struct(&user_register); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 创建的时候生成一个随机数，用于加密密码
	salt := fmt.Sprintf("%06d", rand.Intn(10000))
	// user.Salt = salt
	user_register.Salt = salt
	user_register.Password = utils.EncryptMD5(user_register.Password, salt)

	// 创建失败
	if err := models.CreateUser(&user_register); err != nil {
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
	var user_login system.User_Login
	// var user models.User_Basic
	if err := c.ShouldBindJSON(&user_login); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	user, err := models.FindUserByNameAndPassword(user_login.Name, user_login.Password)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 登陆颁发,那这个东西是存到数据库吗？
	token, err := utils.GenerateToken(user.Username, user.Password)
	println("token>>>>>>>>", token)

	response.OkWithDetailed(user, "登陆成功", c)
}

// UpdateUser
// @Summary 更新用户
// @Tag 更新用户
// @Success 200 {string} json{"code","data"}
// @Router /user/updateUser [put]
func UpdateUserInfo(c *gin.Context) {
	var user system.UpdateUserInfo
	if err := c.ShouldBindJSON(&user); err != nil {
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

	if err := models.UpdateUserInfo(user.Username, user.Phone, user.Email); err != nil {
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
// func UpdateUserPasswd(c *gin.Context) {
// 	var user system.UpdateUserPasswd
// 	if err := c.ShouldBindJSON(&user); err != nil {
// 		response.FailWithMessage(err.Error(), c)
// 		return
// 	}

// 	validate := validator.New()

// 	if err := validate.Var(user.NewPassword, "omitempty,min=3,max=100"); err != nil {
// 		response.FailWithMessage(err.Error(), c)
// 		return
// 	}

// 	//
// 	u, err := models.FindUserByName(user.Username)
// 	if err != nil || u != nil {
// 		response.FailWithMessage(err.Error(), c)
// 		return
// 	}

// 	if !utils.DecryptMD5(user.Salt, newpassword, u.Password) {
// 		response.FailWithMessage("密码输入错误", c)
// 		return
// 	}

// 	if err := models.UpdateUserPasswd(user.Username, newpassword); err != nil {
// 		response.FailWithMessage(err.Error(), c)
// 		return
// 	}

// 	response.OkWithMessage("更新成功", c)

// }

// DeleteUser
// @Summary 注销
// @Tag 删除用户
// @Success 200 {string} json{"code","data"}
// @Router //user/deleteUser [delete]
func DeleteUser(c *gin.Context) {
	var user models.User_Basic
	if err := c.ShouldBindJSON(&user); err != nil {
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

// 查找某个用户
func Finduser(c *gin.Context) {
	var user system.FindUser
	if err := c.ShouldBindJSON(&user); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	founduser, err := models.FindUserByID(user.UserId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(founduser, "查找成功", c)
}

// 查找用户所有好友
func LoadFriends(c *gin.Context) {
	var friendpayload system.LoadFriendsPayload
	// userid, err := strconv.Atoi(c.Query("userid"))
	err := c.ShouldBindJSON(&friendpayload)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	users, err := models.FindFriendsByUserID(friendpayload.UserId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 封装rows
	data := gin.H{
		"Rows": users,
	}
	response.OkWithDetailed(data, "返回所有好友成功", c)
}

// 用户添加好友
func AddFriend(c *gin.Context) {
	var addfriend system.AddFriend

	err := c.ShouldBindJSON(&addfriend)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = models.AddFrend(&addfriend)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("添加成功", c)
}

// 用户创建群组
func CreateGroup(c *gin.Context) {
	// owner, err := strconv.Atoi(c.Query("ownerid"))
	// groupName := c.Query("groupname")
	var creategroup system.CreatGroup
	err := c.ShouldBindJSON(&creategroup)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	group, err := models.CreateGroup(creategroup)
	if err != nil {
		response.FailWithDetailed(group, err.Error(), c)
		return
	}

	response.OkWithDetailed(group, "创建群组成功", c)
}

// 查找单个group
// func FindGroup(c *gin.Context) {
// 	groupName := c.Query("groupname")
// 	group, err := models.FindGroupByName(groupName)
// 	if err != nil {
// 		response.FailWithDetailed(group, err.Error(), c)
// 		return
// 	}
// 	response.OkWithDetailed(group, "查找群组成功", c)
// }

// 返回群列表
func LoadGroups(c *gin.Context) {
	var loadgroups system.LoadGroups
	err := c.ShouldBindJSON(&loadgroups)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	groups, err := models.FindGroupsByUserID(loadgroups.UserId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 封装Rows
	data := gin.H{
		"Rows": groups,
	}
	response.OkWithDetailed(data, "查找群组成功", c)
}

// 加入群组
func AddGroup(c *gin.Context) {
	var addGroup system.AddGroup
	err := c.ShouldBindJSON(&addGroup)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err = models.AddGroup(&addGroup); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("加入成功", c)
}

// websocket需要序列化反序列化数据，而不是绑定
func WsHandler(c *gin.Context) {
	var authPayload system.AuthMessage

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  512,
		WriteBufferSize: 512,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	// 升级为webdsocket连接
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	// 从连接中读取数据
	_, data, err := ws.ReadMessage()
	if err != nil {
		return
	}

	// 反序列化
	err = json.Unmarshal(data, &authPayload)
	if err != nil {
		return
	}

	// 验证token
	// if authPayload.cmd != 1 {

	// }

	// 不从query里面获取参数
	// err:=c.ShouldBindJSON(&authPayload)
	// if err != nil {
	// 	response.FailWithMessage(err.Error(), c)
	// 	return
	// }

	// 验证token
	// 自定义一个命令号, e.g., 1代表认证
	// if cmd == 1{

	// }

	models.Myws(ws, authPayload.UserId)
}
