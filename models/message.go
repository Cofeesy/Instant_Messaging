package models

import (
	"context"
	"encoding/json"
	"fmt"
	// "gin_chat/models"
	"gin_chat/models/system"
	"gin_chat/utils"
	"gin_chat/utils/setting"

	// "net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"google.golang.org/genai"
	"gorm.io/gorm"
)

// 消息源码
// type Message struct {
// 	FromId      int `json:"userid"`
// 	TargetId    int `json:"targetid"`
// 	Content     string `json:"content"`
// 	MessageType int `json:"media"` //发送类型 1私聊  2群聊  3心跳
// 	ContentType int `json:"type"` //消息类型  1文字 2表情包 3语音 4图片 /表情包
// 	Media      int
// }

type Message struct {
	gorm.Model
	UserId     uint   `json:"userid"`
	TargetId   uint   `json:"targetid"`
	Type       int    `json:"type"`  // 1私聊 2群聊 3心跳 4gemini聊天
	Media      int    `json:"media"` // 1文字 2暂时没用 3语音 4图片，表情包
	Content    string `json:"content"`
	CreateTime uint64 `json:"createTime"`
	ReadTime   uint64 `json:"readTime"`
	Pic        string `json:"pic"`
	Url        string `json:"url"`
	Desc       string `json:"desc"`
	Amount     int    `json:"amount"`
}

func (message *Message) TableName() string {
	return "message"
}

// 升级器:将Http协议升级为Websocket协议
// 为什么？
// 定义协议升级的具体细节

var mu sync.RWMutex

// 客户端
type Client struct {
	// 记录用户
	User_id uint
	// *websocket.Conn 类型的对象。这个对象是与单个客户端进行所有通信的唯一凭证和工具。之后的所有操作都是调用这个 conn 对象的方法。
	Conn *websocket.Conn

	// 接受客户端的消息
	// Msg           Message
	HeartbeatTime uint64 //心跳时间
	// 客户端邮箱,存储待发送消息
	SendDataQueue chan []byte
}

// 客户端和用户的映射？
// userid to client
var UserToClient map[uint]*Client = make(map[uint]*Client, 0)

// 监控发送和接受
func Myws(ws *websocket.Conn, userid uint) {

	// 创建客户端实例
	current_time := uint64(time.Now().Unix())
	client := &Client{
		User_id: userid,
		Conn:    ws,
		// Msg:  msg,
		HeartbeatTime: current_time,
		SendDataQueue: make(chan []byte, 256), // 创建带缓冲区的信箱
	}

	// 映射
	// map加读写锁控制并发安全
	mu.Lock()
	UserToClient[userid] = client
	mu.Unlock()

	go client.Send()
	go client.Recieve()

	// 这里应该执行调度

	// 刚启动测试给用户发消息
	// str := "欢迎进入聊天系统"
	// strbyte, _ := json.Marshal(str)
	// SendMsgToUser(userid, strbyte)

}

// 解析收到的消息

// 编码待发送消息

// 每个客户端一直监听即将发送的信息
func (client *Client) Send() {
	for msg := range client.SendDataQueue {
		// 序列化
		// 不需要加指针
		// msg, err := json.Marshal(v)
		// if err != nil {
		// 	fmt.Println("write1:", err.Error())
		// }
		// 传输信息
		err := client.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			fmt.Println("write2:", err.Error())
		}
	}
}

// 监听听并读取接受到的信息，需要一个协程持续监听，这是每一个连接进来的客户端需要做的
func (client *Client) Recieve() {
	// 接受监听消息
	// var Msg Message;
	for {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			fmt.Println("read1:", err.Error())
		}
		// 将数据序列化为代码形式
		// 需要加指针
		// err = json.Unmarshal(msg, &Msg)
		// if err != nil {
		// 	fmt.Println("read2:", err.Error())
		// }
		// fmt.Println("Received message: ", Msg)

		// 接收到消息调度
		client.dispatchMsg(msg)
	}
}

func (client *Client) dispatchMsg(msg []byte) {
	var Msg Message
	json.Unmarshal(msg, &Msg)

	// 确保消息的发送者ID等于当前连接的客户端ID
	// 防止消息中的userid不匹配
	Msg.UserId = client.User_id
	msg, err := json.Marshal(Msg)
	if err != nil {
		fmt.Println("【ERROR】消息序列化失败:", err)
		return
	}

	// fmt.Println("【DEBUG】dispatchMsg: 收到消息，Type=", Msg.Type, ", 发送者ID=", Msg.UserId, ", 目标ID=", Msg.TargetId)

	// TODO:target设置-1的时候会报错，但不影响通信
	switch Msg.Type {
	// 私聊
	case 1:
		// fmt.Println("【DEBUG】处理私聊消息，从", client.User_id, "到", Msg.TargetId)
		SendMsgToUser(client.User_id, Msg.TargetId, msg)
		return
	// 群聊
	case 2:
		// 这里的id是groupid
		// 查找群成员id
		var contacts []Contact = make([]Contact, 0)
		var IDs []uint = make([]uint, 0)
		db.Where("target_id=? AND relation=?", Msg.TargetId, 2).Find(&contacts)

		for _, v := range contacts {
			IDs = append(IDs, v.OwnerId)
		}
		// 不移除发送用户，发送给所有群成员（包括发送者自己）
		// 前端根据 userid 判断是否是自己发的消息，自己发的消息不再渲染
		// remove_ids := client.RemoveId(IDs)

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

	// gemini聊天
	case 4:
		ChatWithGemini(msg)
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
	// send_client := UserToClient[formid]
	recieve_client, ok := UserToClient[targetId]
	mu.RUnlock()

	// TODO:将消息存储到redis上
	// 选择有序表存储
	ctx := context.Background()

	// 私聊消息的 Redis key 格式：包含 type 标识
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

	// 如果初次存储，score设为0
	if score == 0 {
		utils.RDB.ZAdd(ctx, key, redis.Z{Score: 0, Member: msg})
	} else {
		utils.RDB.ZAdd(ctx, key, redis.Z{Score: float64(score), Member: msg})
	}


	if ok && recieve_client != nil {
		// 这里的 select 是为了防止写入阻塞导致协程泄露（可选，但推荐）
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

	// 保存群聊消息到 Redis，key 格式清晰标注为群聊
	// 群聊: chat:group:{groupId}
	ctx := context.Background()
	key := fmt.Sprintf("chat:group:%d", groupId)

	score, err := utils.RDB.ZCard(ctx, key).Result()
	if err != nil {
		fmt.Println("获取群聊消息数量失败:", err)
	}

	// 使用 Redis 有序集合存储群聊消息
	if score == 0 {
		utils.RDB.ZAdd(ctx, key, redis.Z{Score: 0, Member: msg})
	} else {
		utils.RDB.ZAdd(ctx, key, redis.Z{Score: float64(score), Member: msg})
	}

	// 【新增】设置消息过期时间（7天）
	// utils.RDB.Expire(ctx, key, 7*24*time.Hour)

	for _, v := range ids {
		mu.RLock()
		client := UserToClient[v]
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

// 系统广播，比如维护信息
// func SendMsgToAll(){

// }

// 读取start-end的redis数据，并返回给前端
func GetSingleHistoryMsg(redisPayload system.SingleRedisPayload) ([]*Message, error) {
	// 从redis里面查找

	ctx := context.Background()

	// 【修改】组装私聊消息的 key
	key := " "
	if redisPayload.UserId < redisPayload.TargetId {
		key = fmt.Sprintf("chat:private:%d:%d", redisPayload.UserId, redisPayload.TargetId)
	} else {
		key = fmt.Sprintf("chat:private:%d:%d", redisPayload.TargetId, redisPayload.UserId)
	}
	// 前端传来的 Start/End 是相对于最新消息的偏移（Start=0, End=9 表示最近 10 条）
	// 需要把它们转换为 Redis 有序集合的正序索引：
	// 假设总数为 total，则想要的区间为 [total-1-End, total-1-Start]
	// 如果 End == -1，则返回全部（ZRange 0 -1）
	var stringmsgs []string
	// 先获取当前总数
	total, err := utils.RDB.ZCard(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	if redisPayload.End == -1 {
		// 返回所有消息（正序）
		stringmsgs, err = utils.RDB.ZRange(ctx, key, 0, -1).Result()
		if err != nil {
			return nil, err
		}
	} else {
		// 计算正序的 start/end
		if total == 0 {
			return []*Message{}, nil
		}
		// 计算索引
		s := int64(total) - 1 - redisPayload.End
		e := int64(total) - 1 - redisPayload.Start
		if s < 0 {
			s = 0
		}
		if e < 0 {
			// 没有可返回的消息
			return []*Message{}, nil
		}
		stringmsgs, err = utils.RDB.ZRange(ctx, key, s, e).Result()
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}

	msgs := make([]*Message, 0)
	for _, v := range stringmsgs {
		// 每一个string转为[]byte
		stringmsg := []byte(v)
		// 然后将其映射到Message上
		var msg Message
		json.Unmarshal(stringmsg, &msg)
		msgs = append(msgs, &msg)
	}

	return msgs, nil
}

// 【新增】从 Redis 读取群聊消息历史
func GetGroupHistoryMessages(groupRedis *system.GroupRedisPayload) ([]*Message, error) {
	ctx := context.Background()

	// 【修改】群组消息的 Redis key，使用统一的命名规范
	key := fmt.Sprintf("chat:group:%d", groupRedis.GroupId)

	// 从 Redis 有序集合中读取消息
	// ZRange 是从低分数到高分数（正序）
	stringmsgs, err := utils.RDB.ZRange(ctx, key, groupRedis.Start, groupRedis.End).Result()
	if err != nil && err != redis.Nil {
		fmt.Println("从Redis读取群聊消息失败:", err)
		return nil, err
	}

	msgs := make([]*Message, 0)
	for _, v := range stringmsgs {
		var msg Message
		err := json.Unmarshal([]byte(v), &msg)
		if err != nil {
			fmt.Println("解析消息失败:", err)
			continue
		}
		msgs = append(msgs, &msg)
	}

	return msgs, nil
}

// 更新用户心跳
func (client *Client) Heartbeat(currentTime uint64) {
	client.HeartbeatTime = currentTime
}

// 清理超时连接
func CleanConnection(param interface{}) (result bool) {
	result = true
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("cleanConnection err", r)
		}
	}()
	//fmt.Println("定时任务,清理超时连接 ", param)
	//node.IsHeartbeatTimeOut()
	currentTime := uint64(time.Now().Unix())
	for i := range UserToClient {
		client := UserToClient[i]
		if client.IsHeartbeatTimeOut(currentTime) {
			fmt.Println("心跳超时..... 关闭连接：", client)
			client.Conn.Close()
		}
	}
	return result
}

// 用户心跳是否超时
func (client *Client) IsHeartbeatTimeOut(currentTime uint64) (timeout bool) {
	// 每隔多少秒心跳时间
	// HeartbeatHz = 30
	// 最大心跳时间,超过此就下线
	// Time都是用的Unix(),那这个30000是秒
	// HeartbeatMaxTime = 30000
	if client.HeartbeatTime+uint64(setting.HeartbeatMaxTime) <= currentTime {
		fmt.Println("心跳超时...自动下线", client)
		timeout = true
	}
	return
}

// TODO:和gemini聊天后端实现
// 目前仅支持text
// 前端 ——WebSocket—— 后端（Go+Gin+WebSocket） —— 调 AI API（流式） —— 前端展示
func ChatWithGemini(msg []byte) {
	// 客户端创建
	ctx := context.Background()
	cc := genai.ClientConfig{
		APIKey: "AIzaSyCcb4xiwuUgfu56Ie8y8T9oV5Y6-VIIf50",
	}
	aiClient, err := genai.NewClient(ctx, &cc)
	if err != nil {
		fmt.Println(err)
	}

	var Msg Message
	json.Unmarshal(msg, &Msg)
	stream := aiClient.Models.GenerateContentStream(
		ctx,
		"gemini-2.5-flash",
		genai.Text(Msg.Content),
		nil,
	)

	client:=UserToClient[Msg.UserId]

	// 将流消息发送给自己
	for chunk, err := range stream {
		if err==nil{
			part := chunk.Candidates[0].Content.Parts[0]
			fmt.Print(part.Text)
			// 发送给谁呢,发送给自己吧，由前端设置显示
			Text,err:=json.Marshal(part.Text)
			if err==nil{
				client.SendDataQueue <- Text
			}else{
				fmt.Println(err)
			}
		}else{
			fmt.Println(err)
		}
		
	}
}
