package models


// func init(){
// 	go GlobalHub.Management()
// }

// var GlobalHub = Hub{
// 	Register: make(chan *Client),
// 	UnRegister: make(chan *Client),
// 	BroadcastMessage: make([]byte, 0),
// 	UserToClient : make(map[uint]*Client, 0),
// }


// type Hub struct{
// 	// 上线
// 	Register chan *Client
// 	// 下线
// 	UnRegister chan *Client
// 	// 广播消息
// 	BroadcastMessage []byte
// 	// 用户池
// 	UserToClient map[uint]*Client
// }