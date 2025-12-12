package v1

import (
	"encoding/json"
	"ZustChat/model/request"
	"ZustChat/service"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func WsHandler(c *gin.Context) {
	var authPayload request.AuthMessage

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  512,
		WriteBufferSize: 512,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	// 升级为webdsocket连接
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	// 从连接中读取数据
	_, data, err := ws.ReadMessage()
	if err != nil {
		return
	}

	// 反序列化
	err = json.Unmarshal(data, &authPayload)
	if err != nil {
		return
	}

	service.WsConnection(ws, authPayload.UserId)
}