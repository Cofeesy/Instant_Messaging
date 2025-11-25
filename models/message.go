package models

import (
	"context"
	"encoding/json"
	"fmt"
	"gin_chat/models/system"
	"gin_chat/utils"
	"gin_chat/utils/setting"

	// "net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
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
	UserId     uint   //发送者
	TargetId   uint   //接受者
	Type       int    //发送类型  1私聊  2群聊  3心跳
	Media      int    //消息类型  1文字 2表情包 3图片 4音频
	Content    string //消息内容
	CreateTime uint64 //创建时间
	ReadTime   uint64 //读取时间
	Pic        string
	Url        string
	Desc       string
	Amount     int //其他数字统计
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
	// if err != nil {
	// 	fmt.Println("read2:", err.Error())
	// }
	// TODO:target设置-1的时候会报错，但不影响通信
	switch Msg.Type {
	// 私聊
	case 1:
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
		// // 除去发送用户
		// fmt.Println(">>>>>>>>>>>2:",IDs)
		remove_ids := client.RemoveId(IDs)
		SendMsgToGroup(remove_ids, msg)
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
	send_client := UserToClient[formid]
	recieve_client, ok := UserToClient[targetId]
	mu.RUnlock()

	// TODO:将消息存储到redis上
	// 选择有序表存储
	ctx := context.Background()
	key := ""
	if formid < targetId {
		key = fmt.Sprintf("msg_%d_to_%d", formid, targetId)
	} else {
		key = fmt.Sprintf("msg_%d_to_%d", targetId, formid)
	}
	// 使用redis的有序集合
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

	// FIXME:这里出问题？
	// 目前来看，等登陆上才会有client，未登陆上就是空指针
	// 但是事实情况是，服务器一直开着，用户注册并登录过后台就会有记录，因此只需要加上一个判断即可
	// mu.RLock()
	// recieve_client,ok:= UserToClient[targetId]
	// mu.RUnlock()

	// if send_client == nil {
	// 	log.Printf("SendMsgToUser: user %d not connected or not exists, message dropped", targetId)
	// 	return
	// }

	// TODO:回显,用于前端显示消息
	if send_client != nil {
		select {
		case send_client.SendDataQueue <- msg:
			// 发送成功
		default:
			// 发送者队列满了，一般这种情况很少见，除非断网卡死
		}
	}

	if ok && recieve_client != nil {
		// 这里的 select 是为了防止写入阻塞导致协程泄露（可选，但推荐）
		select {
		case recieve_client.SendDataQueue <- msg:
			fmt.Println("消息已实时推送给用户", targetId)
		default:
			fmt.Println("用户", targetId, "的发送队列已满，消息仅存储在Redis")
		}
	} else {
		fmt.Println("用户", targetId, "不在线，消息已保存到Redis")
	}

}

// 群聊
func SendMsgToGroup(ids []uint, msg []byte) {
	if len(ids) == 0 {
		return
	}
	// fmt.Println(">>>>>>>>>>>2:",ids)
	for _, v := range ids {
		mu.RLock()
		client := UserToClient[v]
		mu.RUnlock()
		client.SendDataQueue <- msg
	}
}

// 系统广播，比如维护信息
// func SendMsgToAll(){

// }

// 读取start-end的redis数据，并返回给前端
// 返回的是消息列表吗？
func HistoryMsg(redisPayload system.RedisPayload) ([]*Message, error) {
	// 从redis里面查找

	ctx := context.Background()

	// 组装key
	key := " "
	if redisPayload.UserId < redisPayload.TargetId {
		key = fmt.Sprintf("msg_%d_to_%d", redisPayload.UserId, redisPayload.TargetId)
	} else {
		key = fmt.Sprintf("msg_%d_to_%d", redisPayload.TargetId, redisPayload.UserId)
	}
	// 通过zreverange返回有序集中指定区间内的成员，通过索引，分数从高到低
	// msgs是string类型的
	stringmsgs, err := utils.RDB.ZRevRange(ctx, key, redisPayload.Start, redisPayload.End).Result()
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

func (client *Client) RemoveId(a []uint) []uint {
	for i, v := range a {
		if v == client.User_id {
			return append(a[:i], a[i+1:]...)
		}
	}
	return a
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
