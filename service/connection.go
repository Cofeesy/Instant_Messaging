package service

import (
	"fmt"
	"gin_chat/utils"
	"sync"
	"time"
	"github.com/gorilla/websocket"
)

var mu sync.RWMutex

// 启动项
func init(){
	go GlobalHub.Management()
}

// websocket客户端-->映射user
type Client struct {
	// 用户id
	User_id uint
	// *websocket.Conn 类型的对象。这个对象是与单个客户端进行所有通信的唯一凭证和工具。之后的所有操作都是调用这个 conn 对象的方法。
	Conn *websocket.Conn
	HeartbeatTime uint64 //心跳时间
	// 客户端邮箱,存储待发送消息
	SendDataQueue chan []byte
}

// 监控发送和接受
func WsConnection(ws *websocket.Conn, userid uint) {
	current_time := uint64(time.Now().Unix())
	client := &Client{
		User_id: userid,
		Conn:    ws,
		// Msg:  msg,
		HeartbeatTime: current_time,
		SendDataQueue: make(chan []byte, 256), // 创建带缓冲区的信箱
	}

	GlobalHub.Register<-client
}


// 每个客户端一直监听即将发送的信息
func (client *Client) Send() {
	for msg := range client.SendDataQueue {
		err := client.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			fmt.Println("write2:", err.Error())
		}
	}
}

// 监听听并读取接受到的信息，需要一个协程持续监听，这是每一个连接进来的客户端需要做的
func (client *Client) Recieve() {
	// 异常捕捉
	defer func() {
        GlobalHub.UnRegister <- client
		// 断开wesocket连接
		client.Conn.Close()
    }()
	
	for {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			fmt.Println("read1:", err.Error())
			break
		}

		// 消息调度
		dispatchMsg(msg, client)
	}
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
	currentTime := uint64(time.Now().Unix())
	for i := range GlobalHub.UserToClient {
		client := GlobalHub.UserToClient[i]
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
	if client.HeartbeatTime+uint64(utils.HeartbeatMaxTime) <= currentTime {
		fmt.Println("心跳超时...自动下线", client)
		timeout = true
	}
	return
}