package router

import (
	"gin_chat/service"
	"gin_chat/utils/setting"

	"gin_chat/docs"

	"github.com/gin-gonic/gin"

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
	// toChat是前端页面响应
	r.GET("/toChat", service.ToChat)

	r.POST("/login", service.Login)
	r.POST("/register", service.Register)
	// chat是处理具体wbsocket逻辑
	r.GET("/chat", service.WsHandler)

	// user
	r.GET("/user/getUserList", service.GetUserList)
	r.POST("user/findUser", service.Finduser)
	// r.POST("/user/createUser", service.CreateUser)
	// r.PUT("/user/updateUserPasswd", service.UpdateUserPasswd)
	r.POST("/user/updateUserInfo", service.UpdateUserInfo)
	// 用户注销
	r.DELETE("/user/deleteUser", service.DeleteUser)

	// contact
	r.GET("/findFriend", service.FindFriend)
	r.POST("/findFriends", service.LoadFriends)
	r.POST("/addFriend", service.AddFriend)

	// group
	// r.GET("/findGroup",service.FindGroup)
	r.POST("/group/joinGroup", service.AddGroup)
	r.POST("/group/loadGroups", service.LoadGroups)
	r.POST("/group/createGroup", service.CreateGroup)
	// 以上都已经Api测试过,ok

	r.POST("/user/getSingleMessagesFromRedis", service.GetSingleMessagesFromRedis)

	r.POST("/attach/upload", service.UploadInfo)

	// 【新增】群聊消息接口
	r.POST("/message/getGroupMessagesFromRedis", service.GetGroupMessagesFromRedis)

	return r
}
