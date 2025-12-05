package service

import (
	"context"
	"encoding/json"
	"fmt"
	"gin_chat/models"
	"gin_chat/utils"
	"time"
	"github.com/redis/go-redis/v9"
	"google.golang.org/genai"
)

func dispatchMsg(msg []byte,client *Client) {
	var Msg models.Message
	json.Unmarshal(msg, &Msg)

	// 确保消息的发送者ID等于当前连接的客户端ID
	// 防止消息中的userid不匹配
	Msg.UserId = client.User_id
	msg, err := json.Marshal(Msg)
	if err != nil {
		fmt.Println("【ERROR】消息序列化失败:", err)
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
		var contacts []models.Contact = make([]models.Contact, 0)
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
	mu.RLock()
	recieve_client, ok :=GlobalHub.UserToClient[targetId]
	mu.RUnlock()

	ctx := context.Background()

	// 私聊: chat:private:{小ID}:{大ID}
	key := ""
	if formid < targetId {
		key = fmt.Sprintf("chat:private:%d:%d", formid, targetId)
	} else {
		key = fmt.Sprintf("chat:private:%d:%d", targetId, formid)
	}

	// 使用redis的有序集合
	// zcard获取有序集合的成员数
	score, err := utils.RDB.ZCard(ctx, key).Result()
	if err != nil {
		fmt.Println(">>>>>>>>>err:", err)
		return
	}

	// 初次存储
	if score == 0 {
		utils.RDB.ZAdd(ctx, key, redis.Z{Score: 0, Member: msg})
	} else {
		utils.RDB.ZAdd(ctx, key, redis.Z{Score: float64(score), Member: msg})
	}

	if ok && recieve_client != nil {
		// 防止写入阻塞导致协程泄露
		select {
		case recieve_client.SendDataQueue <- msg:
			fmt.Println("【INFO】消息已实时推送给用户", targetId)
		default:
			fmt.Println("【WARNING】用户", targetId, "的发送队列已满，消息仅存储在Redis")
		}
	} else {
		fmt.Println("【WARNING】用户", targetId, "不在线或连接异常，消息已保存到Redis")
	}

}

// 群聊
func SendMsgToGroup(ids []uint, msg []byte, groupId uint) {
	if len(ids) == 0 {
		return
	}

	// 群聊: chat:group:{groupId}
	ctx := context.Background()
	key := fmt.Sprintf("chat:group:%d", groupId)

	score, err := utils.RDB.ZCard(ctx, key).Result()
	if err != nil {
		fmt.Println("获取群聊消息数量失败:", err)
	}

	if score == 0 {
		utils.RDB.ZAdd(ctx, key, redis.Z{Score: 0, Member: msg})
	} else {
		utils.RDB.ZAdd(ctx, key, redis.Z{Score: float64(score), Member: msg})
	}

	// 【新增】设置消息过期时间（7天）
	// utils.RDB.Expire(ctx, key, 7*24*time.Hour)

	for _, v := range ids {
		mu.RLock()
		client := GlobalHub.UserToClient[v]
		mu.RUnlock()
		if client != nil {
			select {
			case client.SendDataQueue <- msg:
				fmt.Println("消息已发送给用户:", v)
			default:
				fmt.Println("用户", v, "的发送队列已满，消息仅存储在Redis")
			}
		}
	}
}

// gemini聊天后端实现
// 目前仅支持text
// 前端 ——WebSocket—— 后端（Go+Gin+WebSocket） —— 调 AI API（流式） —— 前端展示
func ChatWithGemini(msg []byte) {
	ctx := context.Background()
	cc := genai.ClientConfig{
		APIKey: utils.GeminiKey,
	}

	aiClient, err := genai.NewClient(ctx, &cc)
	if err != nil {
		fmt.Println(err)
	}

	var Msg models.Message
	json.Unmarshal(msg, &Msg)
	stream := aiClient.Models.GenerateContentStream(
		ctx,
		"gemini-2.5-flash",
		genai.Text(Msg.Content),
		nil,
	)

	// redis存ai消息
	key := fmt.Sprintf("aichat:%d:%d", 0, Msg.UserId)

	score, err := utils.RDB.ZCard(ctx, key).Result()
	if err != nil {
		fmt.Println("获取群聊消息数量失败:", err)
	}

	if score == 0 {
		utils.RDB.ZAdd(ctx, key, redis.Z{Score: 0, Member: msg})
	} else {
		utils.RDB.ZAdd(ctx, key, redis.Z{Score: float64(score), Member: msg})
	}

	client := GlobalHub.UserToClient[Msg.UserId]

	aiRedisContent := " "
	// 将流消息发送给自己
	for chunk, err := range stream {
		if err == nil {
			part := chunk.Candidates[0].Content.Parts[0]
			// 发送给谁呢,发送给自己吧，由前端设置显示
			aiRedisContent += part.Text
			Text, err := json.Marshal(part.Text)
			if err == nil {
				client.SendDataQueue <- Text
			} else {
				fmt.Println(err)
			}
		} else {
			fmt.Println(err)
		}

	}
	// 这里也得记录，因为是发给自己
	aiRedisMsg := Msg
	aiRedisMsg.UserId, aiRedisMsg.TargetId = aiRedisMsg.TargetId, aiRedisMsg.UserId
	aiRedisMsg.Content = aiRedisContent
	airedismsg, _ := json.Marshal(aiRedisMsg)
	utils.RDB.ZAdd(ctx, key, redis.Z{Score: float64(score + 1), Member: airedismsg})
}