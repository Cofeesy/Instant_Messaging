package service

import (
	"github.com/gin-gonic/gin"
)

// GetUserList
// @Tag 获取用户列表
// @Success 200 {string} ok
// @Router /index [get]
func Index(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ok",
	})
}
