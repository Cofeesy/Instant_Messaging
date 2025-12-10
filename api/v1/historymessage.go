package v1

import (
	"gin_chat/model/response"
	"gin_chat/model/request"
	"gin_chat/service"

	// "github.com/aws/aws-sdk-go/aws/request"
	"github.com/gin-gonic/gin"
)

func GetSingleMessagesFromRedis(c *gin.Context) {
	var redisPayload request.SingleHistoryMsgReq
	err := c.ShouldBindJSON(&redisPayload)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	redismsg, err := service.GetSingleHistoryMsg(redisPayload)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	data := gin.H{
		"Rows": redismsg,
	}
	response.OkWithDetailed(data, "返回私聊redis消息成功", c)
}

// GetGroupMessagesFromRedis
// @Summary 从Redis获取群组消息历史
// @Tag 消息
// @Param groupId query int true "群组ID"
// @Success 200 {string} json{"code","data"}
// @Router /message/getGroupMessagesFromRedis [post]
func GetGroupMessagesFromRedis(c *gin.Context) {
	var groupRedis request.GroupHistoryMsgReq

	if err := c.ShouldBindJSON(&groupRedis); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 如果没有提供分页参数，默认读取所有消息
	// if groupRedis.End == 0 {
	// 	groupRedis.End = -1
	// }

	// 从 Redis 中读取群聊消息
	redismsg, err := service.GetGroupHistoryMessages(&groupRedis)
	if err != nil {
		response.FailWithMessage("读取消息失败: "+err.Error(), c)
		return
	}
	data := gin.H{
		"Rows": redismsg,
	}
	response.OkWithDetailed(data, "返回群聊redis消息成功", c)

}

// ai聊天历史记录
func GetAiMessagesFromRedis(c *gin.Context) {
	var aiRedis request.AiHistoryMsgReq

	if err := c.ShouldBindJSON(&aiRedis); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 从 Redis 中读取群聊消息
	redismsg, err := service.GetAiHistoryMessages(&aiRedis)
	if err != nil {
		response.FailWithMessage("读取消息失败: "+err.Error(), c)
		return
	}
	
	data := gin.H{
		"Rows": redismsg,
	}
	response.OkWithDetailed(data, "返回ai_redis消息成功", c)
}