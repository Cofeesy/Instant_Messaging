package service

import (
	"gin_chat/common/response"
	"html/template"

	"github.com/gin-gonic/gin"
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


func ToRegister(c *gin.Context){
	t,err:=template.ParseFiles("views/user/register.html")
	if err!=nil{
		response.FailWithMessage(err.Error(), c)
	}
	t.Execute(c.Writer,"register")
}

func ToChat(c *gin.Context){

}
