package webscoket

import (
	"chat-room/api/helper"
	"chat-room/api/models"
	"encoding/json"
	"github.com/gorilla/websocket"
	"time"
)

//TODO webscoket 房间事件处理

type Room struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Member   []Member `json:"member"`
	TimeUnix int64    `json:"time_unix"` //创建时间
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
	newWS(user, roomId, "才不管啥名字", EVENT_JOIN)
}

//创建房间
func Create(user models.User, roomId, roomName string) {
	newWS(user, roomId, roomName, EVENT_CREATE)
}

//封装消息
func newWS(user models.User, roomId, roomName string, eventType EventType) {

	var newRoom Room

	newRoom.TimeUnix = time.Now().Unix()
	msg := "加入成功"

	member := Member{
		UserType: VIEWER,
		UserId:   user.Id,
		UserName: user.Name,
	}

	//房间nil是握手的房间
	if roomId != "nil" {
		helper.Debug("更新房间成员")
		//更新房间成员
		newRoom = updateRoomsMember(roomId, roomName, member)
		msg = user.Name + " 加入了房间 " + newRoom.Name
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
			case EVENT_JOIN:
				//创建房间
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
			helper.Debug("member -- ",member)
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

//初始化函数
func init() {
	//死循环都要开线程去搞他
	go chatRoom()
}
