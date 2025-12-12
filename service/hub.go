package service

var GlobalHub = Hub{
	Register: make(chan *Client),
	UnRegister: make(chan *Client),
	BroadcastMessage: make([]byte, 0),
	UserToClient : make(map[uint]*Client, 0),
}

type Hub struct{
	// 上线
	Register chan *Client
	// 下线
	UnRegister chan *Client
	// 广播消息
	BroadcastMessage []byte
	// 用户池
	UserToClient map[uint]*Client
}

// 管理客户端登入登出，以及系统消息广播
func (hub *Hub)Management(){
	for {
		select {
		case client:=<-hub.Register:
			hub.UserToClient[client.User_id] = client
			go client.Send()
			go client.Recieve()
		case client:=<-hub.UnRegister:
			// 删除
			close(client.SendDataQueue)
			delete(hub.UserToClient,client.User_id)
		// case:广播/待扩展
		}
	}
}