package webscoket

import (
	"chat-room/api/controllers"
	"chat-room/api/helper"
	"encoding/json"
	"strings"

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
			if !helper.Error(err) {
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
				case EVENT_GAME_ANSWER:
					//当前用户是谁
					data := make(map[string]interface{})
					data["user"] = this.User
					m.Data = data

					//TODO 后期要考虑一个用户只能加入一个房间，现在还没实现
					room := getRoom(m.Room.Id)
					//这条消息发给谁呢
					m.Room.Member = room.Member
					//获取自己的身份
					has := false       //是否真的在房间
					needBonus := false //是否需要加分
					for _, u := range room.Member {
						if u.UserId == this.User.Id {
							has = true
							if u.UserType == MASTER {
								//玩野?
							} else {
								needBonus = true
							}
							break
						}
					}
					//有效事件
					if has {
						//替换答案 ***
						msg := m.Msg
						newMsg := strings.Replace(msg, room.KeyWord, "*", -1)
						if newMsg != msg {
							newMsg = "***"
						}
						//推送消息
						m.Msg = newMsg
						publish <- m
						//判断答案
						if msg == room.KeyWord {
							if needBonus {
								//答对啦，奖你一包辣条
								m.EventType = EVENT_GAME_BONUS
								m.Msg = "恭喜" + this.User.Name + "答对了"
								publish <- m
								//获取当前答对人数，计算分数
								//TODO 答对人数应该保存在内存，还没实现
								if m.Room.CorrectNumber == 0 {
									m.Room.CorrectNumber = 1
								}
								m.Msg = "居然答对了！这么厉害~"
								data["Score"] = int(100 / m.Room.CorrectNumber)
								publish <- m
							}
						}
					} else {
						m.Msg = "***"
						m.EventType = EVENT_INVAILD
					}
					break
				default:
					//握手
					member := Member{
						UserType: 0,
						UserId:   this.User.Id,
						UserName: this.User.Name,
					}
					helper.Debug("当前用户:", member)
					ms := make([]Member, 0)
					ms = append(ms, member)
					m.Room.Member = ms
					m.EventType = EVENT_HAND
					m.Msg = "她想握手"
					publish <- m
				}
			}
		}
	}
}
