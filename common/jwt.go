package common

import (
	"github.com/gin-gonic/gin"
	"gin_chat/utils"
	"time"
)

// 目前先这样，后续根据需求改
func JWT()gin.HandlerFunc{
	return func(c *gin.Context){
		token :=c.Query("token")

		claims, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if claims.ExpiresAt-time.Now().Unix() <0{
			c.JSON(401, gin.H{"error": "token已过期"})
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Set("password", claims.Password)

		c.Next()
	}
}