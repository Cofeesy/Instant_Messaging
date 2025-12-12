package common

import (
	"github.com/gin-gonic/gin"
	"ZustChat/utils"
	"ZustChat/model/response"
	"time"
)

// 目前先这样，后续根据需求改
func JWT()gin.HandlerFunc{
	return func(c *gin.Context){
		token := c.GetHeader("x-token")

		claims, err := utils.ParseToken(token)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			c.Abort()
			return
		}

		if claims.ExpiresAt-time.Now().Unix() <0{
			response.FailWithMessage("token已过期", c)
			c.Abort()
			return
		}

		c.Next()
	}
}