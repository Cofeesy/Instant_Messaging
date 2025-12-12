package service

import (
	"context"
	"encoding/json"
	"fmt"
	"ZustChat/global"
	"ZustChat/model"
	"ZustChat/utils"
	"time"

	"github.com/redis/go-redis/v9"
	"google.golang.org/genai"

	"go.uber.org/zap"
)

func dispatchMsg(msg []byte, client *Client) {
	var Msg model.Message
	json.Unmarshal(msg, &Msg)

	// 确保消息的发送者ID等于当前连接的客户端ID
	// 防止消息中的userid不匹配
	Msg.UserId = client.User_id
	msg, err := json.Marshal(Msg)
	if err != nil {
		global.Logger.Error("消息序列化失败",
			zap.Uint("user_id", client.User_id),
			zap.Error(err),
		)
		return
	}

	// TODO:target设置-1的时候会报错，但不影响通信
	switch Msg.Type {
	// 私聊
	case 1:
		if Msg.TargetId == 0 {
			ChatWithGemini(msg)
			return
		} else {
			SendMsgToUser(client.User_id, Msg.TargetId, msg)
			return
		}
	// 群聊
	case 2:
		var contacts []model.Contact = make([]model.Contact, 0)
		var IDs []uint = make([]uint, 0)
		utils.DB.Where("target_id=? AND relation=?", Msg.TargetId, 2).Find(&contacts)

		for _, v := range contacts {
			IDs = append(IDs, v.OwnerId)
		}

		// 如果发送者不在 IDs 中（群成员列表），需要添加发送者
		// 这样确保发送者也能收到自己的消息（用于群聊消息历史同步）
		hasSender := false
		for _, id := range IDs {
			if id == client.User_id {
				hasSender = true
				break
			}
		}
		if !hasSender {
			IDs = append(IDs, client.User_id)
		}

		SendMsgToGroup(IDs, msg, Msg.TargetId)
		return
	// 心跳
	case 3:
		client.HeartbeatTime = uint64(time.Now().Unix())
		return
	}
}

// 私聊
// 写消息进去
// 主要是写是建立在连接双方的写，通过把消息放到其他人的send上，服务器读取到会将该消息发送到该send连接到的client上
// 这里的“写” (conn.WriteMessage)，最终都是发生在某一个具体的 WebSocket 连接 (conn) 上的。
// 即clientA 的 Send goroutine 只会在 clientA.conn 上写，clientB 的 Send goroutine 只会在 clientB.conn 上写。职责非常明确
// 关键点是看谁开启了这个协程，服务器只是负责发送和接受信息
func SendMsgToUser(formid, targetId uint, msg []byte) {
	var mysqlMsg model.Message
	err := json.Unmarshal(msg, &mysqlMsg)
	if err != nil {
		global.Logger.Error("解析消息失败",
			zap.Uint("from_id", formid),
			zap.Uint("target_id", targetId),
			zap.Error(err),
		)
		return
	}

	// 使用服务器的时间
	nowMilli := uint64(time.Now().UnixMilli())
	mysqlMsg.CreateTime = nowMilli
	// 存mysql
	if err := utils.DB.Create(&mysqlMsg).Error; err != nil {
		global.Logger.Error("私聊消息写入MySQL失败",
			zap.Uint("from_id", formid),
			zap.Uint("target_id", targetId),
			zap.Error(err),
		)
		return
	}

	// 更新时间后的msg，redis数据需和mysql保持一致
	newMsgBytes, err := json.Marshal(mysqlMsg)
	if err != nil {
		global.Logger.Error("消息序列化失败",
			zap.Uint("from_id", formid),
			zap.Uint("target_id", targetId),
			zap.Error(err),
		)
		return
	}

	// 存redis
	ctx := context.Background()
	// 私聊: chat:private:{小ID}:{大ID}
	key := ""
	if formid < targetId {
		key = fmt.Sprintf("chat:private:%d:%d", formid, targetId)
	} else {
		key = fmt.Sprintf("chat:private:%d:%d", targetId, formid)
	}

	// 使用管道优化：一次发送多条命令给 Redis，让 Redis 连续执行后再一次性返回结果。
	// 不是原子操作，只是减少网络rtt
	pipe := utils.RDB.Pipeline()
	// 时间戳作为score
	score := float64(mysqlMsg.CreateTime)
	pipe.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: newMsgBytes,
	})

	// 设置过期时间
	pipe.Expire(ctx, key, 7*24*time.Hour)

	// ZRemRangeByRank: 裁剪/清理旧消息 (实现滑动窗口)
	// 只保留最新的 500 条。
	// 逻辑：移除排名为 0 到 -501 的元素 (即保留倒数 500 个)
	pipe.ZRemRangeByRank(ctx, key, 0, -501)

	_, err = pipe.Exec(ctx)
	if err != nil {
		global.Logger.Error("私聊消息写入Redis失败",
			zap.Uint("from_id", formid),
			zap.Uint("target_id", targetId),
			zap.String("key", key),
			zap.Error(err),
		)
		// 注意：Redis 失败通常不回滚 MySQL，记录日志即可
	}

	// TODO:以上内容是否可以异步处理？-->消息队列
	// 推送消息
	mu.RLock()
	recieve_client, ok := GlobalHub.UserToClient[targetId]
	mu.RUnlock()
	if ok && recieve_client != nil {
		// 防止写入阻塞导致协程泄露
		select {
		case recieve_client.SendDataQueue <- newMsgBytes:
			global.Logger.Debug("消息已实时推送给用户",
				zap.Uint("target_id", targetId),
			)
		default:
			global.Logger.Warn("用户发送队列已满，消息仅存储在Redis",
				zap.Uint("target_id", targetId),
			)
		}
	} else {
		global.Logger.Debug("用户不在线或连接异常，消息已保存到Redis",
			zap.Uint("target_id", targetId),
		)
	}

}

// 群聊
func SendMsgToGroup(ids []uint, msg []byte, groupId uint) {
	var mysqlMsg model.Message
	err := json.Unmarshal(msg, &mysqlMsg)
	if err != nil {
		global.Logger.Error("解析群聊消息失败",
			zap.Uint("group_id", groupId),
			zap.Error(err),
		)
		return
	}

	// 使用服务器的时间
	nowMilli := uint64(time.Now().UnixMilli())
	mysqlMsg.CreateTime = nowMilli
	// 存mysql
	if err := utils.DB.Create(&mysqlMsg).Error; err != nil {
		global.Logger.Error("群聊消息写入MySQL失败",
			zap.Uint("group_id", groupId),
			zap.Error(err),
		)
		return
	}

	// 更新时间后的msg，redis数据需和mysql保持一致
	newMsgBytes, err := json.Marshal(mysqlMsg)
	if err != nil {
		global.Logger.Error("群聊消息序列化失败",
			zap.Uint("group_id", groupId),
			zap.Error(err),
		)
		return
	}

	// 存redis
	ctx := context.Background()
	// 私聊: chat:private:{小ID}:{大ID}
	key := fmt.Sprintf("chat:group:%d", groupId)

	// 使用管道优化：一次发送多条命令给 Redis，让 Redis 连续执行后再一次性返回结果。
	// 不是原子操作，只是减少网络rtt
	pipe := utils.RDB.Pipeline()
	// 时间戳作为score
	score := float64(mysqlMsg.CreateTime)
	pipe.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: newMsgBytes,
	})

	// 设置过期时间
	pipe.Expire(ctx, key, 7*24*time.Hour)
	pipe.ZRemRangeByRank(ctx, key, 0, -501)

	_, err = pipe.Exec(ctx)
	if err != nil {
		global.Logger.Error("群聊消息写入Redis失败",
			zap.Uint("group_id", groupId),
			zap.String("key", key),
			zap.Error(err),
		)
	}

	for _, v := range ids {
		mu.RLock()
		client := GlobalHub.UserToClient[v]
		mu.RUnlock()
		if client != nil {
			select {
			case client.SendDataQueue <- newMsgBytes:
				global.Logger.Debug("群聊消息已发送给用户",
					zap.Uint("group_id", groupId),
					zap.Uint("user_id", v),
				)
			default:
				global.Logger.Warn("用户发送队列已满，消息仅存储在Redis",
					zap.Uint("group_id", groupId),
					zap.Uint("user_id", v),
				)
			}
		}
	}
}

// gemini聊天后端实现
// 目前仅支持text
// 前端 ——WebSocket—— 后端（Go+Gin+WebSocket） —— 调 AI API（流式） —— 前端展示
func ChatWithGemini(msg []byte) {
	ctx := context.Background()

	// 初始化 AI 客户端
	cc := genai.ClientConfig{
		APIKey: utils.GeminiKey,
	}
	aiClient, err := genai.NewClient(ctx, &cc)
	if err != nil {
		global.Logger.Error("AI客户端初始化失败",
			zap.Error(err),
		)
		return // 客户端都起不来，直接返回
	}

	// 解析用户消息
	var userMsg model.Message
	if err := json.Unmarshal(msg, &userMsg); err != nil {
		global.Logger.Error("解析AI聊天消息失败",
			zap.Error(err),
		)
		return
	}

	// 补全用户消息字段 (服务器权威时间)
	userMsg.CreateTime = uint64(time.Now().UnixMilli())
	userMsg.TargetId = 0 
	userMsg.Type = 1    


	if err := utils.DB.Create(&userMsg).Error; err != nil {
		global.Logger.Error("用户AI消息写入MySQL失败",
			zap.Uint("user_id", userMsg.UserId),
			zap.Error(err),
		)
		return
	}

	key := fmt.Sprintf("aichat:%d:%d", 0, userMsg.UserId)

	userMsgBytes, _ := json.Marshal(userMsg)

	// 立即写入 Redis，不要等 AI 回复完
	pipe := utils.RDB.Pipeline()
	pipe.ZAdd(ctx, key, redis.Z{
		Score:  float64(userMsg.CreateTime),
		Member: userMsgBytes,
	})
	pipe.Expire(ctx, key, 7*24*time.Hour)
	pipe.ZRemRangeByRank(ctx, key, 0, -501)
	pipe.Exec(ctx)

	stream := aiClient.Models.GenerateContentStream(
		ctx,
		"gemini-2.5-flash",
		genai.Text(userMsg.Content),
		nil,
	)

	// 获取 WebSocket 客户端用于推送
	client, isOnline := GlobalHub.UserToClient[userMsg.UserId]

	fullAIResponse := "" // 收集完整回复

	for chunk, err := range stream {
		if err == nil {
			if len(chunk.Candidates) > 0 && len(chunk.Candidates[0].Content.Parts) > 0 {
				part := chunk.Candidates[0].Content.Parts[0]
				// 将 Part 转换为字符串
				textStr := fmt.Sprintf("%v", part)

				fullAIResponse += textStr

				// 实时推送给前端
				if isOnline && client != nil {
					// 直接推字符串给前端，前端收到后追加显示
					respBytes, _ := json.Marshal(textStr)
					select {
					case client.SendDataQueue <- respBytes:
					default:
					}
				}
			}
		} else {
			global.Logger.Error("AI流式响应错误",
				zap.Uint("user_id", userMsg.UserId),
				zap.Error(err),
			)
			break
		}
	}

	// 构建 AI 消息对象
	aiMsg := model.Message{
		UserId:     0,              // AI 发送
		TargetId:   userMsg.UserId, // 用户接收
		Content:    fullAIResponse,
		Type:       1,
		Media:      1,
		CreateTime: uint64(time.Now().UnixMilli()), // 生成结束的时间
	}

	// 存 MySQL
	if err := utils.DB.Create(&aiMsg).Error; err != nil {
		global.Logger.Error("AI回复消息写入MySQL失败",
			zap.Uint("user_id", userMsg.UserId),
			zap.Error(err),
		)
	}

	// 存 Redis
	aiMsgBytes, _ := json.Marshal(aiMsg)

	pipe = utils.RDB.Pipeline()
	pipe.ZAdd(ctx, key, redis.Z{
		Score:  float64(aiMsg.CreateTime),
		Member: aiMsgBytes,
	})
	pipe.Expire(ctx, key, 7*24*time.Hour)
	pipe.ZRemRangeByRank(ctx, key, 0, -501)

	if _, err := pipe.Exec(ctx); err != nil {
		global.Logger.Error("AI回复消息写入Redis失败",
			zap.Uint("user_id", userMsg.UserId),
			zap.String("key", key),
			zap.Error(err),
		)
	}
}
