package router

import (
	"gin_chat/common"
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
	r.LoadHTMLGlob("views/**/*")

	// swagger
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// 界面渲染
	r.GET("/index", service.GetIndex)
	r.GET("/toLogin", service.ToLogin)
	r.GET("/toRegister", service.ToRegister)
	r.GET("/toChat", service.ToChat)

	// 登陆注册
	r.POST("/login", service.Login)
	r.POST("/register", service.Register)

	// 处理具体wbsocket逻辑
	r.GET("/chat", service.WsHandler)

	auth:=r.Group("/")
	auth.Use(common.JWT())
	{
		auth.GET("/user/getUserList", service.GetUserList)
		auth.POST("user/findUser", service.Finduser)
		// r.POST("/user/createUser", service.CreateUser)
		// r.PUT("/user/updateUserPasswd", service.UpdateUserPasswd)
		auth.POST("/user/updateUserInfo", service.UpdateUserInfo)
		auth.DELETE("/user/deleteUser", service.DeleteUser)

		// contact
		auth.GET("/findFriend", service.FindFriend)
		auth.POST("/findFriends", service.LoadFriends)
		auth.POST("/addFriend", service.AddFriend)

		// group
		auth.POST("/group/joinGroup", service.AddGroup)
		auth.POST("/group/loadGroups", service.LoadGroups)
		auth.POST("/group/createGroup", service.CreateGroup)

		auth.POST("/attach/upload", service.UploadInfo)

		// redis历史消息
		auth.POST("/user/getSingleMessagesFromRedis", service.GetSingleMessagesFromRedis)
		auth.POST("/message/getGroupMessagesFromRedis", service.GetGroupMessagesFromRedis)
		auth.POST("/message/getAiMessagesFromRedis", service.GetAiMessagesFromRedis)
	}

	return r
}
