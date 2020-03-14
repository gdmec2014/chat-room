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
			////helper.Debug("msg --- ", msg)
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
				case EVENT_SYSTEM_MESSAGE:
					publish <- m
					break
				case EVENT_GAME_ANSWER, EVENT_MESSAGE:

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

					//替换答案 ***
					isTrue := false
					msg := m.Msg
					newMsg := strings.ReplaceAll(msg, " ", "")
					if room.KeyWord == newMsg {
						newMsg = "***"
						isTrue = true
						//helper.Debug("答对了")
					} else {
						newMsg = msg
					}
					//推送消息
					m.Msg = newMsg
					publish <- m
					//判断答案
					needBonus := false
					//isTrue = true
					if isTrue {
						//加分
						if m.Room.Mark == nil {
							m.Room.Mark = make(map[int][]Mark, 0)
						}

						mark := Mark{
							Id:        this.User.Id,
							HasAnswer: false,
						}

						masterMark := Mark{
							Id:        0,
							Point:     10,
							HasAnswer: false,
						}

						masterId := int64(0)
						for _, m := range m.Room.Member {
							if m.UserType == MASTER {
								masterMark.Id = m.UserId
								masterId = m.UserId
							}
						}

						mm := make([]Mark, 0)

						hasAdd := false

						helper.Debug("this.uid = ", this.User.Id)

						if _, ok := m.Room.Mark[m.Room.Times]; ok {

							mm = m.Room.Mark[m.Room.Times]

							hasAnswer := false

							for _, m := range mm {
								if m.Id == this.User.Id {
									hasAdd = true
									hasAnswer = m.HasAnswer
									break
								}
							}

							if !hasAdd {
								needBonus = true
								correctNumber := len(m.Room.Mark[m.Room.Times])
								correctNumber++
								mark.Point += int(100 / correctNumber)
								mark.HasAnswer = true
								m.Room.Mark[m.Room.Times] = append(m.Room.Mark[m.Room.Times], mark)
							} else {
								//如果已经添加过了，需要判断是否回答过正确的答案
								if !hasAnswer {
									needBonus = true
								}
							}
						} else {
							needBonus = true
							m.Room.Mark[m.Room.Times] = make([]Mark, 0)
							mark.Point += 100
							mark.HasAnswer = true
							m.Room.Mark[m.Room.Times] = append(m.Room.Mark[m.Room.Times], mark)
						}

						helper.Debug("masterId:", masterId)
						if needBonus {

							//答对啦，奖你一包辣条
							m.EventType = EVENT_SYSTEM_MESSAGE
							m.Msg = "恭喜" + this.User.Name + "答对了"
							publish <- m

							mm = m.Room.Mark[m.Room.Times]
							hasMaster := false

							for i, _ := range mm {
								if mm[i].Id == masterId {
									hasMaster = true
									mm[i].Point += 10
									mm[i].HasAnswer = true
								}
							}

							m.Room.Mark[m.Room.Times] = mm

							if !hasMaster {
								m.Room.Mark[m.Room.Times] = append(m.Room.Mark[m.Room.Times], masterMark)
							}

							answers := 0
							for _, m := range mm {
								//记录不是出题者的人
								if m.Id != masterId {
									answers++
								}
							}

							//判断本轮是不是全部人都答对了
							if answers == len(m.Room.Member)-1 {
								//应该结束本轮
								helper.Debug("应该结束本轮")
								m.Room.IsOver = true
								//发送上个答案
								m.EventType = EVENT_SYSTEM_MESSAGE
								m.Msg = "答案：" + room.KeyWord
								publish <- m
							}

							m.Room.MasterId = masterId

							updateRooms(m.Room)

						}
					}

					break
				default:
					//握手
					member := Member{
						UserType: 0,
						UserId:   this.User.Id,
						UserName: this.User.Name,
					}
					//helper.Debug("当前用户:", member)
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
