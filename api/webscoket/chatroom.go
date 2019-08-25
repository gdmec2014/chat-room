package webscoket

import (
	"chat-room/api/helper"
	"chat-room/api/models"
	"encoding/json"
	"github.com/gorilla/websocket"
	"math/rand"
	"time"
)

//TODO webscoket 房间事件处理

type Room struct {
	Id            string   `json:"id"`
	Name          string   `json:"name"`
	Member        []Member `json:"member"`
	TimeUnix      int64    `json:"time_unix"`      //创建时间
	KeyWord       string   `json:"key_word"`       //当前游戏正确答案
	CorrectNumber int      `json:"correct_number"` //当前轮游戏回答正确人数
	MaxMember     int      `json:"max_member"`     //游戏最大人数
}

type Member struct {
	UserType UserType `json:"user_type"` // 用户类型跟房间，因为不是每一个房间的身份都一样
	UserId   int64    `json:"user_id"`   // 用户ID
	UserName string   `json:"user_name"` // 用户名
}

type Event struct {
	EventType EventType   `json:"event_type"` // 消息类型
	Room      Room        `json:"room"`       // 房间       //前端發來的字段
	Msg       string      `json:"msg"`        // 消息
	TimeUnix  int64       `json:"time_unix"`  // 消息时间戳
	Data      interface{} `json:"data"`       // 附带数据    //返回後端的字段
}

var (
	// 推送消息的阻塞通道
	publish = make(chan Event, 10)
)

//加入房间
func Join(user models.User, roomId string) {
	room := Room{
		Id:            roomId,
		Name:          "才不知道是啥名字",
		TimeUnix:      time.Now().Unix(),
		KeyWord:       "",
		CorrectNumber: 0,
		MaxMember:     6,
	}
	helper.Debug("Join", room)
	newWS(user, room, EVENT_JOIN)
}

//创建房间
func Create(user models.User, roomId, roomName string, maxMember int) {
	room := Room{
		Id:            roomId,
		Name:          roomName,
		TimeUnix:      time.Now().Unix(),
		KeyWord:       "",
		CorrectNumber: 0,
		MaxMember:     maxMember,
	}
	newWS(user, room, EVENT_CREATE)
}

//封装消息
func newWS(user models.User, room Room, eventType EventType) {

	var newRoom Room

	newRoom.TimeUnix = time.Now().Unix()
	msg := "加入成功"

	member := Member{
		UserType: VIEWER,
		UserId:   user.Id,
		UserName: user.Name,
	}

	//房间nil是握手的房间
	if room.Id != "nil" {
		helper.Debug("更新房间成员")
		//更新房间成员
		var code EventType
		newRoom, code = updateRoomsMember(room, member)
		msg = user.Name + " 加入了房间 " + newRoom.Name

		if eventType != EVENT_CREATE {
			eventType = code
			switch code {
			case EVENT_GAME_CAN_START:
				//开始游戏
				msg = "准备好了么~要开始了喔~"
				//赋予玩家身份
				for i, m := range newRoom.Member {
					if i == 0 {
						m.UserType = MASTER
					} else {
						m.UserType = PLAYER
					}
				}
				//随机获取问题答案
				vocabulary := models.Vocabulary
				lenVocabulary := len(vocabulary)
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				key := r.Intn(lenVocabulary)
				newRoom.KeyWord = vocabulary[key]
				go startGame(newRoom)
				break
			case EVENT_GAME_NO_START:
				//人数还不够，不可以开始喔
				msg = "人数还不够，不可以开始喔~"
				break
			case EVENT_NO_PLACE:
				//房间不能加人了
				msg = "房间已经满人了喔！"
				break
			}
		}

	} else {
		eventType = EVENT_INVAILD
		msg = "无效事件"
		helper.Debug("无效事件~")
		members := make([]Member, 0)
		members = append(members, member)
		newRoom = Room{
			Id:     "nil",
			Name:   "无效事件",
			Member: members,
		}
	}

	event := Event{
		Data:      user,
		TimeUnix:  time.Now().Unix(),
		Msg:       msg,
		EventType: eventType,
		Room:      newRoom}

	//推送数据到浏览器
	publish <- event
}

//此函数处理所有传入的阻塞消息。
func chatRoom() {
	for {
		select {
		case event := <-publish:
			switch event.EventType {
			case EVENT_HAND:
				helper.Debug("握手")
				event.EventType = EVENT_HAND
				event.Msg = "握手成功"
				break
			case EVENT_JOIN, EVENT_GAME_NO_START:
				//加入房间
				helper.Debug("加入房间")
				break
			case EVENT_CREATE:
				//创建房间
				helper.Debug("创建房间")
				event.Room.Id = helper.GetRandomString(16)
				break
			case EVENT_MESSAGE:
				helper.Debug("发送信息")
				break
			case EVENT_DRAW:
				event.Msg = "啊！请接收绘图"
				helper.Debug("绘图")
				break
			case EVENT_BREAK_DRAW:
				event.Msg = "啊！请中断绘图"
				helper.Debug("中断绘图")
				break
			default:
				//握手时候，没有房间号
				helper.Debug("假装握手")
				event.EventType = EVENT_HAND
				event.Msg = "握手成功"
			}
			broadcastWebSocket(event)
		}
	}
}

// 广播消息
func broadcastWebSocket(event Event) {

	event.TimeUnix = time.Now().Unix()

	data, err := json.Marshal(event)
	if helper.Error(err) {
		helper.Error("broadcastWebSocket 发生错误，不能发送消息")
		return
	}

	var room Room
	var member []models.User

	switch event.EventType {
	case EVENT_CREATE:
		//创建了房间，要通知所有人更新房间列表
		member = getAllMember()
		break
	case EVENT_HAND:
		helper.Debug("新建房间，握手，需要发送数据给自己 -> ", event.Room.Member[0].UserId)
		//新建房间，握手，需要发送数据给自己
		if len(event.Room.Member) > 0 {
			var has bool
			has, u := hasMember(event.Room.Member[0].UserId)
			if has {
				member = append(member, u)
			} else {
				helper.Error("broadcastWebSocket 发生错误，用户数据丢失")
				return
			}
		}
		break
	default:
		room = getRoom(event.Room.Id)
		if len(room.Member) > 0 {
			member = getMemberByRoom(room)
			helper.Debug("member -- ", member)
		} else {
			//已经没有用户了，应该销毁他
			return
		}
	}

	for _, m := range member {
		ws := m.Conn
		if ws != nil {
			if ws.WriteMessage(websocket.TextMessage, data) != nil {
				//发生错误，这里应该作重连机制
				ws.Close()
			}
		}
	}
}

//开始游戏
func startGame(room Room) {
	t := 0
	t1 := time.NewTimer(time.Second * 1)
	event := Event{
		EventType: 0,
		Room:      room,
		Msg:       "嘿嘿嘿",
		TimeUnix:  0,
		Data:      nil,
	}
	for {
		select {
		case <-t1.C:
			t++
			if t == 90 {
				//开始新的一轮游戏啦
				t = 0
				//转换身份事件
				isOver := 0 //是否应该结束游戏呢? =6 => true
				hasMaster := false
				for i, m := range room.Member {
					switch m.UserType {
					case MASTER:
						event.Room.Member[i].UserType = NO_MASTER
						isOver++
						break
					case NO_MASTER:
						isOver++
						break
					case PLAYER:
						if !hasMaster {
							hasMaster = true
							event.Room.Member[i].UserType = MASTER
						}
						break
					}
				}
				if isOver > len(event.Room.Member) {
					//中断游戏
					event.EventType = EVENT_GAME_OVER
					broadcastWebSocket(event)
					break
				} else {
					//更新身份
					event.EventType = EVENT_GIVE_IDENTITY
					broadcastWebSocket(event)
				}
			} else {
				//推送本轮游戏剩余时间
				event.EventType = EVENT_GAME_TIME
				event.Data = 90 - t
				broadcastWebSocket(event)
			}
			t1.Reset(time.Second * 1)
		}
	}
}

//初始化函数
func init() {
	//死循环都要开线程去搞他
	go chatRoom()
}
