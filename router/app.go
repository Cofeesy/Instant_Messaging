package router

import (
	"gin_chat/service"
	"gin_chat/utils/setting"

	"github.com/gin-gonic/gin"
	"gin_chat/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	gin.SetMode(setting.RunMode)

	// 配置静态资源
	r.Static("/asset", "./asset")
	r.StaticFile("/favicon.ico", "asset/images/favicon.ico")
	//	r.StaticFS()
	r.LoadHTMLGlob("views/**/*")

	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// 首页
	r.GET("/index", service.GetIndex)
	r.GET("/toLogin", service.ToLogin)
	r.GET("/toRegister", service.ToRegister)
	r.GET("/toChat", service.ToChat)

	// user
	r.GET("/user/getUserList", service.GetUserList)
	r.POST("/login", service.Login)
	r.POST("/register", service.Register)
	// r.POST("/user/createUser", service.CreateUser)
	r.PUT("/user/updateUserPasswd", service.UpdateUserPasswd)
	r.PUT("/user/updateUserInfo", service.UpdateUserInfo)
	// 用户注销
	r.DELETE("/user/deleteUser", service.DeleteUser)

	// contact
	r.GET("/user/findFrend",service.FindFrend)
	r.GET("/user/findFrends",service.FindFrends)
	r.POST("/user/addFrend",service.AddFrend)

	// group
	r.GET("/user/findGroup",service.FindGroup)
	r.POST("/user/createGroup",service.CreateGroup)
	
	// chat
	r.GET("/chat", service.WsHandler)

	return r
}
