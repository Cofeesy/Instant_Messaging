package router

import (
	// "gin_chat/common"
	"gin_chat/utils"
	"gin_chat/api/v1"

	"gin_chat/docs"

	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	gin.SetMode(utils.RunMode)

	// 配置静态资源
	r.Static("/asset", "./asset")
	r.StaticFile("/favicon.ico", "asset/images/favicon.ico")
	r.LoadHTMLGlob("views/**/*")

	// swagger
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// 界面渲染
	r.GET("/index", v1.GetIndex)
	r.GET("/toLogin", v1.ToLogin)
	r.GET("/toRegister", v1.ToRegister)
	r.GET("/toChat", v1.ToChat)

	// 登陆注册
	r.POST("/login", v1.Login)
	r.POST("/register", v1.Register)

	// 处理具体wbsocket逻辑
	r.GET("/chat", v1.WsHandler)

	auth:=r.Group("/")
	// auth.Use(common.JWT())
	{
		auth.GET("/user/getUserList", v1.GetUserList)
		auth.POST("user/findUser", v1.Finduser)
		// r.POST("/user/createUser", service.CreateUser)
		// r.PUT("/user/updateUserPasswd", service.UpdateUserPasswd)
		auth.POST("/user/updateUserInfo", v1.UpdateUserInfo)
		auth.DELETE("/user/deleteUser", v1.DeleteUser)

		// contact
		auth.GET("/findFriend", v1.FindFriend)
		auth.POST("/findFriends", v1.LoadFriends)
		auth.POST("/addFriend", v1.AddFriend)

		// group
		auth.POST("/group/joinGroup", v1.AddGroup)
		auth.POST("/group/loadGroups", v1.LoadGroups)
		auth.POST("/group/createGroup", v1.CreateGroup)

		auth.POST("/attach/upload", v1.UploadInfo)

		// redis历史消息
		auth.POST("/user/getSingleMessagesFromRedis", v1.GetSingleMessagesFromRedis)
		auth.POST("/message/getGroupMessagesFromRedis", v1.GetGroupMessagesFromRedis)
		auth.POST("/message/getAiMessagesFromRedis", v1.GetAiMessagesFromRedis)
	}

	return r
}
