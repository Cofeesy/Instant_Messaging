package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 关键：添加 CheckOrigin 函数，并总是返回 true 来允许所有来源的连接
	// CheckOrigin: func(r *http.Request) bool {
	//     return true
	// },
}

func wsHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// 输出客户端发送的消息到控制台
		log.Printf("Received message: %s\n", p)

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func main() {
	r := gin.Default()
	r.GET("/", wsHandler)

	log.Println("Server is running on :3000...")
	if err := r.Run(":3000"); err != nil {
        log.Fatal("Failed to run server:", err)
    }
}
