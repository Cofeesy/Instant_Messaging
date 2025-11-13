package models

import (
	"encoding/json"
	"fmt"
	"gin_chat/utils/setting"
	"log"
	// "net/http"
	"sync"
	"time"

	// "github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 消息源码
type Message struct {
	FromId      int
	TargetId    int
	Content     string
	MessageType int //0私发 1群发 2广播
	ContentType int //文字 图片 表情包
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
		Conn: ws,
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
	str := "欢迎进入聊天系统"
	strbyte, _ := json.Marshal(str)
	SendMsgToUser(userid, strbyte)

}

// 解析收到的消息

// 编码待发送消息

// 每个客户端一直监听即将发送的信息
func (client *Client) Send() {
	// client.send <- []byte("你好啊")
	for {
		select {
		case msg := <-client.SendDataQueue:
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
}

// 监听听并读取接受到的信息，需要一个协程持续监听，这是每一个连接进来的客户端需要做的
func (client *Client) Recieve() {
	// 接受监听消息
	var Msg Message;
	for {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			fmt.Println("read1:", err.Error())
		}
		// 将数据序列化为代码形式
		// 需要加指针
		err = json.Unmarshal(msg, &Msg)
		if err != nil {
			fmt.Println("read2:", err.Error())
		}
		fmt.Println("Received message: ", Msg)

		// FIXME:解释
		if Msg.Content == "心跳" {
			client.HeartbeatTime = uint64(time.Now().Unix())
			continue
		}

		// 调度
	}
}

// 私聊
// 写消息进去
// 主要是写是建立在连接双方的写，通过把消息放到其他人的send上，服务器读取到会将该消息发送到该send连接到的client上
// 这里的“写” (conn.WriteMessage)，最终都是发生在某一个具体的 WebSocket 连接 (conn) 上的。
// 即clientA 的 Send goroutine 只会在 clientA.conn 上写，clientB 的 Send goroutine 只会在 clientB.conn 上写。职责非常明确
// 关键点是看谁开启了这个协程，服务器只是负责发送和接受信息
func SendMsgToUser(userId uint, msg []byte) {

	mu.RLock()
	client := UserToClient[userId]
	mu.RUnlock()

	if client == nil {
		log.Printf("SendMsgToUser: user %d not connected, message dropped", userId)
		return
	}

	select {
	case client.SendDataQueue <- msg:
	default:
		log.Printf("SendMsgToUser: user's send channel full, dropping message for user %d", userId)
	}

}

// 群聊单发
func SendMsgToGroup(GroupId int, msg Message) {

}

// 系统广播，比如维护信息
// func SendMsgToAll(){

// }

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
	if client.HeartbeatTime+uint64(setting.HeartbeatMaxTime) <= currentTime {
		fmt.Println("心跳超时...自动下线", client)
		timeout = true
	}
	return
}
