package service

import (
	"gin_chat/common/response"
	"gin_chat/models"
	"html/template"
	"github.com/gin-gonic/gin"
	"strconv"
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
	t.Execute(c.Writer, "index is ok")
}

func ToRegister(c *gin.Context) {
	t, err := template.ParseFiles("views/user/register.html")
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	}
	t.Execute(c.Writer, "register")
}

func ToLogin(c *gin.Context) {
	t, err := template.ParseFiles("views/user/login.html")
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	}
	t.Execute(c.Writer, "login")
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
	userId, _ := strconv.Atoi(c.Query("userId"))
	token := c.Query("token")
	user := models.User_Basic{}
	user.ID = uint(userId)
	user.LoginToken = token	
	//fmt.Println("ToChat>>>>>>>>", user)
	t.Execute(c.Writer, user)
}
