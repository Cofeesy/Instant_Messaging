package v1

import (
	"fmt"
	"gin_chat/model"
	"gin_chat/model/response"
	"gin_chat/model/request"
	"gin_chat/service"
	"gin_chat/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"html/template"
	"math/rand"
)

// Server-Side Rendering, SSR，后端渲染的项目，不是前后端分离的

// GetUserList
// @Tag 获取用户列表
// @Success 200 {string} ok
// @Router /index [get]
func GetIndex(c *gin.Context) {
	// 加载前端文件
	// 被嵌套的模版要在后面解析
	t, err := template.ParseFiles("index.html", "views/chat/head.html")
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	}

	// 将生成的 HTML写入 HTTP 响应体中。
	t.Execute(c.Writer, "index 跳转成功")
}

func ToRegister(c *gin.Context) {
	t, err := template.ParseFiles("views/user/register.html")
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	}
	t.Execute(c.Writer, "register跳转成功")
}

func ToLogin(c *gin.Context) {
	t, err := template.ParseFiles("views/user/login.html")
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	}
	t.Execute(c.Writer, "login跳转成功")
}

func ToChat(c *gin.Context) {
	t, err := template.ParseFiles("views/chat/index.html",
		"views/chat/head.html",
		"views/chat/foot.html",
		"views/chat/tabmenu.html",
		"views/chat/concat.html",
		"views/chat/group.html",
		"views/chat/profile.html",
		"views/chat/createcom.html",
		"views/chat/userinfo.html",
		"views/chat/main.html")

	if err != nil {
		panic(err)
	}
	// 获取信息
	// TODO:改其他方式，比如session的地方获取
	// userId, _ := strconv.Atoi(c.Query("userId"))
	// token := c.Query("token")
	// user := models.User_Basic{}
	// user.ID = uint(userId)
	// user.LoginToken = token
	user := model.User_Basic{}
	// fmt.Println("ToChat>>>>>>>>", user)
	// 返回给前端的数据
	t.Execute(c.Writer, user)
}

func Register(c *gin.Context) {
	var user_register request.User_Register
	if err := c.ShouldBindJSON(&user_register); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if user_register.Password != user_register.Repassword {
		response.FailWithMessage("两次密码不一致", c)
		return
	}

	// 如果用户存在，则返回错误
	if data, err := service.FindUserByName(user_register.Name); err != nil {
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
	if err := service.CreateUser(&user_register); err != nil {
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
func Login(c *gin.Context) {
	var user_login request.User_Login
	if err := c.ShouldBindJSON(&user_login); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	user, err := service.FindUserByNameAndPassword(user_login.Name, user_login.Password)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		response.FailWithMessage("生成Token失败", c)
		return
	}
	println("token>>>>>>>>", token)
	response.OkWithDetailed(gin.H{
		"user":  user,
		"token": token,
	}, "登陆成功", c)
}
