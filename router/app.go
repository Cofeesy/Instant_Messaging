package router

import (
	"gin_chat/service"
	"gin_chat/utils/setting"

	"github.com/gin-gonic/gin"

	// docs "github.com/go-project-name/docs"
	"gin_chat/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	gin.SetMode(setting.RunMode)

	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// 首页
	r.GET("/index", service.Index)
	r.POST("/register", service.Register)
	r.POST("/login", service.Login)

	// 用户
	r.GET("/user/getUserList", service.GetUserList)
	// r.POST("/user/createUser", service.CreateUser)
	r.PUT("/user/updateUserPasswd", service.UpdateUserPasswd)
	r.PUT("/user/updateUserInfo", service.UpdateUserInfo)
	r.DELETE("/user/deleteUser", service.DeleteUser)

	// chat
	r.GET("/ws", service.WsHandler)

	return r
}
