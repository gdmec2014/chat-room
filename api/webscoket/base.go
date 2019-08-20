package webscoket

import (
	"chat-room/api/controllers"
	"chat-room/api/helper"
	"encoding/json"

	"github.com/gorilla/websocket"
)

// WebSocketController
type WebSocketController struct {
	controllers.NeedLoginController
}

func (this *WebSocketController) join() {

	// 开启链接
	ws, err := websocket.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		this.SetReturnData(helper.FAILED, "1.链接失败了", err.Error())
	} else if err != nil {
		this.SetReturnData(helper.FAILED, "2.链接失败了", err.Error())
	}

	// 加入房间
	this.User.Conn = ws
	updateUserConn(this.User)

	// 轮询读取消息
	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			return
		}
		if len(p) > 3 {
			var m Event
			msg := string(p)
			helper.Debug("msg --- ", msg)
			err := json.Unmarshal([]byte(msg), &m)
			if err == nil {
				switch m.EventType {
				case EVENT_CREATE:
					if len(m.Room.Name) < 1 {
						m.EventType = EVENT_INVAILD
						m.Msg = "你需要输入房间名字才能创建喔"
						publish <- m
						return
					}
					roomId := helper.GetRandomString(16)
					Create(this.User, roomId, m.Room.Name)
					break
				case EVENT_JOIN:
					Join(this.User, m.Room.Id)
					break
				case EVENT_MESSAGE:
					this.User.Token = ""
					m.Data = this.User
					publish <- m
					break
				case EVENT_DRAW:
					data := make(map[string]interface{})
					data["user"] = this.User
					data["position"] = m.Data
					m.Data = data
					publish <- m
					break
				case EVENT_BREAK_DRAW:
					data := make(map[string]interface{})
					data["user"] = this.User
					m.Data = data
					publish <- m
					break
				default:
					//握手
					member := Member{
						UserType: 0,
						UserId:   this.User.Id,
						UserName: this.User.Name,
					}
					helper.Debug("当前用户:", member)
					ms := make([]Member,0)
					ms = append(ms,member)
					m.Room.Member = ms
					m.EventType = EVENT_HAND
					m.Msg = "她想握手"
					publish <- m
				}
			}
		}
	}
}
