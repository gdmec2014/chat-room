package webscoket

import (
	"chat-room/api/helper"
	"encoding/json"

	"chat-room/api/controllers"

	"github.com/gorilla/websocket"
)

// WebSocketController
type WebSocketController struct {
	controllers.NeedLoginController
}

// 广播消息
func broadcastWebSocket(event Event) {

	data, err := json.Marshal(event)
	if helper.Error(err) {
		helper.Error("broadcastWebSocket 发生错误，不能发送消息")
		return
	}

	subscriber := getRoomUsers(event.RoomId)

	for _, sub := range subscriber {
		ws := sub.Conn
		if ws != nil {
			if ws.WriteMessage(websocket.TextMessage, data) != nil {
				//发生错误，这里应该作重连机制
				//但我不想，直接断了他
				unsubscribe <- sub.User.Id
			}
		}
	}
}

func (this *WebSocketController) join(roomId string, userType UserType) {
	// 开启链接
	ws, err := websocket.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		this.SetReturnData(helper.FAILED, "1.链接失败了", err.Error())
	} else if err != nil {
		this.SetReturnData(helper.FAILED, "2.链接失败了", err.Error())
	}

	// 加入房间
	Join(this.User, ws, roomId, userType)
	defer Leave(this.User.Id)

	// 轮询读取消息
	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			return
		}
		publish <- newEvent(EVENT_MESSAGE, this.User, string(p), roomId)
	}
}
