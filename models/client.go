package models

// import (
// 	// "gorm.io/gorm"
// 	"github.com/gorilla/websocket"
// )

// type Client struct {
// 	// 记录用户
// 	User_id uint
// 	// *websocket.Conn 类型的对象。这个对象是与单个客户端进行所有通信的唯一凭证和工具。之后的所有操作都是调用这个 conn 对象的方法。
// 	Conn *websocket.Conn

// 	// 接受客户端的消息
// 	// Msg           Message
// 	HeartbeatTime uint64 //心跳时间
// 	// 客户端邮箱,存储待发送消息
// 	SendDataQueue chan []byte
// }