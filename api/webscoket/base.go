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
			//helper.Debug("msg --- ", msg)
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
				case EVENT_DRAW, EVENT_NEW_DRAW:
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
				case EVENT_GAME_ANSWER, EVENT_MESSAGE:
					helper.Debug("替换答案 ***")
					this.User.Token = ""
					m.Data = this.User

					//当前用户是谁
					data := make(map[string]interface{})
					data["user"] = this.User

					//TODO 后期要考虑一个用户只能加入一个房间，现在还没实现
					room := getRoom(m.Room.Id)
					//这条消息发给谁呢
					m.Room = room

					if !isGameStart(room) {
						publish <- m
						return
					}

					//获取自己的身份
					has := false       //是否真的在房间
					needBonus := false //是否需要加分
					for _, u := range room.Member {
						if u.UserId == this.User.Id {
							helper.DebugStructToString(u)
							has = true
							if u.UserType == MASTER {
								//玩野?
							} else if u.UserType == PLAYER {
								needBonus = true
							}
							break
						}
					}

					//有效事件
					if has {
						//替换答案 ***
						isTrue := false
						msg := m.Msg
						newMsg := strings.Replace(msg, room.KeyWord, "*", -1)
						if newMsg != msg {
							newMsg = "***"
							isTrue = true
						}
						//推送消息
						m.Msg = newMsg
						publish <- m
						//判断答案
						//isTrue = true
						//needBonus = true
						if isTrue {
							if needBonus {
								//答对啦，奖你一包辣条
								m.EventType = EVENT_GAME_BONUS
								m.Msg = "恭喜" + this.User.Name + "答对了"
								publish <- m

								m.Msg = "居然答对了！这么厉害~"

								mark := Mark{
									Id: this.User.Id,
								}

								//加分
								if m.Room.Mark == nil {
									m.Room.Mark = make(map[int][]Mark, 0)
								}
								if _, ok := m.Room.Mark[m.Room.Times]; ok {
									//存在
									mm := m.Room.Mark[m.Room.Times]
									hasAdd := false
									for i, _ := range mm {
										if mm[i].Id == this.User.Id {
											//已经加过分了
											hasAdd = true
										}
									}
									if !hasAdd {
										correctNumber := len(m.Room.Mark[m.Room.Times])
										correctNumber++
										mark.Point = int(100 / correctNumber)
										m.Room.Mark[m.Room.Times] = append(m.Room.Mark[m.Room.Times], mark)
										updateRooms(m.Room)
									}
								} else {
									m.Room.Mark[m.Room.Times] = make([]Mark, 0)
									mark.Point = 100
									m.Room.Mark[m.Room.Times] = append(m.Room.Mark[m.Room.Times], mark)
									updateRooms(m.Room)
								}
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
